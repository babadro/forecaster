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
}

func (s *Service) GetPollByID(ctx context.Context, id int32) (bot.Poll, error) {
	return s.db.GetPollByID(ctx, id)
}
