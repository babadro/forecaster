package models

import (
	"context"
	"time"

	"github.com/babadro/forecaster/internal/models"

	swModels "github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MaxCountInRow = 8

	Days365     = 365
	Hours24     = 24
	Seconds3600 = 3600

	Percent100 = 100
)

type DB interface {
	GetSeriesByID(ctx context.Context, id int32) (swModels.Series, error)
	GetPollByID(ctx context.Context, id int32) (swModels.PollWithOptions, error)
	GetUserVote(ctx context.Context, userID int64, pollID int32) (swModels.Vote, error)
	GetPolls(ctx context.Context, offset, limit uint64, filter PollFilter, sort PollSort) ([]swModels.Poll, int32, error)
	GetForecasts(ctx context.Context, offset, limit uint64) ([]models.Forecast, int32, error)

	CreateSeries(ctx context.Context, s swModels.CreateSeries, now time.Time) (swModels.Series, error)
	CreatePoll(ctx context.Context, poll swModels.CreatePoll, now time.Time) (swModels.Poll, error)
	CreateOption(ctx context.Context, option swModels.CreateOption, now time.Time) (swModels.Option, error)
	CreateVote(ctx context.Context, vote swModels.CreateVote, nowUnixTimestamp int64) (swModels.Vote, error)

	UpdateSeries(ctx context.Context, id int32, s swModels.UpdateSeries, now time.Time) (swModels.Series, error)
	UpdatePoll(ctx context.Context, id int32, poll swModels.UpdatePoll, now time.Time) (swModels.Poll, error)
	UpdateOption(
		ctx context.Context, pollID int32, optionID int16, option swModels.UpdateOption, now time.Time,
	) (swModels.Option, error)

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

type EditButton[T any] struct {
	Text  string
	Field T
}

type Nullable[T any] struct {
	Value   T
	Defined bool
}

func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Value: value, Defined: true}
}

type PollFilter struct {
	TelegramUserID Nullable[int64]
}

func NewPollFilter() PollFilter {
	return PollFilter{}
}

func (f PollFilter) WithTelegramUserID(id int64) PollFilter {
	f.TelegramUserID = NewNullable(id)

	return f
}

type PollSortType byte

const (
	DefaultPollSort PollSortType = iota
	CreatedAtPollSort
	PopularityPollSort
)

type PollSort struct {
	By  PollSortType
	Asc bool
}
