package telegram_test

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *TelegramServiceSuite) TestEditOptionBackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	createPollPage := s.goToCreatePollPage(userID, &sentMsg)

	_ = createPollPage
	//
}

func (s *TelegramServiceSuite) findButtonByLowerText(text string, markup interface{}) tgbotapi.InlineKeyboardButton {
	return findItemByCriteria(s, s.buttonsFromMarkup(markup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return strings.ToLower(button.Text) == text
		},
	)
}

func (s *TelegramServiceSuite) findButtonByContainsLowerText(text string, markup interface{}) tgbotapi.InlineKeyboardButton {
	return findItemByCriteria(s, s.buttonsFromMarkup(markup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return strings.Contains(strings.ToLower(button.Text), text)
		},
	)
}
