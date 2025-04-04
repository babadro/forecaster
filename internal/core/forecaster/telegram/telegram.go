package telegram

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/deleteoption"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/deletepoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/editoptionfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/editpollfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/errorpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/forecast"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/forecasts"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/mainpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/mypolls"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/polls"
	userpollresult "github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/userpoll_result"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/vote"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/pages/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type pageServices struct {
	main            *mainpage.Service
	votePreview     *votepreview.Service
	vote            *vote.Service
	poll            *poll.Service
	polls           *polls.Service
	userPollResult  *userpollresult.Service
	forecasts       *forecasts.Service
	forecast        *forecast.Service
	editPoll        *editpoll.Service
	editpollfield   *editpollfield.Service
	myPolls         *mypolls.Service
	editOption      *editoption.Service
	editOptionField *editoptionfield.Service
	deletePoll      *deletepoll.Service
	deleteOption    *deleteoption.Service
}

type Service struct {
	db  models.DB
	bot models.TgBot

	pages            pageServices
	callbackHandlers [256]handlerFunc
}

func NewService(db models.DB, b models.TgBot, botName string) *Service {
	pages := pageServices{
		main:            mainpage.New(db),
		votePreview:     votepreview.New(db),
		vote:            vote.New(db),
		poll:            poll.New(db),
		userPollResult:  userpollresult.New(db, botName),
		polls:           polls.New(db),
		forecasts:       forecasts.New(db),
		forecast:        forecast.New(db),
		editPoll:        editpoll.New(db),
		editpollfield:   editpollfield.New(db),
		myPolls:         mypolls.New(db),
		editOption:      editoption.New(db),
		editOptionField: editoptionfield.New(db),
		deletePoll:      deletepoll.New(db),
		deleteOption:    deleteoption.New(db),
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

const (
	showMainStartCommandUpdateType       = "show_main_start_command_update_type"
	showPollStartCommandUpdateType       = "show_poll_start_command_update_type"
	renderCallbackUpdateType             = "render_callback_update_type"
	showUserResultStartCommandUpdateType = "show_user_result_start_command_update_type"
	showPollsStartCommandUpdateType      = "show_polls_start_command_update_type"
	showForecastsStartCommandUpdateType  = "show_forecasts_start_command_update_type"
	showForecastStartCommandUpdateType   = "show_forecast_start_command_update_type"
	editPollUpdateType                   = "edit_poll_update_type"
	editOptionUpdateType                 = "edit_option_update_type"
)

func (s *Service) switcher(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	var (
		msg                tgbotapi.Chattable
		errMsg, updateType string
		route              byte
		err                error
	)

	switch {
	case upd.Message != nil && upd.Message.ReplyToMessage != nil:
		if upd.Message.ReplyToMessage != nil {
			parentText := upd.Message.ReplyToMessage.Text

			if strings.HasPrefix(parentText, models.EditPollCommand) {
				updateType = editPollUpdateType
				msg, errMsg, err = validateCommandInput(s.pages.editPoll.RenderCommand)(ctx, upd)
			} else if strings.HasPrefix(parentText, models.EditOptionCommand) {
				updateType = editOptionUpdateType
				msg, errMsg, err = validateCommandInput(s.pages.editOption.RenderCommand)(ctx, upd)
			}
		}
	case upd.Message != nil:
		switch {
		case strings.HasPrefix(upd.Message.Text, models.ShowMainStartCommandPrefix):
			updateType = showMainStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.main.RenderStartCommand)(ctx, upd)
		case strings.HasPrefix(upd.Message.Text, models.ShowPollStartCommandPrefix):
			updateType = showPollStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.poll.RenderStartCommand)(ctx, upd)
		case strings.HasPrefix(upd.Message.Text, models.ShowUserResultCommandPrefix):
			updateType = showUserResultStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.userPollResult.RenderStartCommand)(ctx, upd)
		case strings.HasPrefix(upd.Message.Text, models.ShowPollsStartCommandPrefix):
			updateType = showPollsStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.polls.RenderStartCommand)(ctx, upd)
		case strings.HasPrefix(upd.Message.Text, models.ShowForecastsStartCommandPrefix):
			updateType = showForecastsStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.forecasts.RenderStartCommand)(ctx, upd)
		case strings.HasPrefix(upd.Message.Text, models.ShowForecastStartCommandPrefix):
			updateType = showForecastStartCommandUpdateType
			msg, errMsg, err = validateCommandInput(s.pages.forecast.RenderStartCommand)(ctx, upd)
		}
	case upd.CallbackData() != "":
		updateType = renderCallbackUpdateType

		var decoded []byte

		decoded, err = base64.StdEncoding.DecodeString(upd.CallbackQuery.Data)
		if err != nil {
			err = fmt.Errorf("can't decode base64: %s", err.Error())
			break
		}

		route = decoded[0]
		upd.CallbackQuery.Data = string(decoded)

		msg, errMsg, err = s.callbackHandlers[route](ctx, upd)
	}

	if updateType == "" {
		if upd.Message != nil {
			return nil, "I don't know this command", fmt.Errorf("unknown command %s", upd.Message.Text)
		}

		return nil, "", fmt.Errorf("unknown update type")
	}

	if err != nil {
		if updateType == renderCallbackUpdateType {
			return nil, errMsg, fmt.Errorf("unable to handle %s with route %d: %w", updateType, route, err)
		}

		return nil, errMsg, fmt.Errorf("unable to handle %s: %w", updateType, err)
	}

	return msg, errMsg, nil
}
