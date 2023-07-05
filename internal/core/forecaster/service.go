package forecaster

import (
	"context"

	bot "github.com/babadro/forecaster/pkg/fcasterbot"
)

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db: db}
}

type DB interface {
	GetPollByID(ctx context.Context, id int32) (bot.Poll, error)
	CreatePoll(ctx context.Context, poll bot.Poll) (bot.Poll, error)
	CreateOption(ctx context.Context, option bot.Option) (bot.Option, error)
	UpdatePoll(ctx context.Context, poll bot.Poll) (bot.Poll, error)
	UpdateOption(ctx context.Context, option bot.Option) (bot.Option, error)
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, id int32) error
}

func (s *Service) GetPollByID(ctx context.Context, id int32) (bot.Poll, error) {
	return s.db.GetPollByID(ctx, id)
}
