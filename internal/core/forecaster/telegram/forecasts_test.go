package telegram_test

import (
	"context"
	"github.com/babadro/forecaster/internal/models/swagger"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// create several polls, go to the last page, and then go back to the first page
// every page check that text contains expected forecasts and keyboard contains expected buttons...
func (s *TelegramServiceSuite) TestForecasts_pagination() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createRandomPolls(24)
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	userID := int64(gofakeit.IntRange(1, math.MaxInt64))

	// vote in each poll so that every poll has votes
	for _, p := range polls {
		op := p.Options[0]

		_, err := s.db.CreateVote(context.Background(), swagger.CreateVote{
			OptionID: op.ID,
			PollID:   p.ID,
			UserID:   userID,
		}, time.Now().Unix())

		s.Require().NoError(err)
	}

	// send /start showforecasts_1 command

	update := startShowForecasts(1, userID)

	s.sendMessage(update)

	forecastsPageStartCommand := s.asMessage(sentMsg)
	txt, buttons := forecastsPageStartCommand.Text, s.buttonsFromInterface(forecastsPageStartCommand.ReplyMarkup)

	// verify the first page
	s.verifyForecastsPage(txt, buttons, polls, 1, 10, false, true)

	// send "Next" button
	nextButton := buttons[0]
	s.Require().Contains(nextButton.Text, "Next")
	s.sendCallback(nextButton, userID)

	forecastsPage2 := s.asEditMessage(sentMsg)
	txt, buttons = forecastsPage2.Text, s.buttonsFromInterface(forecastsPage2.ReplyMarkup)

	// verify the second page
	s.verifyForecastsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Next" button
	nextButton = buttons[1]
	s.Require().Contains(nextButton.Text, "Next")
	s.sendCallback(nextButton, userID)

	pollsPage3 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage3.Text, s.buttonsFromInterface(pollsPage3.ReplyMarkup)

	// verify the third page
	s.verifyPollsPage(txt, buttons, polls, 21, 24, true, false)

	// send "Prev" button
	prevButton := buttons[0]
	s.Require().Contains(prevButton.Text, "Prev")
	s.sendCallback(prevButton, userID)

	forecastsPage2 = s.asEditMessage(sentMsg)
	txt, buttons = forecastsPage2.Text, s.buttonsFromInterface(forecastsPage2.ReplyMarkup)

	// verify the second page
	s.verifyForecastsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Prev" button
	prevButton = buttons[0]
	s.Require().Contains(prevButton.Text, "Prev")
	s.sendCallback(prevButton, userID)

	forecastsPage1 := s.asEditMessage(sentMsg)
	txt, buttons = forecastsPage1.Text, s.buttonsFromInterface(forecastsPage1.ReplyMarkup)

	// verify the first page
	s.verifyForecastsPage(txt, buttons, polls, 1, 10, false, true)
}

// check page contains expected forecasts and keyboard contains expected buttons...
func (s *TelegramServiceSuite) verifyForecastsPage(
	txt string, buttons []tgbotapi.InlineKeyboardButton, allPolls []swagger.PollWithOptions,
	first, last int, prevButton, nextButton bool,
) {
	s.T().Helper()

	for i, forecast := range allPolls {
		idx := i + 1
		if idx >= first && idx <= last {
			s.Require().Contains(txt, forecast.Title)
		} else {
			s.Require().NotContains(txt, forecast.Title)
		}
	}

	forecastsCount := last - first + 1

	expectedButtonsCount := forecastsCount

	if prevButton {
		expectedButtonsCount++

		s.buttonsContainsText(buttons, "Prev")
	}

	if nextButton {
		expectedButtonsCount++

		s.buttonsContainsText(buttons, "Next")
	}

	s.Require().Len(buttons, expectedButtonsCount)

	for i := 0; i < forecastsCount; i++ {
		s.buttonsContainsText(buttons, strconv.Itoa(i+1))
	}
}

/*
// create several polls, chose the first poll, go to the poll page, and then go back to the polls page...
func (s *TelegramServiceSuite) TestForecasts_chose_forecast() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createRandomPolls(2)
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	// send /start showpolls_1 command
	userID := int64(gofakeit.IntRange(1, math.MaxInt64))
	update := startShowForecasts(1, userID)

	s.sendMessage(update)

	forecastsPageStartCommand := s.asMessage(sentMsg)

	firstPollButton, found := tgbotapi.InlineKeyboardButton{}, false

	for _, button := range s.buttonsFromInterface(forecastsPageStartCommand.ReplyMarkup) {
		if button.Text == "1" {
			firstPollButton = button
			found = true
		}
	}

	s.Require().True(found)

	s.sendCallback(firstPollButton, userID)

	// verify the poll message
	pollMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(pollMsg.Text, polls[0].Title)

	// verify AllPolls button
	buttons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	allPollButtons := buttons[len(buttons)-1]
	s.Require().Contains(allPollButtons.Text, "All Polls")

	// send AllPolls button
	s.sendCallback(allPollButtons, userID)

	// verify the polls page
	pollsMessage := s.asEditMessage(sentMsg)
	s.verifyPollsPage(pollsMessage.Text, s.buttonsFromInterface(pollsMessage.ReplyMarkup), polls, 1, 2, false, false)
}
*/
