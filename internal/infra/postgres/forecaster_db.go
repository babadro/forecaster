package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/babadro/forecaster/internal/domain"
	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lann/builder"

	sq "github.com/Masterminds/squirrel"
)

type ForecasterDB struct {
	db *pgxpool.Pool
	q  sq.StatementBuilderType
}

func NewForecasterDB(db *pgxpool.Pool) *ForecasterDB {
	return &ForecasterDB{
		db: db,
		q:  sq.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(sq.Dollar),
	}
}

func (db *ForecasterDB) GetSeriesByID(ctx context.Context, id int32) (models.Series, error) {
	seriesSQL, _, err := db.q.Select(
		"id", "title", "description", "created_at", "updated_at",
	).From("forecaster.series").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return models.Series{}, buildingQueryFailed("select series", err)
	}

	var series models.Series
	err = db.db.
		QueryRow(ctx, seriesSQL, id).
		Scan(&series.ID, &series.Title, &series.Description, &series.CreatedAt, &series.UpdatedAt)

	selectSeries := "select series"
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Series{}, errNotFound(selectSeries, err)
		}

		return models.Series{}, scanFailed(selectSeries, err)
	}

	return series, nil
}

func (db *ForecasterDB) GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error) {
	pollSQL, args, err := db.q.Select(
		"id", "series_id", "title", "description", "start", "finish", "created_at", "updated_at",
	).From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return models.PollWithOptions{}, buildingQueryFailed("select poll", err)
	}

	var poll models.PollWithOptions
	err = db.db.
		QueryRow(ctx, pollSQL, args...).
		Scan(
			&poll.ID, &poll.SeriesID, &poll.Title, &poll.Description, &poll.Start, &poll.Finish, &poll.CreatedAt, &poll.UpdatedAt,
		)

	selectPoll := "select poll"
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PollWithOptions{}, errNotFound(selectPoll, err)
		}

		return models.PollWithOptions{}, scanFailed(selectPoll, err)
	}

	optionsSQL, args, err := db.q.Select(
		"id", "poll_id", "title", "description",
	).From("forecaster.options").Where(sq.Eq{"poll_id": id}).ToSql()
	if err != nil {
		return models.PollWithOptions{}, buildingQueryFailed("select options", err)
	}

	rows, err := db.db.Query(ctx, optionsSQL, args...)
	if err != nil {
		return models.PollWithOptions{}, queryFailed("select options", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option models.Option

		err = rows.Scan(&option.ID, &option.Title, &option.Description, &option.PollID)
		if err != nil {
			return models.PollWithOptions{}, scanFailed("select options", err)
		}

		poll.Options = append(poll.Options, &option)
	}

	if err = rows.Err(); err != nil {
		return models.PollWithOptions{}, rowsError("select options", err)
	}

	return poll, nil
}

func (db *ForecasterDB) CreateSeries(ctx context.Context, s models.CreateSeries) (res models.Series, err error) {
	now := time.Now()

	seriesSQL, args, err := db.q.
		Insert("forecaster.series").Columns("title", "description", "updated_at", "created_at").
		Values(s.Title, s.Description, now, now).
		Suffix("RETURNING id, title, description, created_at, updated_at").
		ToSql()

	if err != nil {
		return models.Series{}, buildingQueryFailed("insert series", err)
	}

	err = db.db.QueryRow(ctx, seriesSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.UpdatedAt, &res.CreatedAt)
	if err != nil {
		return models.Series{}, scanFailed("insert series", err)
	}

	return res, nil
}

func (db *ForecasterDB) CreatePoll(ctx context.Context, poll models.CreatePoll) (res models.Poll, err error) {
	now := time.Now()

	pollSQL, args, err := db.q.
		Insert("forecaster.polls").
		Columns("title", "description", "start", "finish", "created_at", "updated_at").
		Values(poll.Title, poll.Description, poll.Start, poll.Finish, now, now).
		Suffix("RETURNING id, title, description, start, finish, created_at, updated_at").
		ToSql()

	if err != nil {
		return models.Poll{}, buildingQueryFailed("insert poll", err)
	}

	err = db.db.QueryRow(ctx, pollSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.Start, &res.Finish, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return models.Poll{}, scanFailed("insert poll", err)
	}

	return res, nil
}

