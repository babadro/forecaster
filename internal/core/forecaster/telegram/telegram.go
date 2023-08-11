package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/errorpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type pageServices struct {
	votePreview *votepreview.Service
	poll        *poll.Service
}

type Service struct {
	db  models.DB
	bot models.TgBot

	pages            pageServices
	callbackHandlers [256]callbackHandlerFunc
}

func NewService(db models.DB, b models.TgBot) *Service {
	pages := pageServices{
		votePreview: votepreview.New(db),
		poll:        poll.New(db),
	}

	callbackHandlers := newCallbackHandlers(pages)

	return &Service{db: db, bot: b, pages: pages, callbackHandlers: callbackHandlers}
}

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.bot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result, errMsg, switcherErr := s.switcher(ctx, upd)
	if switcherErr != nil {
		if errMsg == "" {
			errMsg = "Something went wrong"
		}

		result = errorpage.ErrorPage(errMsg)
	}

	var sendErr error
	if result != nil {
		if _, err := s.bot.Send(result); err != nil {
			sendErr = fmt.Errorf("unable to send message: %s", err.Error())
		}
	}

	if switcherErr != nil && sendErr != nil {
		return fmt.Errorf("process error: %s; send error: %s", switcherErr.Error(), sendErr.Error())
	}

	if switcherErr != nil {
		return switcherErr
	}

	return sendErr
}

func (s *Service) switcher(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	if upd.Message != nil {
		text := upd.Message.Text

		prefix := "/start showpoll_"
		if strings.HasPrefix(text, prefix) {
			pollIDStr := text[len(prefix):]

			return s.pages.poll.Render(ctx, pollIDStr, upd.Message.From.ID)
		}
	} else if callbackData := upd.CallbackData(); callbackData != "" {
		route := callbackData[0]

		return s.callbackHandlers[route](ctx, callbackData)
	}

	return nil, "", nil
}
