package telegram_test

import (
	"context"
	"regexp"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
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
	pollInput.Start = strfmt.DateTime(time.Now().Add(-time.Hour))
	pollInput.Finish = strfmt.DateTime(time.Now().Add(time.Hour))

	poll, err := s.db.CreatePoll(ctx, pollInput, time.Now())
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

	// verify the poll message
	msg, ok := sentMsg.(tgbotapi.MessageConfig)
	s.Require().True(ok)

	pollKeyboard, ok := msg.ReplyMarkup.(tgbotapi.InlineKeyboardMarkup)
	s.Require().True(ok)

	// each keyboard button is a poll option
	pollButtons := getButtons(pollKeyboard)
	s.Require().Len(pollButtons, len(poll.Options))

	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton)

	// verify the result votepreview message
	editMsg, ok := sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	// verify contains the poll title and description
	txt := editMsg.Text
	option := s.findOptionByCallbackData(poll, firstButton.CallbackData)
	s.Require().Contains(txt, option.Title)
	s.Require().Contains(txt, option.Description)

	// verify message has two buttons
	votePreviewKeyboard := *editMsg.ReplyMarkup
	votePreviewButtons := getButtons(votePreviewKeyboard)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0])

	// verify the vote message
	editMsg, ok = sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	s.Require().Contains(editMsg.Text, "Success")

	// push back to poll button
	voteKeyboard := getButtons(*editMsg.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton := voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton)

	// verify the poll message
	editMsg, ok = sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	s.Require().Contains(editMsg.Text, poll.Title)
	pattern := "Last time you voted for:.+" + option.Title
	regex := regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(editMsg.Text), "expected %s to match regex %s", editMsg.Text, pattern)

	// each keyboard button is a poll option
	pollButtons = getButtons(pollKeyboard)
	s.Require().Len(pollButtons, len(poll.Options))

	// chose option I didn't vote earlier
	anotherOptionButton, found := tgbotapi.InlineKeyboardButton{}, false
	for _, button := range pollButtons {
		op := s.findOptionByCallbackData(poll, button.CallbackData)
		if op.ID != option.ID {
			anotherOptionButton, found = button, true
			break
		}
	}

	// push the button to vote for another option this time
	s.Require().True(found)
	// sleep for second to make sure vote timestamp (which used second precision) is different
	time.Sleep(time.Second)
	s.sendCallback(anotherOptionButton)

	// verify the votepreview message
	editMsg, ok = sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	// verify the poll contains title and description
	txt = editMsg.Text
	anotherOption := s.findOptionByCallbackData(poll, anotherOptionButton.CallbackData)
	s.Require().Contains(txt, anotherOption.Title)
	s.Require().Contains(txt, anotherOption.Description)

	// verify message has two buttons
	votePreviewKeyboard = *editMsg.ReplyMarkup
	votePreviewButtons = getButtons(votePreviewKeyboard)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0])

	// verify the vote message
	editMsg, ok = sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	s.Require().Contains(editMsg.Text, "Success")

	// push back to poll button
	voteKeyboard = getButtons(*editMsg.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton = voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton)

	// verify the poll message
	editMsg, ok = sentMsg.(tgbotapi.EditMessageTextConfig)
	s.Require().True(ok)

	s.Require().Contains(editMsg.Text, poll.Title)
	pattern = "Last time you voted for:.+" + anotherOption.Title
	regex = regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(editMsg.Text), "expected %s to match regex: %q", editMsg.Text, pattern)

	// each keyboard button is a poll option
	pollButtons = getButtons(pollKeyboard)
	s.Require().Len(pollButtons, len(poll.Options))
}

func (s *TelegramServiceSuite) sendCallback(button tgbotapi.InlineKeyboardButton) {
	s.T().Helper()

	s.Require().NotNil(button.CallbackData)

	update := callback(*button.CallbackData)
	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().NoError(err)
}

func (s *TelegramServiceSuite) findOptionByCallbackData(poll swagger.PollWithOptions, callbackData *string) *swagger.Option {
	s.T().Helper()

	s.Require().NotNil(callbackData)
	votepreview := &votepreview2.VotePreview{}
	err := proto.UnmarshalCallbackData(*callbackData, votepreview)
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
			From: &tgbotapi.User{
				ID: 456,
			},
		},
	}
}
