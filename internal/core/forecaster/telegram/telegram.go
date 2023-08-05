package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	db  db
	bot tgBot
}

func NewService(db db, b tgBot) *Service {
	return &Service{db: db, bot: b}
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

type tgBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.bot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result := s.processTelegramUpdate(ctx, upd)

	if result.msgText != "" {
		logger.Info().Msg(result.msgText)

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, result.msgText)
		msg.ParseMode = "HTML"

		if _, sendErr := s.bot.Send(msg); sendErr != nil {
			return fmt.Errorf("unable to send message: %s", sendErr.Error())
		}
	}

	return nil
}

type processTGResult struct {
	msgText        string
	inlineKeyboard tgbotapi.InlineKeyboardMarkup
}

func (s *Service) processTelegramUpdate(ctx context.Context, upd tgbotapi.Update) processTGResult {
	if upd.Message != nil {
		text := upd.Message.Text

		prefix := "/start showpoll_"
		if strings.HasPrefix(text, prefix) {
			pollIDStr := text[len(prefix):]

			return s.poll(ctx, pollIDStr)
		}
	}

	return processTGResult{
		msgText: "I don't understand you",
	}
}

func formatTime[T time.Time | strfmt.DateTime](t T) string {
	return time.Time(t).Format(time.RFC822)
}

func fPrintf(sb *strings.Builder, format string, a ...any) {
	_, _ = fmt.Fprintf(sb, format, a...)
}

func fPrint(sb *strings.Builder, a ...any) {
	_, _ = fmt.Fprint(sb, a...)
}
