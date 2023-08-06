package telegram_test

import (
	"context"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/mock"
)

func (s *TelegramServiceSuite) TestProcessTelegramUpdate_happyPath() {
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
			Text: "/start showpoll_" + strconv.Itoa(int(poll.ID)),
		},
	}

	s.mockTgBot.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		s.Require().Equal(update.Message.Chat.ID, msg.ChatID)

		s.Require().Contains(msg.Text, poll.Title)

		for _, op := range createdOptions {
			s.Require().Contains(msg.Text, op.Title)
		}

		return true
	})).Return(tgbotapi.Message{}, nil)

	err = s.telegramService.ProcessTelegramUpdate(&s.logger, update)

	s.Require().NoError(err)
	s.mockTgBot.AssertExpectations(s.T())
}

func (s *TelegramServiceSuite) TestProcessTelegramUpdate_notFound() {
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 123,
			},
			Text: "/start showpoll_999",
		},
	}

	s.mockTgBot.On("Send", mock.MatchedBy(func(msg tgbotapi.MessageConfig) bool {
		s.Require().Equal(update.Message.Chat.ID, msg.ChatID)

		s.Require().Contains(msg.Text, "can't find poll")

		return true
	})).Return(tgbotapi.Message{}, nil)

	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)

	s.Require().NoError(err)

	logOutput := s.logOutput.String()
	s.Require().Contains(logOutput, "error")
	s.Require().Contains(logOutput, "unable to get poll by id")

	s.mockTgBot.AssertExpectations(s.T())
}
