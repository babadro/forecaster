package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	models3 "github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	models2 "github.com/babadro/forecaster/internal/models"

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

const pollColumns = "id, series_id, telegram_user_id, title, description, start, finish, popularity, " +
	"created_at, updated_at"

func (db *ForecasterDB) GetPollByID(ctx context.Context, id int32) (models.Poll, error) {
	pollSQL, args, err := db.q.Select(pollColumns).From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return models.Poll{}, buildingQueryFailed("select poll", err)
	}

	var poll models.Poll
	err = db.db.
		QueryRow(ctx, pollSQL, args...).
		Scan(
			&poll.ID, &poll.SeriesID, &poll.TelegramUserID, &poll.Title, &poll.Description,
			&poll.Start, &poll.Finish, &poll.Popularity, &poll.CreatedAt, &poll.UpdatedAt,
		)

	selectPoll := "select poll"

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Poll{}, errNotFound(selectPoll, err)
		}

		return models.Poll{}, scanFailed(selectPoll, err)
	}

	return poll, nil
}

func (db *ForecasterDB) GetPollWithOptionsByID(ctx context.Context, id int32) (models.PollWithOptions, error) {
	pollSQL, args, err := db.q.Select(pollColumns).From("forecaster.polls").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return models.PollWithOptions{}, buildingQueryFailed("select poll", err)
	}

	var poll models.PollWithOptions
	err = db.db.
		QueryRow(ctx, pollSQL, args...).
		Scan(
			&poll.ID, &poll.SeriesID, &poll.TelegramUserID, &poll.Title, &poll.Description,
			&poll.Start, &poll.Finish, &poll.Popularity, &poll.CreatedAt, &poll.UpdatedAt,
		)

	selectPoll := "select poll"

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.PollWithOptions{}, errNotFound(selectPoll, err)
		}

		return models.PollWithOptions{}, scanFailed(selectPoll, err)
	}

	optionsSQL, args, err := db.q.Select(
		"id", "poll_id", "title", "description", "is_actual_outcome", "total_votes", "updated_at",
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

		err = rows.Scan(
			&option.ID, &option.PollID, &option.Title, &option.Description,
			&option.IsActualOutcome, &option.TotalVotes, &option.UpdatedAt,
		)
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

func (db *ForecasterDB) GetPolls(
	ctx context.Context, offset, limit uint64, filter models3.PollFilter, sort models3.PollSort,
) ([]models.Poll, int32, error) {
	var rowsCount sql.NullInt32

	err := db.db.
		QueryRow(ctx, "SELECT count(*) FROM forecaster.polls").
		Scan(&rowsCount)
	if err != nil {
		return nil, 0, scanFailed("select count(*) from forecaster.polls", err)
	}

	if rowsCount.Int32 == 0 {
		return nil, 0, nil
	}

	orderBy, err := pollOrderBy(sort)
	if err != nil {
		return nil, 0, err
	}

	b := db.q.
		Select(
			"id", "series_id", "telegram_user_id", "title", "description", "start", "finish", "popularity",
			"created_at", "updated_at",
		).
		From("forecaster.polls").OrderBy(orderBy).
		Limit(limit).Offset(offset)

	if filter.TelegramUserID.Defined {
		b = b.Where(sq.Eq{"telegram_user_id": filter.TelegramUserID.Value})
	}

	pollsSQL, args, err := b.ToSql()

	if err != nil {
		return nil, 0, buildingQueryFailed("select polls", err)
	}

	rows, err := db.db.Query(ctx, pollsSQL, args...)
	if err != nil {
		return nil, 0, queryFailed("select polls", err)
	}

	defer rows.Close()

	polls := make([]models.Poll, 0, limit)

	for rows.Next() {
		var poll models.Poll

		err = rows.Scan(
			&poll.ID, &poll.SeriesID, &poll.TelegramUserID, &poll.Title, &poll.Description,
			&poll.Start, &poll.Finish, &poll.Popularity, &poll.CreatedAt, &poll.UpdatedAt,
		)
		if err != nil {
			return nil, 0, scanFailed("select polls", err)
		}

		polls = append(polls, poll)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, rowsError("select polls", err)
	}

	return polls, rowsCount.Int32, nil
}

func pollOrderBy(sort models3.PollSort) (string, error) {
	switch sort.By {
	case models3.PopularityPollSort:
		if sort.Asc {
			return "popularity ASC", nil
		}

		return "popularity DESC", nil
	case models3.CreatedAtPollSort, models3.DefaultPollSort:
		if sort.Asc {
			return "created_at ASC", nil
		}

		return "created_at DESC", nil
	default:
		return "", fmt.Errorf("unknown sort type: %d", sort.By)
	}
}

func (db *ForecasterDB) GetForecasts(
	ctx context.Context, offset, limit uint64,
) ([]models2.Forecast, int32, error) {
	var rowsCount sql.NullInt32

	err := db.db.
		QueryRow(ctx, `
			SELECT count(distinct p.id) FROM forecaster.polls p
			    INNER JOIN forecaster.options o ON p.id = o.poll_id
			    WHERE o.total_votes > 0`).
		Scan(&rowsCount)
	if err != nil {
		return nil, 0, scanFailed("select count(*) from forecaster.polls", err)
	}

	if rowsCount.Int32 == 0 {
		return nil, 0, nil
	}

	forecastsSQL, args, err := db.q.
		Select("o.poll_id, p.title, popularity").
		From("forecaster.options o").
		Join("forecaster.polls p ON o.poll_id = p.id").
		Where("o.total_votes > 0").
		Distinct().
		OrderBy("p.popularity DESC").
		Offset(offset).
		Limit(limit).
		ToSql()

	if err != nil {
		return nil, 0, buildingQueryFailed("select forecasts", err)
	}

	rows, err := db.db.Query(ctx, forecastsSQL, args...)
	if err != nil {
		return nil, 0, queryFailed("select forecasts", err)
	}

	defer rows.Close()

	forecasts := make([]models2.Forecast, 0, limit)

	for rows.Next() {
		var forecast models2.Forecast

		err = rows.Scan(&forecast.PollID, &forecast.PollTitle, nil)
		if err != nil {
			return nil, 0, scanFailed("select forecasts", err)
		}

		forecasts = append(forecasts, forecast)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, rowsError("select forecasts", err)
	}

	pollIDs := make([]int32, len(forecasts))
	for i := range forecasts {
		pollIDs[i] = forecasts[i].PollID
	}

	optionsSQL, args, err := db.q.Select("poll_id", "id", "title", "total_votes").
		From("forecaster.options").Where(sq.Eq{"poll_id": pollIDs}).ToSql()
	if err != nil {
		return nil, 0, buildingQueryFailed("select forecast options", err)
	}

	rows, err = db.db.Query(ctx, optionsSQL, args...)
	if err != nil {
		return nil, 0, queryFailed("select forecast options", err)
	}

	defer rows.Close()

	pollIDToForecastIDx := make(map[int32]int, len(forecasts))
	for i := range forecasts {
		pollIDToForecastIDx[forecasts[i].PollID] = i
	}

	for rows.Next() {
		pollID, option := int32(0), models2.ForecastOption{}

		err = rows.Scan(&pollID, &option.ID, &option.Title, &option.TotalVotes)
		if err != nil {
			return nil, 0, scanFailed("select forecast options", err)
		}

		forecastIDx := pollIDToForecastIDx[pollID]
		forecasts[forecastIDx].Options = append(forecasts[forecastIDx].Options, option)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, rowsError("select forecast options", err)
	}

	return forecasts, rowsCount.Int32, nil
}

func (db *ForecasterDB) CreateSeries(ctx context.Context, s models.CreateSeries, now time.Time) (models.Series, error) {
	seriesSQL, args, err := db.q.
		Insert("forecaster.series").Columns("title", "description", "updated_at", "created_at").
		Values(s.Title, s.Description, now, now).
		Suffix("RETURNING id, title, description, created_at, updated_at").
		ToSql()

	if err != nil {
		return models.Series{}, buildingQueryFailed("insert series", err)
	}

	var res models.Series

	err = db.db.QueryRow(ctx, seriesSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.UpdatedAt, &res.CreatedAt)
	if err != nil {
		return models.Series{}, scanFailed("insert series", err)
	}

	return res, nil
}

func (db *ForecasterDB) CreatePoll(ctx context.Context, poll models.CreatePoll, now time.Time) (models.Poll, error) {
	pollSQL, args, err := db.q.
		Insert("forecaster.polls").
		Columns("series_id", "telegram_user_id", "title", "description", "start", "finish", "status",
			"created_at", "updated_at").
		Values(poll.SeriesID, poll.TelegramUserID, poll.Title, poll.Description, poll.Start, poll.Finish, "draft",
			now, now).
		Suffix("RETURNING id, series_id, telegram_user_id, title, " +
			"description, start, finish, popularity, created_at, updated_at").
		ToSql()

	if err != nil {
		return models.Poll{}, buildingQueryFailed("insert poll", err)
	}

	var res models.Poll

	err = db.db.QueryRow(ctx, pollSQL, args...).
		Scan(&res.ID, &res.SeriesID, &res.TelegramUserID, &res.Title,
			&res.Description, &res.Start, &res.Finish, &res.Popularity, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return models.Poll{}, scanFailed("insert poll", err)
	}

	return res, nil
}

func (db *ForecasterDB) CreateOption(
	ctx context.Context, option models.CreateOption, now time.Time,
) (models.Option, error) {
	var rowsCount sql.NullInt32

	err := db.db.
		QueryRow(ctx, "SELECT count(*) FROM forecaster.polls WHERE id = $1", option.PollID).
		Scan(&rowsCount)
	if err != nil {
		return models.Option{}, scanFailed("select count(*) from forecaster.polls", err)
	}

	if rowsCount.Int32 == 0 {
		return models.Option{}, fmt.Errorf("%w: poll with id %d does not exist", domain.ErrNotFound, option.PollID)
	}

	query, args, err := db.q.Select("MAX(id)").
		From("forecaster.options").
		Where(sq.Eq{"poll_id": option.PollID}).
		ToSql()

	if err != nil {
		return models.Option{}, buildingQueryFailed("select max id", err)
	}

	var maxOptionID sql.NullInt16
	if err = db.db.QueryRow(ctx, query, args...).Scan(&maxOptionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Option{}, errNotFound("select max id", err)
		}

		return models.Option{}, scanFailed("select max id", err)
	}

	optionID := maxOptionID.Int16 + 1

	optionSQL, args, err := db.q.
		Insert("forecaster.options").
		Columns("id", "poll_id", "title", "description", "updated_at").
		Values(optionID, option.PollID, option.Title, option.Description, now).
		Suffix("RETURNING id, poll_id, title, description, is_actual_outcome, updated_at").
		ToSql()

	if err != nil {
		return models.Option{}, buildingQueryFailed("insert option", err)
	}

	var res models.Option

	err = db.db.QueryRow(ctx, optionSQL, args...).
		Scan(&res.ID, &res.PollID, &res.Title, &res.Description, &res.IsActualOutcome, &res.UpdatedAt)
	if err != nil {
		return models.Option{}, scanFailed("insert option", err)
	}

	return res, nil
}

