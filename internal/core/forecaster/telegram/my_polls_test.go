package telegram_test

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (s *TelegramServiceSuite) TestMyPollsBackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

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

	myPollsPage := s.asEditMessage(sentMsg)

	// click back button
	backButton := findItemByCriteria(s, s.buttonsFromMarkup(myPollsPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Main Menu"
		},
	)

	s.sendCallback(backButton, userID)

	// check that we are on main page
	mainPageEditMessage := s.asEditMessage(sentMsg)

	s.Require().Contains(mainPageEditMessage.Text, sentenceFromMainPage)
}
