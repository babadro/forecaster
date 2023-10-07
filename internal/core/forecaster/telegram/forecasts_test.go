package telegram_test

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// create several polls, go to the last page, and then go back to the first page
// every page check that text contains expected forecasts and keyboard contains expected buttons...
func (s *TelegramServiceSuite) TestForecasts_pagination() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createForecasts(24)

	// send /start showforecasts_1 command
	userID := randomPositiveInt64()
	update := startShowForecasts(userID)

	s.sendMessage(update)

	forecastsPageStartCommand := s.asMessage(sentMsg)
	txt, buttons := forecastsPageStartCommand.Text, s.buttonsFromInterface(forecastsPageStartCommand.ReplyMarkup)

	// verify the first page
	s.verifyForecastsPage(txt, buttons, polls, 1, 10, false, true)

	// send "Next" button
	nextButton := findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	forecastsPage2 := s.asEditMessage(sentMsg)
	txt, buttons = forecastsPage2.Text, s.buttonsFromInterface(forecastsPage2.ReplyMarkup)

	// verify the second page
	s.verifyForecastsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Next" button
	nextButton = findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	pollsPage3 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage3.Text, s.buttonsFromInterface(pollsPage3.ReplyMarkup)

	// verify the third page
	s.verifyPollsPage(txt, buttons, polls, 21, 24, true, false)

	// send "Prev" button
	prevButton := findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Prev")
	})
	s.sendCallback(prevButton, userID)

	forecastsPage2 = s.asEditMessage(sentMsg)
	txt, buttons = forecastsPage2.Text, s.buttonsFromInterface(forecastsPage2.ReplyMarkup)

	// verify the second page
	s.verifyForecastsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Prev" button
	prevButton = findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Prev")
	})
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

	// check polls title
	for i, forecast := range allPolls {
		idx := i + 1
		if idx >= first && idx <= last {
			s.Require().Contains(txt, forecast.Title)
		} else {
			s.Require().NotContains(txt, forecast.Title)
		}
	}

	forecastsCount := last - first + 1

	// check statistics
	// each poll has 3 options and each option has 1 vote,
	// so expected percentage for each option is 33%
	s.Require().Equal(strings.Count(txt, "33%"), forecastsCount)

	// check buttons
	expectedButtonsCount := forecastsCount + 1 // +1 for Main Menu button

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

// create several polls, chose the first poll, go to the poll page, and then go back to the polls page...
func (s *TelegramServiceSuite) TestForecasts_chose_forecast() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createForecasts(2)

	// send /start showforecasts_1 command
	userID := randomPositiveInt64()
	update := startShowForecasts(userID)

	s.sendMessage(update)

	forecastsPageStartCommand := s.asMessage(sentMsg)

	firstForecastButton := findItemByCriteria(s,
		s.buttonsFromInterface(forecastsPageStartCommand.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "1"
		})

	s.sendCallback(firstForecastButton, userID)

	// verify the forecast message
	forecastMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(forecastMsg.Text, polls[0].Title)

	// verify AllForecasts button
	buttons := s.buttonsFromInterface(forecastMsg.ReplyMarkup)
	allForecastsButtons := buttons[len(buttons)-2]
	s.Require().Contains(allForecastsButtons.Text, "All Forecasts")

	// send AllForecasts button
	s.sendCallback(allForecastsButtons, userID)

	// verify the forecasts page
	forecastsMessage := s.asEditMessage(sentMsg)
	s.verifyForecastsPage(
		forecastsMessage.Text, s.buttonsFromInterface(forecastsMessage.ReplyMarkup), polls,
		1, 2, false, false,
	)
}

func (s *TelegramServiceSuite) TestShowForecastStartCommand() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createForecast()

	userID := randomPositiveInt64()
	update := startShowForecast(poll.ID, userID)

	s.sendMessage(update)

	forecastMsg := s.asMessage(sentMsg)

	s.verifyForecastPage(forecastMsg.Text, s.buttonsFromInterface(forecastMsg.ReplyMarkup), poll)
}

func (s *TelegramServiceSuite) TestForecastRenderCallback() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createForecast()

	userID := randomPositiveInt64()
	update := startShowForecasts(userID)

	s.sendMessage(update)

	forecastsPage := s.asMessage(sentMsg)

	firstForecastButton := findItemByCriteria(s,
		s.buttonsFromInterface(forecastsPage.ReplyMarkup), func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "1"
		})

	s.sendCallback(firstForecastButton, userID)

	forecastMsg := s.asEditMessage(sentMsg)

	s.verifyForecastPage(forecastMsg.Text, s.buttonsFromInterface(forecastMsg.ReplyMarkup), poll)
}

func (s *TelegramServiceSuite) Test_forecast_unavailable() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createRandomPoll(time.Now())

	userID := randomPositiveInt64()
	update := startShowForecast(poll.ID, userID)

	s.sendMessage(update)

	// verify the forecast message
	forecastMsg := s.asMessage(sentMsg)
	s.Require().Contains(forecastMsg.Text, "Forecast Unavailable")
}

