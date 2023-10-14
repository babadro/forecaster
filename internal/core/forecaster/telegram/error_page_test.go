package telegram_test

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *TelegramServiceSuite) TestUnknownCommand() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	// send unknown command
	userID := randomPositiveInt64()
	update := command("some-fake-unknown-command", userID)

	err := s.telegramService.ProcessTelegramUpdate(&s.logger, update)
	s.Require().ErrorContains(err, "unknown command")

	unknownCommand := s.asMessage(sentMsg)

	// check that we are on the error page
	s.Require().Contains(unknownCommand.Text, "I don't know this command")

	backToMainButton := findItemByCriteria(s, s.buttonsFromMarkup(unknownCommand.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Back to main"
		})

	s.sendCallback(backToMainButton, userID)

	mainPage := s.asEditMessage(sentMsg)

	// check that we are on the main page
	s.Require().Contains(mainPage.Text, sentenceFromMainPage)
}