func (db *ForecasterDB) UpdateSeries(
	ctx context.Context, id int32, s models.UpdateSeries, now time.Time,
) (models.Series, error) {
	b := db.q.Update("forecaster.series").
		Set("updated_at", now).
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

	var res models.Series

	err = db.db.QueryRow(ctx, seriesSQL, args...).
		Scan(&res.ID, &res.Title, &res.Description, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return models.Series{}, scanFailed("update series", err)
	}

	return res, err
}

func (db *ForecasterDB) UpdatePoll(
	ctx context.Context, id int32, in models.UpdatePoll, now time.Time,
) (models.Poll, error) {
	b := db.q.Update("forecaster.polls").
		Set("updated_at", now).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, series_id, telegram_user_id, title, description, " +
			"start, finish, popularity, updated_at, created_at")

	if in.SeriesID != nil {
		b = b.Set("series_id", in.SeriesID)
	}

	if in.TelegramUserID != nil {
		b = b.Set("telegram_user_id", in.TelegramUserID)
	}

	if in.Title != nil {
		b = b.Set("title", in.Title)
	}

	if in.Description != nil {
		b = b.Set("description", in.Description)
	}

	if in.Start != nil {
		b = b.Set("start", in.Start)
	}

	if in.Finish != nil {
		b = b.Set("finish", in.Finish)
	}

	pollSQL, args, err := b.ToSql()

	if err != nil {
		return models.Poll{}, buildingQueryFailed("update poll", err)
	}

	var res models.Poll

	err = db.db.QueryRow(ctx, pollSQL, args...).
		Scan(
			&res.ID, &res.SeriesID, &res.TelegramUserID, &res.Title,
			&res.Description, &res.Start, &res.Finish, &res.Popularity, &res.UpdatedAt, &res.CreatedAt,
		)
	if err != nil {
		return models.Poll{}, scanFailed("update poll", err)
	}

	return res, nil
}

