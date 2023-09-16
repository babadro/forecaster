package models

import (
	"context"
	"time"

	models "github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MaxCountInRow = 8

	Days365     = 365
	Hours24     = 24
	Seconds3600 = 3600
)

type DB interface {
	GetSeriesByID(ctx context.Context, id int32) (models.Series, error)
	GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error)
	GetUserVote(ctx context.Context, userID int64, pollID int32) (models.Vote, error)
	GetPolls(ctx context.Context, currentPage int32, pageSize int32) ([]models.Poll, int32, error)

	CreateSeries(ctx context.Context, s models.CreateSeries, now time.Time) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll, now time.Time) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption, now time.Time) (models.Option, error)
	CreateVote(ctx context.Context, vote models.CreateVote, nowUnixTimestamp int64) (models.Vote, error)

	UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries, now time.Time) (models.Series, error)
	UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll, now time.Time) (models.Poll, error)
	UpdateOption(
		ctx context.Context, pollID int32, optionID int16, option models.UpdateOption, now time.Time,
	) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, pollID int32, optionID int16) error
}

type TgBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Scope struct {
	DB  DB
	Bot TgBot
}
