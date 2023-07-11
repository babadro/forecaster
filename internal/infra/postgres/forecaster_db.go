package postgres

import (
	"context"
	"fmt"
	"time"

	models "github.com/babadro/forecaster/internal/models/swagger"
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

func (db *ForecasterDB) GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error) {
	pollSQL, _, err := db.q.Select("*").From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return models.PollWithOptions{}, buildingQueryFailed("select poll", err)
	}

	var poll models.PollWithOptions
	err = db.db.
		QueryRow(ctx, pollSQL, id).
		Scan(&poll.ID, &poll.Title, &poll.Description, &poll.Start, &poll.Finish, &poll.UpdatedAt)

	if err != nil {
		return models.PollWithOptions{}, scanFailed("select poll", err)
	}

	optionsSQL, _, err := db.q.Select("*").From("forecaster.options").Where(sq.Eq{"poll_id": id}).ToSql()
	if err != nil {
		return models.PollWithOptions{}, buildingQueryFailed("select options", err)
	}

	rows, err := db.db.Query(ctx, optionsSQL)
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
		Insert("forecaster").Columns("title", "description", "updated_at", "created_at").
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

func (db *ForecasterDB) UpdateSeries(ctx context.Context, s models.UpdateSeries) (res models.Series, err error) {
	b := db.q.Update("forecaster.series").
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": s.ID}).
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

func (db *ForecasterDB) UpdateOption(ctx context.Context, option models.Option) (models.Option, error) {
	optionSQL, args, err := db.q.
		Update("forecaster.options").
		Set("title", option.Title).
		Set("description", option.Description).
		Where(sq.Eq{"id": option.ID}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return models.Option{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	err = db.db.QueryRow(ctx, optionSQL, args...).Scan(&option.ID)
	if err != nil {
		return models.Option{}, fmt.Errorf("unable to update option: %w", err)
	}

	return option, nil
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
	return fmt.Errorf("%s: building query failed: %w", queryName, err)
}

func queryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: query failed: %w", queryName, err)
}

func rowsError(queryName string, err error) error {
	return fmt.Errorf("%s: rows error: %w", queryName, err)
}

func scanFailed(queryName string, err error) error {
	return fmt.Errorf("%s: scan rows failed: %w", queryName, err)
}

func execFailed(queryName string, err error) error {
	return fmt.Errorf("%s: exec failed: %w", queryName, err)
}
