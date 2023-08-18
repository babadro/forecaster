package telegram_test

import (
	"context"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

func (s *TelegramServiceSuite) TestShowPollStartCommand() {
	ctx := context.Background()
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0

	poll, err := s.db.CreatePoll(ctx, pollInput, time.Now())
	s.Require().NoError(err)

	createdOptions := make([]*swagger.Option, 3)
	for i := range createdOptions {
		optionInput := randomModel[swagger.CreateOption](s.T())
		optionInput.PollID = poll.ID

		var op swagger.Option
		op, err = s.db.CreateOption(ctx, optionInput, time.Now())
		s.Require().NoError(err)

		createdOptions[i] = &op
	}

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: 456,
			},
			Text: "/start showpoll_" + strconv.Itoa(int(poll.ID)),
		},
	}

	s.mockTgBot.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		s.Assert().Equal(update.Message.Chat.ID, msg.ChatID)

		s.Assert().Contains(msg.Text, poll.Title)

		for _, op := range createdOptions {
			s.Assert().Contains(msg.Text, op.Title)
		}

		return true
	})).Return(tgbotapi.Message{}, nil)

	err = s.telegramService.ProcessTelegramUpdate(&s.logger, update)

	s.Require().NoError(err)
	s.mockTgBot.AssertExpectations(s.T())
}

func (s *TelegramServiceSuite) TestShowPollStartCommand_notFound() {
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: 456,
			},
			Text: "/start showpoll_999",
		},
	}

	s.mockTgBot.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		s.Assert().Equal(update.Message.Chat.ID, msg.ChatID)

		s.Assert().Contains(msg.Text, "can't find poll")

		return true
	})).Return(tgbotapi.Message{}, nil)

	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().Error(err)
	s.ErrorContains(err, "unable to get poll by id")

	s.mockTgBot.AssertExpectations(s.T())
}

func (s *TelegramServiceSuite) createRandomPoll() swagger.PollWithOptions {
	s.T().Helper()

	ctx := context.Background()
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0

	poll, err := s.db.CreatePoll(ctx, pollInput, time.Now())
	s.Require().NoError(err)

	createdOptions := make([]*swagger.Option, 3)
	for i := range createdOptions {
		optionInput := randomModel[swagger.CreateOption](s.T())
		optionInput.PollID = poll.ID

		var op swagger.Option
		op, err = s.db.CreateOption(ctx, optionInput, time.Now())
		s.Require().NoError(err)

		createdOptions[i] = &op
	}

	pollWithOptions, err := s.db.GetPollByID(ctx, poll.ID)
	s.Require().NoError(err)

	return pollWithOptions
}

func (s *TelegramServiceSuite) TestVoting() {
	var sentMsg interface{}

	s.mockTgBot.On("Send", mock.Anything).
		Return(tgbotapi.Message{}, nil).
		Run(func(args mock.Arguments) {
			sentMsg = args.Get(0)
		})

	poll := s.createRandomPoll()

	update := startShowPoll(poll.ID)

	// send /start showpoll_<poll_id> command
	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().NoError(err)

	// verify the result poll message
	msg, ok := sentMsg.(tgbotapi.MessageConfig)
	s.Require().True(ok)

	keyboard, ok := msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	s.Require().True(ok)

	// each keyboard button is a poll option
	buttons := getButtons(keyboard)
	s.Require().Len(buttons, len(poll.Options))

	// choose the first option
	firstButton := buttons[0]
	s.Require().NotNil(firstButton.CallbackData)
	votepreview := &votepreview2.VotePreview{}
	err = proto.UnmarshalCallbackData(*firstButton.CallbackData, votepreview)
	s.Require().NoError(err)

	// send the first option
	update = callback(*firstButton.CallbackData)
	err = s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().NoError(err)

	// verify the result votepreview message
	editMsg, ok := sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)
	// verify contains the poll title
	// verify has two buttons
	// push the first button (yes)
	txt := editMsg.Text
	s.Require().Contains(txt, poll.Title)

	_ = editMsg
}

func getButtons(keyboard tgbotapi.InlineKeyboardMarkup) []tgbotapi.InlineKeyboardButton {
	var buttons []tgbotapi.InlineKeyboardButton
	for _, row := range keyboard.InlineKeyboard {
		for _, button := range row {
			buttons = append(buttons, button)
		}
	}

	return buttons
}

func startShowPoll(pollID int32) tgbotapi.Update {
	return tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			From: &tgbotapi.User{
				ID: 456,
			},
			Text: "/start showpoll_" + strconv.Itoa(int(pollID)),
		},
	}
}

func callback(data string) tgbotapi.Update {
	return tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			Message: &tgbotapi.Message{
				MessageID: 1,                       // to pass validation
				Chat:      &tgbotapi.Chat{ID: 123}, // to pass validation
			},
			Data: data,
		},
	}
}
