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

	s.Require().Contains(myPollsPage.Text, "There are no polls yet")
}

func (s *TelegramServiceSuite) TestCreatePoll() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	var pollKeyboard any = s.goToCreatePollPage(userID, &sentMsg).ReplyMarkup

	var pollMessage tgbotapi.MessageConfig

	for _, buttonName := range []string{"title", "description"} {
		button := s.findButtonByLowerText(buttonName, pollKeyboard)

		s.sendCallback(button, userID)

		// check that we are on the edit page for the button
		editPage := s.asEditMessage(sentMsg)

		s.Require().Contains(strings.ToLower(editPage.Text), buttonName)

		// reply to the edit page with some text
		userInput := randomSentence()
		reply := replyMessageUpdate(userInput, editPage.Text, userID)

		s.sendMessage(reply)

		// check that we are on the edit poll page and this page contains user input
		pollMessage = s.asMessage(sentMsg)

		s.Require().Contains(pollMessage.Text, userInput)

		pollKeyboard = pollMessage.ReplyMarkup
	}

	for _, dateButtonName := range []string{"start", "finish"} {
		dateButton := s.findButtonByContainsLowerText(dateButtonName, pollMessage.ReplyMarkup)

		s.sendCallback(dateButton, userID)

		// check that we are on the edit page for the button
		editPage := s.asEditMessage(sentMsg)

		s.Require().Contains(strings.ToLower(editPage.Text), dateButtonName)

		// reply to the edit page with some text
		validDate := gofakeit.Date()
		reply := replyMessageUpdate(validDate.Format(time.RFC3339), editPage.Text, userID)

		s.sendMessage(reply)

		// check that we are on the edit poll page and this page contains user input
		pollMessage = s.asMessage(sentMsg)

		const dateErrorMsg = "Can't parse date format"

		s.Require().NotContains(pollMessage.Text, dateErrorMsg)

		// invalid date case
		invalidDate := "some-invalid-date"
		reply = replyMessageUpdate(invalidDate, editPage.Text, userID)

		s.sendMessage(reply)

		pollMessage = s.asMessage(sentMsg)

		s.Require().Contains(pollMessage.Text, dateErrorMsg)
	}
}

func (s *TelegramServiceSuite) goToCreatePollPage(
	userID int64, sentMsgPtr *interface{},
) tgbotapi.EditMessageTextConfig {
	s.T().Helper()

	myPollsMessage := s.goToMyPollsPage(userID, sentMsgPtr)

	// click create poll button
	createPollButton := s.findButtonByLowerText("create poll", myPollsMessage.ReplyMarkup)

	s.sendCallback(createPollButton, userID)

	return s.asEditMessage(*sentMsgPtr)
}

func (s *TelegramServiceSuite) TestCreatePoll_error_title_should_be_created_first() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	createPollKeyboard := s.goToCreatePollPage(userID, &sentMsg).ReplyMarkup

	// click description button
	descriptionButton := s.findButtonByLowerText("description", createPollKeyboard)

	s.sendCallback(descriptionButton, userID)

	// check that we are on the edit page for the button
	editPage := s.asEditMessage(sentMsg)

	s.Require().Contains(editPage.Text, "First create Title, please, and then you can create other fields")
}

func (s *TelegramServiceSuite) TestEditPollStatus_validation_error() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	// create poll with only title, that is not enough to activate poll
	// for activation poll should have all fields and at least 2 options
	pollKeyboard := s.createPollWithTitleOnlyAndGoToEditPollPage(userID, &sentMsg).ReplyMarkup

	activateButton := s.findButtonByContainsLowerText("activate", pollKeyboard)

	s.sendCallback(activateButton, userID)

	confirmationPageKeyboard := s.asEditMessage(sentMsg).ReplyMarkup

	activateButton = s.findButtonByContainsLowerText("activate", confirmationPageKeyboard)

	s.sendCallback(activateButton, userID)

	validationErrorMessage := s.asEditMessage(sentMsg)

	s.Require().Contains(validationErrorMessage.Text, "Can't change status to active, due to validation errors")
}

func (s *TelegramServiceSuite) TestEditPollStatus_success() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	editPollPage, p := s.createPollReadyForActivationAndGoToEditPollPage(userID, &sentMsg)

	activateButton := s.findButtonByContainsLowerText("activate", editPollPage.ReplyMarkup)

	s.sendCallback(activateButton, userID)

	confirmationPageKeyboard := s.asEditMessage(sentMsg).ReplyMarkup

	activateButton = s.findButtonByContainsLowerText("activate", confirmationPageKeyboard)

	s.sendCallback(activateButton, userID)

	successPage := s.asEditMessage(sentMsg)

	s.Require().Contains(successPage.Text, p.Title)
	s.Require().Contains(successPage.Text, "is active now")
}
