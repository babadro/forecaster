package telegram_test

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// open create poll page and click back button
func (s *TelegramServiceSuite) TestBackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	// open create poll page
	userID := randomPositiveInt64()
	update := startMainPage(userID)

	s.sendMessage(update)

	mainPage := s.asMessage(sentMsg)

	// click myPolls button
	myPollsButton := findItemByCriteria(s, s.buttonsFromMarkup(mainPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "My polls"
		},
	)

	s.sendCallback(myPollsButton, userID)

	myPollsMessage := s.asEditMessage(sentMsg)

	// click create poll button
	createPollButton := findItemByCriteria(s, s.buttonsFromMarkup(myPollsMessage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Create poll"
		},
	)

	s.sendCallback(createPollButton, userID)

	createPollPage := s.asEditMessage(sentMsg)

	// click back button
	backButton := findItemByCriteria(s, s.buttonsFromMarkup(createPollPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Go back"
		},
	)

	s.sendCallback(backButton, userID)

	// check that we are on "my polls" page
	myPollsPage := s.asEditMessage(sentMsg)

	s.Require().Contains(myPollsPage.Text, "Getting polls are not implemented yet")
}

func (s *TelegramServiceSuite) goToCreatePollPage(userID int64, sentMsg interface{}) tgbotapi.EditMessageTextConfig {
	s.T().Helper()

	update := startMainPage(userID)

	s.sendMessage(update)

	mainPage := s.asMessage(sentMsg)

	// click myPolls button
	myPollsButton := findItemByCriteria(s, s.buttonsFromMarkup(mainPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "My polls"
		},
	)

	s.sendCallback(myPollsButton, userID)

	myPollsMessage := s.asEditMessage(sentMsg)

	// click create poll button
	createPollButton := findItemByCriteria(s, s.buttonsFromMarkup(myPollsMessage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Create poll"
		},
	)

	s.sendCallback(createPollButton, userID)

	return s.asEditMessage(sentMsg)
}
