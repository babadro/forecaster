package forecaster

import (
	"context"
	models "github.com/babadro/forecaster/internal/models/swagger"
)

type Service struct {
	db DB
}

func NewService(db DB) *Service {
	return &Service{db: db}
}

type DB interface {
	GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error)

	CreateSeries(ctx context.Context, s models.CreateSeries) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error)

	UpdateSeries(ctx context.Context, s models.UpdateSeries) (models.Series, error)
	UpdatePoll(ctx context.Context, poll models.UpdatePoll) (models.Poll, error)
	UpdateOption(ctx context.Context, option models.UpdateOption) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, id int32) error
}

func (s *Service) GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error) {
	return s.db.GetPollByID(ctx, id)
}

func (s *Service) CreateSeries(ctx context.Context, series models.CreateSeries) (models.Series, error) {
	return s.db.CreateSeries(ctx, series)
}

func (s *Service) CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error) {
	return s.db.CreatePoll(ctx, poll)
}

func (s *Service) CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error) {
	return s.db.CreateOption(ctx, option)
}

func (s *Service) UpdateSeries(ctx context.Context, series models.UpdateSeries) (models.Series, error) {
	return s.db.UpdateSeries(ctx, series)
}

func (s *Service) UpdatePoll(ctx context.Context, poll models.UpdatePoll) (models.Poll, error) {
	return s.db.UpdatePoll(ctx, poll)
}

func (s *Service) UpdateOption(ctx context.Context, option models.UpdateOption) (models.Option, error) {
	return s.db.UpdateOption(ctx, option)
}

func (s *Service) DeleteSeries(ctx context.Context, id int32) error {
	return s.db.DeleteSeries(ctx, id)
}

func (s *Service) DeletePoll(ctx context.Context, id int32) error {
	return s.db.DeletePoll(ctx, id)
}

func (s *Service) DeleteOption(ctx context.Context, id int32) error {
	return s.db.DeleteOption(ctx, id)
}
