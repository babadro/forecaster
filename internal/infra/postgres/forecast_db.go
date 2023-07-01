package postgres

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/pkg/fcasterbot"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lann/builder"

	sq "github.com/Masterminds/squirrel"
)

type ForecastDB struct {
	db *pgxpool.Pool
	q  sq.StatementBuilderType
}

func NewForecastDB(db *pgxpool.Pool) *ForecastDB {
	return &ForecastDB{
		db: db,
		q:  sq.StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(sq.Dollar),
	}
}

func (db *ForecastDB) GetByID(ctx context.Context, id int32) (fcasterbot.Poll, error) {
	pollSQL, _, err := db.q.Select("*").From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fcasterbot.Poll{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	var poll fcasterbot.Poll
	err = db.db.
		QueryRow(ctx, pollSQL, id).
		Scan(&poll.ID, &poll.Title, &poll.Description, &poll.Start, &poll.Finish, &poll.UpdatedAt)

	if err != nil {
		return fcasterbot.Poll{}, fmt.Errorf("unable to get poll: %w", err)
	}

	optionsSQL, _, err := db.q.Select("*").From("forecaster.options").Where(sq.Eq{"poll_id": id}).ToSql()
	if err != nil {
		return fcasterbot.Poll{}, fmt.Errorf("unable to build SQL: %w", err)
	}

	rows, err := db.db.Query(ctx, optionsSQL)
	if err != nil {
		return fcasterbot.Poll{}, fmt.Errorf("unable to get poll options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option fcasterbot.Option

		err = rows.Scan(&option.ID, &option.Title, &option.Description, &option.PollID)
		if err != nil {
			return fcasterbot.Poll{}, fmt.Errorf("unable to scan poll option: %w", err)
		}

		poll.Options = append(poll.Options, option)
	}

	if err = rows.Err(); err != nil {
		return fcasterbot.Poll{}, fmt.Errorf("rows error: %w", err)
	}

	return poll, nil
}
