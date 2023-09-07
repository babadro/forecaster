package telegram_test

import (
	"context"
	"encoding/base64"
	"regexp"
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

func (s *TelegramServiceSuite) TestShowPollStartCommand() {
	poll := s.createRandomPoll()

	update := startShowPoll(poll.ID)

	s.mockTgBot.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		s.Assert().Equal(update.Message.Chat.ID, msg.ChatID)

		s.Assert().Contains(msg.Text, poll.Title)

		for _, op := range poll.Options {
			s.Assert().Contains(msg.Text, op.Title)
		}

		return true
	})).Return(tgbotapi.Message{}, nil)

	s.sendMessage(update)

	s.mockTgBot.AssertExpectations(s.T())
}

func (s *TelegramServiceSuite) TestShowPollStartCommand_notFound() {
	update := startShowPoll(999)

	var msg tgbotapi.MessageConfig

	s.mockTgBot.On("Send", mock.MatchedBy(func(sentMsg tgbotapi.MessageConfig) bool {
		s.Assert().Equal(update.Message.Chat.ID, sentMsg.ChatID)

		s.Assert().Contains(sentMsg.Text, "can't find poll")

		msg = sentMsg

		return true
	})).Return(tgbotapi.Message{}, nil)

	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().Error(err)
	s.ErrorContains(err, "unable to get poll by id")

	buttons := s.buttonsFromInterface(msg.ReplyMarkup)
	s.Require().Len(buttons, 1)
	s.Require().Contains(buttons[0].Text, "Back to main")

	s.mockTgBot.AssertExpectations(s.T())
}

func (s *TelegramServiceSuite) createRandomPoll() swagger.PollWithOptions {
	s.T().Helper()

	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0
	pollInput.Start = strfmt.DateTime(time.Now().Add(-time.Hour))
	pollInput.Finish = strfmt.DateTime(time.Now().Add(time.Hour))

	return s.createPoll(pollInput)
}

func (s *TelegramServiceSuite) createPoll(pollInput swagger.CreatePoll) swagger.PollWithOptions {
	s.T().Helper()

	ctx := context.Background()

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

// Vote 2 times for different options of the same poll and verify the results...
func (s *TelegramServiceSuite) TestVoting() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll()

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	// each keyboard button is a poll option
	s.Require().Len(pollButtons, len(poll.Options))

	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton)

	// verify the result votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)

	// verify contains the poll title and description
	txt := votePreviewMsg.Text
	option := s.findOptionByCallbackData(poll, firstButton.CallbackData)
	s.Require().Contains(txt, option.Title)
	s.Require().Contains(txt, option.Description)

	// verify message has two buttons
	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0])

	// verify the vote message
	voteMsg := s.asEditMessage(sentMsg)

	s.Require().Contains(voteMsg.Text, "Success")

	// push back to poll button
	voteKeyboard := getButtons(*voteMsg.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton := voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton)

	// verify the poll message
	pollMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg2.Text, poll.Title)
	pattern := "Last time you voted for:.+" + option.Title
	regex := regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(pollMsg2.Text), "expected %s to match regex %s", pollMsg2.Text, pattern)

	// each keyboard button is a poll option
	pollButtons2 := getButtons(*pollMsg2.ReplyMarkup)
	s.Require().Len(pollButtons2, len(poll.Options))

	// chose option I didn't vote earlier
	anotherOptionButton, found := tgbotapi.InlineKeyboardButton{}, false

	for _, button := range pollButtons2 {
		op := s.findOptionByCallbackData(poll, button.CallbackData)
		if op.ID != option.ID {
			anotherOptionButton, found = button, true
			break
		}
	}

	s.Require().True(found)

	// sleep for second to make sure vote timestamp (which used second precision) is different
	time.Sleep(time.Second)
	// push the button to vote for another option this time
	s.sendCallback(anotherOptionButton)

	// verify the votepreview message
	votePreviewMsg2 := s.asEditMessage(sentMsg)

	// verify the poll contains title and description
	txt = votePreviewMsg2.Text
	anotherOption := s.findOptionByCallbackData(poll, anotherOptionButton.CallbackData)
	s.Require().Contains(txt, anotherOption.Title)
	s.Require().Contains(txt, anotherOption.Description)

	// verify message has two buttons
	votePreviewButtons = getButtons(*votePreviewMsg2.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the first button (yes)
	s.sendCallback(votePreviewButtons[0])

	// verify the vote message
	voteMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(voteMsg2.Text, "Success")

	// push back to poll button
	voteKeyboard = getButtons(*voteMsg2.ReplyMarkup)
	s.Require().Len(voteKeyboard, 1)

	backButton = voteKeyboard[0]
	s.Contains(backButton.Text, "Back")
	s.sendCallback(backButton)

	// verify the poll message
	pollMsg3 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg3.Text, poll.Title)
	pattern = "Last time you voted for:.+" + anotherOption.Title
	regex = regexp.MustCompile(pattern)
	s.Require().True(regex.MatchString(pollMsg3.Text), "expected %s to match regex: %q", pollMsg3.Text, pattern)

	// each keyboard button is a poll option
	pollButtons3 := getButtons(*pollMsg3.ReplyMarkup)
	s.Require().Len(pollButtons3, len(poll.Options))
}

func (s *TelegramServiceSuite) TestVotePreview_BackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll()

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)

	// send the first option
	firstButton := pollButtons[0]
	s.sendCallback(firstButton)

	// verify the result votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)

	// verify message has two buttons
	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 2)

	// push the back button
	s.sendCallback(votePreviewButtons[1])

	// verify the poll message
	pollMsg2 := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMsg2.Text, poll.Title)
}

func (s *TelegramServiceSuite) Test_expiredPoll() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0
	pollInput.Finish = strfmt.DateTime(time.Now().Add(-time.Hour)) // expired

	poll := s.createPoll(pollInput)

	// send /start showpoll_<poll_id> command
	update := startShowPoll(poll.ID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)
	// verify the poll message
	s.Require().Contains(pollMsg.Text, "poll has expired")

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	// send the first option
	s.sendCallback(pollButtons[0])

	// verify votepreview message
	votePreviewMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(votePreviewMsg.Text, "poll is expired")

	votePreviewButtons := getButtons(*votePreviewMsg.ReplyMarkup)
	s.Require().Len(votePreviewButtons, 1)
	// the only button is "Back"
	s.Require().Contains(votePreviewButtons[0].Text, "Back")
}

func (s *TelegramServiceSuite) mockTelegramSender(sentMsg *interface{}) {
	s.mockTgBot.On("Send", mock.Anything).
		Return(tgbotapi.Message{}, nil).
		Run(func(args mock.Arguments) {
			*sentMsg = args.Get(0)
		})
}

func (s *TelegramServiceSuite) sendCallback(button tgbotapi.InlineKeyboardButton) {
	s.T().Helper()

	s.Require().NotNil(button.CallbackData)

	update := callback(*button.CallbackData)
	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().NoError(err)
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

	keyboard, ok := in.(tgbotapi.InlineKeyboardMarkup)
	s.Require().True(ok)

	return getButtons(keyboard)
}
