package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/poll"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	db  models.DB
	bot models.TgBot
}

func NewService(db models.DB, b models.TgBot) *Service {
	return &Service{db: db, bot: b}
}

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.bot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result := s.processTelegramUpdate(ctx, upd)

	if result.MsgText != "" {
		logger.Info().Msg(result.MsgText)

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, result.MsgText)
		msg.ParseMode = "HTML"

		if _, sendErr := s.bot.Send(msg); sendErr != nil {
			return fmt.Errorf("unable to send message: %s", sendErr.Error())
		}
	}

	return nil
}

func (s *Service) processTelegramUpdate(ctx context.Context, upd tgbotapi.Update) models.ProcessTgResult {
	sc := models.Scope{
		DB:  s.db,
		Bot: s.bot,
	}

	if upd.Message != nil {
		text := upd.Message.Text

		prefix := "/start showpoll_"
		if strings.HasPrefix(text, prefix) {
			pollIDStr := text[len(prefix):]

			return poll.Poll(ctx, pollIDStr, sc)
		}
	}

	return models.ProcessTgResult{}
}