func findItemByCriteria[T any](s *TelegramServiceSuite, items []T, f func(item T) bool) T {
	s.T().Helper()

	for _, item := range items {
		if f(item) {
			return item
		}
	}

	s.Fail("item not found")

	var zero T

	return zero
}

func (s *TelegramServiceSuite) createForecast() swagger.PollWithOptions {
	s.T().Helper()

	poll := s.createRandomPoll(time.Now())

	for optionIDx, votesCount := range []int{3, 2, 1} {
		for i := 0; i < votesCount; i++ {
			_, err := s.db.CreateVote(context.Background(), swagger.CreateVote{
				OptionID: poll.Options[optionIDx].ID,
				PollID:   poll.ID,
				UserID:   randomPositiveInt64(),
			}, time.Now().Unix())

			s.Require().NoError(err)
		}
	}

	s.Require().NoError(s.db.CalculateStatistics(context.Background(), poll.ID))

	return poll
}

func (s *TelegramServiceSuite) verifyForecastPage(
	gotText string, gotButtons []tgbotapi.InlineKeyboardButton, p swagger.PollWithOptions,
) {
	s.Require().Contains(gotText, p.Title)

	// verify options
	for _, op := range p.Options {
		s.Require().Contains(gotText, op.Title)
	}

	// verify statistics
	// first option has 3 votes, second option has 2 votes, third option has 1 vote
	s.Require().Contains(gotText, "50%")
	s.Require().Contains(gotText, "3 votes")
	s.Require().Contains(gotText, "33%")
	s.Require().Contains(gotText, "2 votes")
	s.Require().Contains(gotText, "17%")
	s.Require().Contains(gotText, "1 vote")

	// verify AllForecasts button
	allForecastsButtons := gotButtons[len(gotButtons)-2]
	s.Require().Contains(allForecastsButtons.Text, "All Forecasts")

	// verify ShowPoll button
	showPollButtons := gotButtons[len(gotButtons)-1]
	s.Require().Contains(showPollButtons.Text, "Show Poll")
}

// go to poll from the forecast and get back to the forecast...
func (s *TelegramServiceSuite) TestForecastShowPollAndBack() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	poll := s.createForecast()

	userID := randomPositiveInt64()
	update := startShowForecast(poll.ID, userID)

	s.sendMessage(update)

	forecastMsg := s.asMessage(sentMsg)

	// verify ShowPoll button
	showPollButton := findItemByCriteria(s,
		s.buttonsFromInterface(forecastMsg.ReplyMarkup), func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Show Poll"
		})

	s.sendCallback(showPollButton, userID)

	// verify the poll message
	pollMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(pollMsg.Text, poll.Title)

	// vote for the first option
	firstOptionButton := findItemByCriteria(s,
		s.buttonsFromInterface(pollMsg.ReplyMarkup), func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "1"
		})

	s.sendCallback(firstOptionButton, userID)

	// verify the vote message
	voteMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(voteMsg.Text, "Vote for this option?")

	// get back to poll
	backToPollButton := findItemByCriteria(s,
		s.buttonsFromInterface(voteMsg.ReplyMarkup), func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Back"
		})

	s.sendCallback(backToPollButton, userID)

	// verify the poll message
	pollMsg = s.asEditMessage(sentMsg)
	s.Require().Contains(pollMsg.Text, poll.Title)

	// get back to forecast
	backToForecastButton := findItemByCriteria(s,
		s.buttonsFromInterface(pollMsg.ReplyMarkup), func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Show Forecast"
		})

	s.sendCallback(backToForecastButton, userID)

	// verify the forecast message
	forecastEditMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(forecastEditMsg.Text, poll.Title)
}

func (s *TelegramServiceSuite) TestForecasts_polls_without_total_votes_should_not_be_shown() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	pollWithStatistic := s.createForecast()

	pollWithoutStatistic := s.createRandomPoll(time.Now())

	// send /start showforecasts_1 command
	userID := randomPositiveInt64()
	update := startShowForecasts(userID)

	s.sendMessage(update)

	forecastsPageStartCommand := s.asMessage(sentMsg)

	s.Require().Contains(forecastsPageStartCommand.Text, pollWithStatistic.Title)
	s.Require().NotContains(forecastsPageStartCommand.Text, pollWithoutStatistic.Title)
}

func (s *TelegramServiceSuite) createForecasts(count int) []swagger.PollWithOptions {
	s.T().Helper()

	polls := s.createRandomPolls(count)
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	ctx := context.Background()

	// vote in each poll so that every poll has votes
	for _, p := range polls {
		// vote for each option
		for _, op := range p.Options {
			_, err := s.db.CreateVote(ctx, swagger.CreateVote{
				OptionID: op.ID,
				PollID:   p.ID,
				UserID:   randomPositiveInt64(),
			}, time.Now().Unix())

			s.Require().NoError(err)
		}
	}

	// calculate statistic to fill total_votes field
	for _, p := range polls {
		err := s.db.CalculateStatistics(ctx, p.ID)
		s.Require().NoError(err)
	}

	return polls
}
