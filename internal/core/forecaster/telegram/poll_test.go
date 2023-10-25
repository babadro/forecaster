package telegram_test

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

func (s *TelegramServiceSuite) TestShowPollStartCommand() {
	poll := s.createRandomPoll(withNow(time.Now()))

	update := startShowPoll(poll.ID, 456)

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
	update := startShowPoll(999, 456)

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

	buttons := s.buttonsFromMarkup(msg.ReplyMarkup)
	s.Require().Len(buttons, 1)
	s.Require().Contains(buttons[0].Text, "Back to main")

	s.mockTgBot.AssertExpectations(s.T())
}