func (db *ForecasterDB) CreateOption(ctx context.Context, option models.CreateOption) (res models.Option, err error) {
	optionSQL, args, err := db.q.
		Insert("forecaster.options").
		Columns("poll_id", "title", "description").
		Values(option.PollID, option.Title, option.Description).
		Suffix("RETURNING id, poll_id, title, description").
		ToSql()

	if err != nil {
		return models.Option{}, buildingQueryFailed("insert option", err)
	}

	err = db.db.QueryRow(ctx, optionSQL, args...).
		Scan(&res.ID, &res.PollID, &res.Title, &res.Description)
	if err != nil {
		return models.Option{}, scanFailed("insert option", err)
	}

	return res, nil
}

func (db *ForecasterDB) UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries) (res models.Series, err error) {
	b := db.q.Update("forecaster.series").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, title, description, updated_at, created_at")

	if s.Title != nil {
		b = b.Set("title", s.Title)
	}

	if s.Description != nil {
		b = b.Set("description", s.Description)
	}

	seriesSQL, args, err := b.ToSql()

	if err != nil {
		return models.Series{}, buildingQueryFailed("update series", err)
	}

	err = db.db.QueryRow(ctx, seriesSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return models.Series{}, scanFailed("update series", err)
	}

	return res, err
}

func (db *ForecasterDB) UpdatePoll(ctx context.Context, in models.UpdatePoll) (res models.Poll, err error) {
	b := db.q.Update("forecaster.polls").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": in.ID}).
		Suffix("RETURNING id, series_id, title, description, start, finish, updated_at, created_at")

	if in.SeriesID != nil {
		b.Set("series_id", in.SeriesID)
	}

	if in.Title != nil {
		b.Set("title", in.Title)
	}

	if in.Description != nil {
		b.Set("description", in.Description)
	}

	if in.Start != nil {
		b.Set("start", in.Start)
	}

	if in.Finish != nil {
		b.Set("finish", in.Finish)
	}

	pollSQL, args, err := b.ToSql()

	if err != nil {
		return models.Poll{}, buildingQueryFailed("update poll", err)
	}

	err = db.db.QueryRow(ctx, pollSQL, args...).
		Scan(
			&res.ID, &res.SeriesID, &res.Title, &res.Description, &res.Start, &res.Finish, &res.UpdatedAt, &res.CreatedAt,
		)
	if err != nil {
		return models.Poll{}, scanFailed("update poll", err)
	}

	return res, nil
}

func (db *ForecasterDB) UpdateOption(ctx context.Context, in models.UpdateOption) (res models.Option, err error) {
	b := db.q.
		Update("forecaster.options").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": in.ID}).
		Suffix("RETURNING id")

	if in.Title != nil {
		b.Set("title", in.Title)
	}

	if in.Description != nil {
		b.Set("description", in.Description)
	}

	optionSQL, args, err := b.ToSql()

	if err != nil {
		return models.Option{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	err = db.db.QueryRow(ctx, optionSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description)
	if err != nil {
		return models.Option{}, fmt.Errorf("unable to update option: %w", err)
	}

	return res, nil
}

func (db *ForecasterDB) DeleteSeries(ctx context.Context, id int32) error {
	seriesSQL, args, err := db.q.
		Delete("forecaster.series").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return buildingQueryFailed("delete series", err)
	}

	_, err = db.db.Exec(ctx, seriesSQL, args...)
	if err != nil {
		return execFailed("delete series", err)
	}

	return nil
}

func (db *ForecasterDB) DeletePoll(ctx context.Context, id int32) error {
	pollSQL, args, err := db.q.
		Delete("forecaster.polls").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return buildingQueryFailed("delete poll", err)
	}

	_, err = db.db.Exec(ctx, pollSQL, args...)
	if err != nil {
		return execFailed("delete poll", err)
	}

	return nil
}

func (db *ForecasterDB) DeleteOption(ctx context.Context, id int32) error {
	optionSQL, args, err := db.q.
		Delete("forecaster.options").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return buildingQueryFailed("delete option", err)
	}

	_, err = db.db.Exec(ctx, optionSQL, args...)
	if err != nil {
		return execFailed("delete option", err)
	}

	return nil
}

func buildingQueryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: building query failed: %v", queryName, err)
}

func queryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: query failed: %v", queryName, err)
}

func rowsError(queryName string, err error) error {
	return fmt.Errorf("%s: rows error: %v", queryName, err)
}

func scanFailed(queryName string, err error) error {
	return fmt.Errorf("%s: scan rows failed: %v", queryName, err)
}

func execFailed(queryName string, err error) error {
	return fmt.Errorf("%s: exec failed: %v", queryName, err)
}

func errNotFound(queryName string, err error) error {
	return fmt.Errorf("%w: %s", domain.ErrNotFound, err)
}
