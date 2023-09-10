package telegram

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/errorpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/poll"
	userpollresult "github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/userpoll_result"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/vote"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type pageServices struct {
	votePreview    *votepreview.Service
	vote           *vote.Service
	poll           *poll.Service
	userPollResult *userpollresult.Service
}

type Service struct {
	db  models.DB
	bot models.TgBot

	pages            pageServices
	callbackHandlers [256]handlerFunc
}

func NewService(db models.DB, b models.TgBot, botName string) *Service {
	pages := pageServices{
		votePreview:    votepreview.New(db),
		vote:           vote.New(db),
		poll:           poll.New(db),
		userPollResult: userpollresult.New(db, botName),
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

		result = errorpage.ErrorPage(logger, errMsg, upd)
	}

	var sendErr error

	if result != nil {
		if _, err := s.bot.Send(result); err != nil {
			sendErr = fmt.Errorf("unable to send message: %s", err.Error())
		}
	}

	if switcherErr != nil && sendErr != nil {
		return fmt.Errorf("switcher error: %s; send error: %s", switcherErr.Error(), sendErr.Error())
	}

	if switcherErr != nil {
		return switcherErr
	}

	return sendErr
}

// update type int8 iota
const (
	unknownUpdateType byte = iota
	showPollStartCommandUpdateType
	renderCallbackUpdateType
)

func (s *Service) switcher(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	var msg tgbotapi.Chattable

	var errMsg string
	var updateType, route byte
	var err error

	switch {
	case upd.Message != nil:
		if strings.HasPrefix(upd.Message.Text, models.ShowPollStartCommandPrefix) {
			updateType = showPollStartCommandUpdateType
			msg, errMsg, err = s.pages.poll.RenderStartCommand(ctx, upd)
		}
	case upd.CallbackData() != "":
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(upd.CallbackQuery.Data)
		if err != nil {
			return nil, "", fmt.Errorf("decode error: %s", err.Error())
		}

		route = decoded[0]
		upd.CallbackQuery.Data = string(decoded)

		updateType = renderCallbackUpdateType

		msg, errMsg, err = s.callbackHandlers[route](ctx, upd)
	}

	if updateType != unknownUpdateType {
		if err != nil {
			return nil, errMsg, unableToHandleUpdate(updateType, route, err)
		}

		return msg, errMsg, nil
	}

	return nil, "", errors.New("unknown update type")
}

func unableToHandleUpdate(updateType, route byte, err error) error {
	if updateType == renderCallbackUpdateType {
		return fmt.Errorf("unable to handle callback with route %d: %w", route, err)
	}

	uType := "unknown update type"
	if updateType == showPollStartCommandUpdateType {
		uType = "show poll start command update type"
	}

	return fmt.Errorf("unable to handle %s: %w", uType, err)
}
