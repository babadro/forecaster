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

func (db *ForecasterDB) GetByID(ctx context.Context, id int32) (fcasterbot.Poll, error) {
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

func (db *ForecasterDB) CreatePoll(ctx context.Context, poll fcasterbot.Poll) (fcasterbot.Poll, error) {
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return fcasterbot.Poll{}, unableToStartTransaction(err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	pollSQL, args, err := db.q.
		Insert("forecaster.polls").
		Columns("title", "description", "start", "finish", "created_at", "updated_at").
		Values(poll.Title, poll.Description, poll.Start, poll.Finish).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("insert poll", err)
	}

	var pollID int32

	err = tx.QueryRow(ctx, pollSQL, args...).Scan(&pollID)
	if err != nil {
		return fcasterbot.Poll{}, scanFailed("insert poll", err)
	}

	poll.ID = pollID

	q := db.q.Insert("forecaster.options").
		Columns("title", "description", "poll_id")

	for _, option := range poll.Options {
		q = q.Values(option.Title, option.Description, pollID)
	}

	optionSQL, args, err := q.Suffix("RETURNING id").ToSql()
	if err != nil {
		return fcasterbot.Poll{}, buildingQueryFailed("insert options", err)
	}

	rows, err := tx.Query(ctx, optionSQL, args...)
	if err != nil {
		return fcasterbot.Poll{}, queryFailed("insert options", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option fcasterbot.Option

		err = rows.Scan(&option.ID)
		if err != nil {
			return fcasterbot.Poll{}, scanFailed("insert options", err)
		}

		poll.Options = append(poll.Options, option)
	}

	if err = rows.Err(); err != nil {
		return fcasterbot.Poll{}, rowsError("insert options", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fcasterbot.Poll{}, unableToCommitTransaction(err)
	}

	return poll, nil
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
		return fcasterbot.Poll{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	var updatedAt time.Time

	err = db.db.QueryRow(ctx, pollSQL, args...).Scan(&updatedAt)
	if err != nil {
		return fcasterbot.Poll{}, scanFailed("update poll", err)
	}

	poll.UpdatedAt = updatedAt

	return poll, nil
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
	return fmt.Errorf("%s: exec query failed: %w", queryName, err)
}

func unableToStartTransaction(err error) error {
	return fmt.Errorf("unable to start transaction: %w", err)
}

func unableToCommitTransaction(err error) error {
	return fmt.Errorf("unable to commit transaction: %w", err)
}
