package telegram_test

import (
	"math"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	sentenceFromMainPage      = "Welcome to Forecaster Bot"
	sentenceFromForecastsPage = "Most popular option"
)

// open main page with command, go to forecasts page, go back to main page, go to polls page, go back to main page
func (s *TelegramServiceSuite) TestMainPage() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createForecasts(2)

	// send /show main command
	userID := int64(gofakeit.IntRange(1, math.MaxInt64))
	update := startMainPage(userID)

	s.sendMessage(update)

	mainPage := s.asMessage(sentMsg)

	// check that we are on the main page
	s.Require().Contains(mainPage.Text, sentenceFromMainPage)

	buttons := s.buttonsFromInterface(mainPage.ReplyMarkup)

	// go to forecasts page
	forecastsButton := findItemByCriteria(s, buttons,
		func(button tgbotapi.InlineKeyboardButton) bool {
			return strings.Contains(button.Text, "forecasts")
		})

	s.sendCallback(forecastsButton, userID)

	forecastsPage := s.asEditMessage(sentMsg)

	// check that we are on the forecasts page
	for _, poll := range polls {
		s.Require().Contains(forecastsPage.Text, poll.Title)
	}

	// by checking this text we check that it's a forecasts page with some statistics,
	// not a polls page
	s.Require().Contains(forecastsPage.Text, sentenceFromForecastsPage)

	// go back to main page
	mainMenuButton := s.findMainButton(s.buttonsFromInterface(forecastsPage.ReplyMarkup))

	s.sendCallback(mainMenuButton, userID)

	mainEditPage := s.asEditMessage(sentMsg)
	s.Require().Contains(mainEditPage.Text, sentenceFromMainPage)

	// go to polls page
	pollsButton := findItemByCriteria(s, s.buttonsFromInterface(mainEditPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return strings.Contains(button.Text, "polls")
		})

	s.sendCallback(pollsButton, userID)

	pollsPage := s.asEditMessage(sentMsg)

	// check that we are on the polls page
	for _, poll := range polls {
		s.Require().Contains(pollsPage.Text, poll.Title)
	}

	s.Require().NotContains(pollsPage.Text, sentenceFromForecastsPage)

	// go back to main page

	mainMenuButton = s.findMainButton(s.buttonsFromInterface(pollsPage.ReplyMarkup))

	s.sendCallback(mainMenuButton, userID)

	mainEditPage = s.asEditMessage(sentMsg)

	s.Require().Contains(mainEditPage.Text, sentenceFromMainPage)
}

func (s *TelegramServiceSuite) findMainButton(buttons []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardButton {
	return findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return button.Text == "Main Menu"
	})
}
