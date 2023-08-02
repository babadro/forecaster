package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	models "github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	db    db
	tgBot *tgbotapi.BotAPI
}

func NewService(db db, tgBot *tgbotapi.BotAPI) *Service {
	return &Service{db: db, tgBot: tgBot}
}

type db interface {
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

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.tgBot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result := s.processTGUpdate(ctx, upd)

	// todo remove
	result.msgText = upd.Message.Text

	if result.msgText != "" {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, result.msgText)
		if _, sendErr := s.tgBot.Send(msg); sendErr != nil {
			return fmt.Errorf("unable to send message: %s", sendErr.Error())
		}
	}

	return nil
}

type processTGResult struct {
	msgText string
}

func (s *Service) processTGUpdate(ctx context.Context, upd tgbotapi.Update) processTGResult {
	l := zerolog.Ctx(ctx)

	if upd.Message != nil {
		text := upd.Message.Text

		prefix := "/start showpoll"
		if strings.HasPrefix(prefix, text) {
			pollIDStr := strings.TrimPrefix(text, prefix)

			pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
			if err != nil {
				l.Error().Msgf("unable to convert poll id to int: %v\n", err)

				return processTGResult{
					msgText: fmt.Sprintf("invalid poll id: %s", pollIDStr),
				}
			}

			poll, err := s.db.GetPollByID(ctx, int32(pollID))
			_ = poll // todo remove

			if err != nil {
				l.Error().Msgf("unable to get poll by id: %v\n", err)

				return processTGResult{
					msgText: fmt.Sprintf("oops, can't find poll with id %d", pollID),
				}
			}
		}
	}

	return processTGResult{}
}
