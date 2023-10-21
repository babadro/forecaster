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
	backButton := s.findButtonByLowerText("go back", createPollPage.ReplyMarkup)

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
		button := s.findButtonByLowerText(buttonName, createOrEditPollKeyboard)

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
		dateButton := s.findButtonByContainsLowerText(dateButtonName, createPollMessage.ReplyMarkup)

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
	myPollsButton := s.findButtonByLowerText("my polls", mainPage.ReplyMarkup)

	s.sendCallback(myPollsButton, userID)

	myPollsMessage := s.asEditMessage(*sentMsgPtr)

	// click create poll button
	createPollButton := s.findButtonByLowerText("create poll", myPollsMessage.ReplyMarkup)

	s.sendCallback(createPollButton, userID)

	return s.asEditMessage(*sentMsgPtr)
}
