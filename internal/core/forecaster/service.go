package forecaster

import (
	"context"
	"fmt"
	"sync"

	models "github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	db    DB
	tgBot *tgbotapi.BotAPI
}

func NewService(db DB, tgBot *tgbotapi.BotAPI) *Service {
	return &Service{db: db, tgBot: tgBot}
}

type Tg struct {
	tgBot *tgbotapi.BotAPI
	wg    *sync.WaitGroup
}

type DB interface {
	GetSeriesByID(ctx context.Context, id int32) (models.Series, error)
	GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error)

	CreateSeries(ctx context.Context, s models.CreateSeries) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error)

	UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries) (models.Series, error)
	UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll) (models.Poll, error)
	UpdateOption(ctx context.Context, id int32, option models.UpdateOption) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, id int32) error
}

func (s *Service) GetSeriesByID(ctx context.Context, id int32) (models.Series, error) {
	return s.db.GetSeriesByID(ctx, id)
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

func (s *Service) UpdateSeries(ctx context.Context, id int32, series models.UpdateSeries) (models.Series, error) {
	return s.db.UpdateSeries(ctx, id, series)
}

func (s *Service) UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll) (models.Poll, error) {
	return s.db.UpdatePoll(ctx, id, poll)
}

func (s *Service) UpdateOption(ctx context.Context, id int32, option models.UpdateOption) (models.Option, error) {
	return s.db.UpdateOption(ctx, id, option)
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

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.tgBot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
	if _, sendErr := s.tgBot.Send(msg); sendErr != nil {
		return fmt.Errorf("Unable to send message: %v\n", sendErr)
	}

	return nil
}