func (db *ForecasterDB) UpdateOption(
	ctx context.Context, pollID int32, optionID int16, in models.UpdateOption, now time.Time,
) (models.Option, error) {
	// Build the query
	b := db.q.Update("forecaster.options").
		Set("updated_at", now).
		Where(sq.Eq{"poll_id": pollID}).
		Where(sq.Eq{"id": optionID}).
		Suffix("RETURNING id, poll_id, title, description, is_actual_outcome, total_votes, updated_at")

	if in.Title != nil {
		b = b.Set("title", in.Title)
	}

	if in.Description != nil {
		b = b.Set("description", in.Description)
	}

	if in.IsActualOutcome != nil {
		b = b.Set("is_actual_outcome", in.IsActualOutcome)
	}

	optionSQL, args, err := b.ToSql()
	if err != nil {
		return models.Option{}, buildingQueryFailed("update option", err)
	}

	var tx pgx.Tx

	// Start a new transaction only if IsActualOutcome is being set to true
	if in.IsActualOutcome != nil && *in.IsActualOutcome {
		tx, err = db.db.Begin(ctx)
		if err != nil {
			return models.Option{}, fmt.Errorf("unable to start transaction: %w", err)
		}

		var existingOptionID int16

		err = tx.QueryRow(ctx,
			"SELECT id FROM forecaster.options WHERE poll_id = $1 AND is_actual_outcome = TRUE AND id != $2",
			pollID, optionID,
		).Scan(&existingOptionID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return models.Option{}, rollback(ctx, tx, scanFailed("searching existing actual outcome for poll", err))
		}

		if err == nil {
			return models.Option{}, rollback(ctx, tx, domain.OptionWithOutcomeFlagAlreadyExistsError{
				PollID:   pollID,
				OptionID: existingOptionID,
			})
		}
	}

	var res models.Option
	if tx == nil {
		err = db.db.QueryRow(ctx, optionSQL, args...).
			Scan(&res.ID, &res.PollID, &res.Title, &res.Description, &res.IsActualOutcome, &res.TotalVotes, &res.UpdatedAt)
		if err != nil {
			return models.Option{}, scanFailed("update option", err)
		}

		return res, nil
	}

	err = tx.QueryRow(ctx, optionSQL, args...).
		Scan(&res.ID, &res.PollID, &res.Title, &res.Description, &res.IsActualOutcome, &res.TotalVotes, &res.UpdatedAt)
	if err != nil {
		return models.Option{}, rollback(ctx, tx, scanFailed("update option", err))
	}

	if err = tx.Commit(ctx); err != nil {
		return models.Option{}, commitTxFailed("update option", err)
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

func (db *ForecasterDB) DeleteOption(ctx context.Context, pollID int32, optionID int16) error {
	optionSQL, args, err := db.q.
		Delete("forecaster.options").
		Where(sq.Eq{"poll_id": pollID}).
		Where(sq.Eq{"id": optionID}).
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

func (db *ForecasterDB) CreateVote(
	ctx context.Context, vote models.CreateVote, nowUnixTimestamp int64,
) (models.Vote, error) {
	voteSQL, args, err := db.q.
		Insert("forecaster.votes").
		Columns("poll_id", "option_id", "user_id", "epoch_unix_timestamp").
		Values(vote.PollID, vote.OptionID, vote.UserID, nowUnixTimestamp).
		Suffix(`ON CONFLICT (poll_id, user_id) DO UPDATE 
					SET option_id = EXCLUDED.option_id, epoch_unix_timestamp = EXCLUDED.epoch_unix_timestamp
					WHERE forecaster.votes.option_id != EXCLUDED.option_id
					RETURNING poll_id, option_id, user_id, epoch_unix_timestamp`).
		ToSql()

	if err != nil {
		return models.Vote{}, buildingQueryFailed("insert vote", err)
	}

	var res models.Vote

	err = db.db.QueryRow(ctx, voteSQL, args...).
		Scan(&res.PollID, &res.OptionID, &res.UserID, &res.EpochUnixTimestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Vote{}, domainError("insert or update vote", domain.ErrVoteWithSameOptionAlreadyExists, err)
		}

		return models.Vote{}, scanFailed("insert or update vote", err)
	}

	return res, nil
}

func (db *ForecasterDB) GetUserVote(ctx context.Context, userID int64, pollID int32) (models.Vote, error) {
	voteSQL, args, err := db.q.
		Select("poll_id", "option_id", "user_id", "position", "epoch_unix_timestamp").
		From("forecaster.votes").
		Where(sq.Eq{"poll_id": pollID}).
		Where(sq.Eq{"user_id": userID}).
		OrderBy("epoch_unix_timestamp DESC").
		Limit(1).
		ToSql()

	if err != nil {
		return models.Vote{}, buildingQueryFailed("select vote", err)
	}

	var res models.Vote

	err = db.db.QueryRow(ctx, voteSQL, args...).
		Scan(&res.PollID, &res.OptionID, &res.UserID, &res.Position, &res.EpochUnixTimestamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Vote{}, errNotFound("select vote", err)
		}

		return models.Vote{}, scanFailed("select vote", err)
	}

	return res, nil
}

const updateVotePositionsQuery = `
WITH numbered_votes AS
(
    SELECT user_id, ROW_NUMBER() OVER (ORDER BY epoch_unix_timestamp ASC) AS rn
    FROM forecaster.votes
    WHERE poll_id = $1 AND epoch_unix_timestamp >= $2 AND epoch_unix_timestamp <= $3
	AND option_id = (
        SELECT id FROM forecaster.options
        WHERE poll_id = $1 AND is_actual_outcome = true
    )
	
)
UPDATE forecaster.votes
SET position = numbered_votes.rn
FROM numbered_votes
WHERE forecaster.votes.poll_id = $1 AND forecaster.votes.user_id = numbered_votes.user_id;
`

func (db *ForecasterDB) CalculateStatistics(ctx context.Context, pollID int32) error {
	var start, finish time.Time

	selectStartFinish := "select start, finish from forecaster.polls"

	if err := db.db.QueryRow(ctx,
		"select start, finish from forecaster.polls where polls.id = $1", pollID).
		Scan(&start, &finish); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errNotFound(selectStartFinish, err)
		}

		return scanFailed(selectStartFinish, err)
	}

	startUnix, finishUnix := start.Unix(), finish.Unix()
	// Update total_votes for each option
	updateTotalVotesSQL, args, err := db.q.
		Update("forecaster.options").
		Set("total_votes",
			sq.Expr(`(
					SELECT COUNT(*) FROM forecaster.votes
						WHERE poll_id = ? AND option_id = forecaster.options.id
						AND epoch_unix_timestamp >= ? AND epoch_unix_timestamp <= ?
					)`,
				pollID, startUnix, finishUnix,
			),
		).
		Where(sq.Eq{"poll_id": pollID}).
		ToSql()

	if err != nil {
		return buildingQueryFailed("update total_votes", err)
	}

	_, err = db.db.Exec(ctx, updateTotalVotesSQL, args...)
	if err != nil {
		return execFailed("update total_votes", err)
	}

	if _, err = db.db.Exec(ctx, updateVotePositionsQuery, pollID, startUnix, finishUnix); err != nil {
		return execFailed("update position", err)
	}

	return nil
}

func buildingQueryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: building query failed: %s", queryName, err.Error())
}

func queryFailed(queryName string, err error) error {
	return fmt.Errorf("%s: query failed: %s", queryName, err.Error())
}

func rowsError(queryName string, err error) error {
	return fmt.Errorf("%s: rows error: %s", queryName, err.Error())
}

func scanFailed(queryName string, err error) error {
	return fmt.Errorf("%s: scan rows failed: %s", queryName, err.Error())
}

func execFailed(queryName string, err error) error {
	return fmt.Errorf("%s: exec failed: %s", queryName, err.Error())
}

func errNotFound(queryName string, err error) error {
	return fmt.Errorf("%s: %w: %s", queryName, domain.ErrNotFound, err.Error())
}

func domainError(queryName string, domainErr, dbErr error) error {
	return fmt.Errorf("%s: %w: %s", queryName, domainErr, dbErr.Error())
}

func rollback(ctx context.Context, tx pgx.Tx, err error) error {
	tErr := tx.Rollback(ctx)
	if tErr != nil {
		return fmt.Errorf("forecasterDB: rollback failed: %s. Original error: %w", tErr.Error(), err)
	}

	return err
}

func commitTxFailed(queryName string, err error) error {
	return fmt.Errorf("%s: commit transaction failed: %s", queryName, err.Error())
}
