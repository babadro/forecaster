package telegram_test

import (
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// open create poll page and click back button...
func (s *TelegramServiceSuite) TestEditPollBackButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	createPollPage := s.goToCreatePollPage(userID, &sentMsg)

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

func (s *TelegramServiceSuite) TestCreatePoll() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	var createOrEditPollKeyboard any = s.goToCreatePollPage(userID, &sentMsg).ReplyMarkup

	var createPollMessage tgbotapi.MessageConfig

	for _, buttonName := range []string{"title", "description"} {
		button := findItemByCriteria(s, s.buttonsFromMarkup(createOrEditPollKeyboard),
			func(button tgbotapi.InlineKeyboardButton) bool {
				return strings.ToLower(button.Text) == buttonName
			},
		)

		s.sendCallback(button, userID)

		// check that we are on the edit page for the button
		editPage := s.asEditMessage(sentMsg)

		s.Require().Contains(strings.ToLower(editPage.Text), buttonName)

		// reply to the edit page with some text
		userInput := randomSentence()
		reply := replyMessageUpdate(userInput, editPage.Text, userID)

		s.sendMessage(reply)

		// check that we are on the edit poll page and this page contains user input
		createPollMessage = s.asMessage(sentMsg)

		s.Require().Contains(createPollMessage.Text, userInput)

		createOrEditPollKeyboard = createPollMessage.ReplyMarkup
	}

	for _, dateButtonName := range []string{"start", "finish"} {
		dateButton := findItemByCriteria(s, s.buttonsFromMarkup(createPollMessage.ReplyMarkup),
			func(button tgbotapi.InlineKeyboardButton) bool {
				return strings.Contains(strings.ToLower(button.Text), dateButtonName)
			},
		)

		s.sendCallback(dateButton, userID)

		// check that we are on the edit page for the button
		editPage := s.asEditMessage(sentMsg)

		s.Require().Contains(strings.ToLower(editPage.Text), dateButtonName)

		// reply to the edit page with some text
		validDate := gofakeit.Date()
		reply := replyMessageUpdate(validDate.Format(time.RFC3339), editPage.Text, userID)

		s.sendMessage(reply)

		// check that we are on the edit poll page and this page contains user input
		createPollMessage = s.asMessage(sentMsg)

		const dateErrorMsg = "Can't parse date format"

		s.Require().NotContains(createPollMessage.Text, dateErrorMsg)

		// invalid date case
		invalidDate := "some-invalid-date"
		reply = replyMessageUpdate(invalidDate, editPage.Text, userID)

		s.sendMessage(reply)

		createPollMessage = s.asMessage(sentMsg)

		s.Require().Contains(createPollMessage.Text, dateErrorMsg)
	}
}

func (s *TelegramServiceSuite) goToCreatePollPage(
	userID int64, sentMsgPtr *interface{},
) tgbotapi.EditMessageTextConfig {
	s.T().Helper()

	update := startMainPage(userID)

	s.sendMessage(update)

	mainPage := s.asMessage(*sentMsgPtr)

	// click myPolls button
	myPollsButton := findItemByCriteria(s, s.buttonsFromMarkup(mainPage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "My polls"
		},
	)

	s.sendCallback(myPollsButton, userID)

	myPollsMessage := s.asEditMessage(*sentMsgPtr)

	// click create poll button
	createPollButton := findItemByCriteria(s, s.buttonsFromMarkup(myPollsMessage.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "Create poll"
		},
	)

	s.sendCallback(createPollButton, userID)

	return s.asEditMessage(*sentMsgPtr)
}
