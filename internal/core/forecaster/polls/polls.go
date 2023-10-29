package polls

import (
	"context"
	"time"

	models "github.com/babadro/forecaster/internal/models/swagger"
)

type Service struct {
	db db
}

func NewService(db db) *Service {
	return &Service{db: db}
}

type db interface {
	GetSeriesByID(ctx context.Context, id int32) (models.Series, error)
	GetPollWithOptionsByID(ctx context.Context, id int32) (models.PollWithOptions, error)

	CreateSeries(ctx context.Context, s models.CreateSeries, now time.Time) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll, now time.Time) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption, now time.Time) (models.Option, error)

	UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries, now time.Time) (models.Series, error)
	UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll, now time.Time) (models.Poll, error)
	UpdateOption(
		ctx context.Context, pollID int32, optionID int16, option models.UpdateOption, now time.Time,
	) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, pollID int32, optionID int16) error

	CalculateStatistics(ctx context.Context, pollID int32) error
}

func (s *Service) GetSeriesByID(ctx context.Context, id int32) (models.Series, error) {
	return s.db.GetSeriesByID(ctx, id)
}

func (s *Service) GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error) {
	return s.db.GetPollWithOptionsByID(ctx, id)
}

func (s *Service) CreateSeries(ctx context.Context, series models.CreateSeries) (models.Series, error) {
	return s.db.CreateSeries(ctx, series, time.Now())
}

func (s *Service) CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error) {
	return s.db.CreatePoll(ctx, poll, time.Now())
}

func (s *Service) CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error) {
	return s.db.CreateOption(ctx, option, time.Now())
}

func (s *Service) UpdateSeries(ctx context.Context, id int32, series models.UpdateSeries) (models.Series, error) {
	return s.db.UpdateSeries(ctx, id, series, time.Now())
}

func (s *Service) UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll) (models.Poll, error) {
	return s.db.UpdatePoll(ctx, id, poll, time.Now())
}

func (s *Service) UpdateOption(
	ctx context.Context, pollID int32, optionID int16, option models.UpdateOption,
) (models.Option, error) {
	return s.db.UpdateOption(ctx, pollID, optionID, option, time.Now())
}

func (s *Service) DeleteSeries(ctx context.Context, id int32) error {
	return s.db.DeleteSeries(ctx, id)
}

func (s *Service) DeletePoll(ctx context.Context, id int32) error {
	return s.db.DeletePoll(ctx, id)
}

func (s *Service) DeleteOption(ctx context.Context, pollID int32, optionID int16) error {
	return s.db.DeleteOption(ctx, pollID, optionID)
}

func (s *Service) CalculateStatistics(ctx context.Context, pollID int32) error {
	return s.db.CalculateStatistics(ctx, pollID)
}
