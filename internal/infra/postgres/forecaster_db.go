package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/babadro/forecaster/pkg/fcasterbot"
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

func (db *ForecasterDB) GetPollByID(ctx context.Context, id int32) (fcasterbot.Poll, error) {
	pollSQL, _, err := db.q.Select("*").From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("select poll", err)
	}

	var poll fcasterbot.Poll
	err = db.db.
		QueryRow(ctx, pollSQL, id).
		Scan(&poll.ID, &poll.Title, &poll.Description, &poll.Start, &poll.Finish, &poll.UpdatedAt)

	if err != nil {
		return fcasterbot.Poll{}, scanFailed("select poll", err)
	}

	optionsSQL, _, err := db.q.Select("*").From("forecaster.options").Where(sq.Eq{"poll_id": id}).ToSql()
	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("select options", err)
	}

	rows, err := db.db.Query(ctx, optionsSQL)
	if err != nil {
		return fcasterbot.Poll{}, queryFailed("select options", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option fcasterbot.Option

		err = rows.Scan(&option.ID, &option.Title, &option.Description, &option.PollID)
		if err != nil {
			return fcasterbot.Poll{}, scanFailed("select options", err)
		}

		poll.Options = append(poll.Options, option)
	}

	if err = rows.Err(); err != nil {
		return fcasterbot.Poll{}, rowsError("select options", err)
	}

	return poll, nil
}

func (db *ForecasterDB) CreateSeries(ctx context.Context, s fcasterbot.Series) (fcasterbot.Series, error) {
	seriesSQL, args, err := db.q.
		Insert("forecaster").Columns("title", "description").
		Values(s.Title, s.Description).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fcasterbot.Series{}, buildingQueryFailed("insert series", err)
	}

	err = db.db.QueryRow(ctx, seriesSQL, args...).Scan(&s.ID, &s.Title, &s.Description)
	if err != nil {
		return fcasterbot.Series{}, scanFailed("insert series", err)
	}

	return s, nil
}

func (db *ForecasterDB) CreatePoll(ctx context.Context, poll fcasterbot.Poll) (fcasterbot.Poll, error) {
	pollSQL, args, err := db.q.
		Insert("forecaster.polls").
		Columns("title", "description", "start", "finish", "created_at", "updated_at").
		Values(poll.Title, poll.Description, poll.Start, poll.Finish).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()

	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("insert poll", err)
	}

	err = db.db.QueryRow(ctx, pollSQL, args...).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return fcasterbot.Poll{}, scanFailed("insert poll", err)
	}

	return poll, nil
}

func (db *ForecasterDB) CreateOption(ctx context.Context, option fcasterbot.Option) (fcasterbot.Option, error) {
	optionSQL, args, err := db.q.
		Insert("forecaster.options").
		Columns("title", "description", "poll_id").
		Values(option.Title, option.Description, option.PollID).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fcasterbot.Option{}, buildingQueryFailed("insert option", err)
	}

	err = db.db.QueryRow(ctx, optionSQL, args...).Scan(&option.ID)
	if err != nil {
		return fcasterbot.Option{}, scanFailed("insert option", err)
	}

	return option, nil
}

func (db *ForecasterDB) UpdateSeries(ctx context.Context, s fcasterbot.Series) (fcasterbot.Series, error) {
	seriesSQL, args, err := db.q.
		Update("forecaster.series").
		Set("title", s.Title).
		Set("description", s.Description).
		Where(sq.Eq{"id": s.ID}).
		Suffix("RETURNING updated_at").
		ToSql()

	if err != nil {
		return fcasterbot.Series{}, buildingQueryFailed("update series", err)
	}

	var updatedAt time.Time

	err = db.db.QueryRow(ctx, seriesSQL, args...).Scan(&updatedAt)
	if err != nil {
		return fcasterbot.Series{}, scanFailed("update series", err)
	}

	s.UpdatedAt = updatedAt

	return s, nil
}

func (db *ForecasterDB) UpdatePoll(ctx context.Context, poll fcasterbot.Poll) (fcasterbot.Poll, error) {
	pollSQL, args, err := db.q.
		Update("forecaster.polls").
		Set("title", poll.Title).
		Set("description", poll.Description).
		Set("start", poll.Start).
		Set("finish", poll.Finish).
		Where(sq.Eq{"id": poll.ID}).
		Suffix("RETURNING updated_at").
		ToSql()

	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("update poll", err)
	}

	var updatedAt time.Time

	err = db.db.QueryRow(ctx, pollSQL, args...).Scan(&updatedAt)
	if err != nil {
		return fcasterbot.Poll{}, scanFailed("update poll", err)
	}

	poll.UpdatedAt = updatedAt

	return poll, nil
}

func (db *ForecasterDB) UpdateOption(ctx context.Context, option fcasterbot.Option) (fcasterbot.Option, error) {
	optionSQL, args, err := db.q.
		Update("forecaster.options").
		Set("title", option.Title).
		Set("description", option.Description).
		Where(sq.Eq{"id": option.ID}).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fcasterbot.Option{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	err = db.db.QueryRow(ctx, optionSQL, args...).Scan(&option.ID)
	if err != nil {
		return fcasterbot.Option{}, fmt.Errorf("unable to update option: %w", err)
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
