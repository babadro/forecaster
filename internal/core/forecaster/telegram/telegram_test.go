package telegram_test

import (
	"context"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (s *TelegramServiceSuite) createRandomPolls(count int) []swagger.PollWithOptions {
	s.T().Helper()

	polls := make([]swagger.PollWithOptions, count)

	for i := range polls {
		polls[i] = s.createRandomPoll(time.Now().Add(time.Second * time.Duration(i)))
	}

	return polls
}

func (s *TelegramServiceSuite) createRandomPoll(now time.Time) swagger.PollWithOptions {
	s.T().Helper()

	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0
	pollInput.Start = strfmt.DateTime(time.Now().Add(-time.Hour))
	pollInput.Finish = strfmt.DateTime(time.Now().Add(time.Hour))

	return s.createPoll(pollInput, now)
}

func (s *TelegramServiceSuite) createPoll(pollInput swagger.CreatePoll, now time.Time) swagger.PollWithOptions {
	s.T().Helper()

	ctx := context.Background()

	poll, err := s.db.CreatePoll(ctx, pollInput, now)
	s.Require().NoError(err)

	createdOptions := make([]*swagger.Option, 3)
	for i := range createdOptions {
		optionInput := randomModel[swagger.CreateOption](s.T())
		optionInput.PollID = poll.ID
		optionInput.Title = "option " + strconv.Itoa(i+1)

		var op swagger.Option
		op, err = s.db.CreateOption(ctx, optionInput, time.Now())
		s.Require().NoError(err)

		createdOptions[i] = &op
	}

	pollWithOptions, err := s.db.GetPollByID(ctx, poll.ID)
	s.Require().NoError(err)

	return pollWithOptions
}

func (s *TelegramServiceSuite) mockTelegramSender(sentMsg *interface{}) {
	s.mockTgBot.On("Send", mock.Anything).
		Return(tgbotapi.Message{}, nil).
		Run(func(args mock.Arguments) {
			*sentMsg = args.Get(0)
		})
}

func (s *TelegramServiceSuite) sendCallback(button tgbotapi.InlineKeyboardButton, userID int64) tgbotapi.Update {
	s.T().Helper()

	s.Require().NotNil(button.CallbackData)

	update := callback(*button.CallbackData, userID)
	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().NoError(err)

	return update
}

func (s *TelegramServiceSuite) findOptionByCallbackData(
	poll swagger.PollWithOptions, callbackData *string) *swagger.Option {
	s.T().Helper()

	s.Require().NotNil(callbackData)

	decoded, err := base64.StdEncoding.DecodeString(*callbackData)
	require.NoError(s.T(), err)

	votepreview := &votepreview2.VotePreview{}
	err = proto.UnmarshalCallbackData(string(decoded), votepreview)

	s.Require().NoError(err)
	s.Require().Equal(poll.ID, *votepreview.PollId)

	for _, op := range poll.Options {
		if int32(op.ID) == *votepreview.OptionId {
			return op
		}
	}

	s.Fail("unable to find option with id %d", *votepreview.OptionId)

	return nil
}

func getButtons(keyboard tgbotapi.InlineKeyboardMarkup) []tgbotapi.InlineKeyboardButton {
	var buttons []tgbotapi.InlineKeyboardButton

	for _, row := range keyboard.InlineKeyboard {
		buttons = append(buttons, row...)
	}

	return buttons
}

func startShowPoll(pollID int32, userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: userID,
			},
			Text: "/start showpoll_" + strconv.Itoa(int(pollID)),
		},
	}
}

func startShowForecast(pollID int32, userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: userID,
			},
			Text: "/start showforecast_" + strconv.Itoa(int(pollID)),
		},
	}
}

func startShowPolls(currentPage int32, userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: userID,
			},
			Text: "/start showpolls_" + strconv.Itoa(int(currentPage)),
		},
	}
}

func startShowForecasts(userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: userID,
			},
			Text: "/start showforecasts_1",
		},
	}
}

func startShowUserRes(pollID int32, userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: userID,
			},
			Text: "/start showuserres_" + strconv.Itoa(int(pollID)) + "_" + strconv.Itoa(int(userID)),
		},
	}
}

func callback(data string, userID int64) tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Message: &tgbotapi.Message{
				MessageID: 1,                       // to pass validation
				Chat:      &tgbotapi.Chat{ID: 123}, // to pass validation

			},
			Data: data,
			From: &tgbotapi.User{
				ID:       userID,
				UserName: "user" + strconv.Itoa(int(userID)),
			},
		},
	}
}

func (s *TelegramServiceSuite) sendMessage(upd tgbotapi.Update) {
	s.T().Helper()

	err := s.telegramService.ProcessTelegramUpdate(&s.logger, upd)
	s.Require().NoError(err)
}

func (s *TelegramServiceSuite) asMessage(sentMsg interface{}) tgbotapi.MessageConfig {
	s.T().Helper()

	msg, ok := sentMsg.(tgbotapi.MessageConfig)
	s.Require().True(ok)

	return msg
}

func (s *TelegramServiceSuite) asEditMessage(sentMsg interface{}) tgbotapi.EditMessageTextConfig {
	s.T().Helper()

	msg, ok := sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	return msg
}

func (s *TelegramServiceSuite) buttonsFromInterface(in interface{}) []tgbotapi.InlineKeyboardButton {
	s.T().Helper()

	switch keyboard := in.(type) {
	case tgbotapi.InlineKeyboardMarkup:
		return getButtons(keyboard)
	case *tgbotapi.InlineKeyboardMarkup:
		return getButtons(*keyboard)
	default:
		s.Failf("can't get buttons from interface", "unexpected type %T", in)
	}

	return nil
}
