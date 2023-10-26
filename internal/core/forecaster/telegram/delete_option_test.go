package telegram_test

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// open create poll page and click back button...
func (s *TelegramServiceSuite) TestDeleteOption() {
	userID := randomPositiveInt64()

	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	editOptionPage, optionTitle := s.createPollAndOptionAndGoToToEditOptionPage(userID, &sentMsg)

	s.Require().Contains(editOptionPage.Text, optionTitle)

	optionKeyboard := editOptionPage.ReplyMarkup

	deleteButton := s.findButtonByLowerText("delete option", optionKeyboard)

	s.sendCallback(deleteButton, userID)

	deleteConfirmation := s.asEditMessage(sentMsg)

	s.Require().Contains(deleteConfirmation.Text, "Are you sure you want to delete this option?")

	deleteButton = s.findButtonByLowerText("delete", deleteConfirmation.ReplyMarkup)

	s.sendCallback(deleteButton, userID)

	successDeletionResultMsg := s.asEditMessage(sentMsg)

	s.Require().Contains(successDeletionResultMsg.Text, "Option was successfully deleted!")

	backButton := s.findButtonByLowerText("go back", successDeletionResultMsg.ReplyMarkup)

	s.sendCallback(backButton, userID)

	pollMessage := s.asEditMessage(sentMsg)

	// check that the option is deleted
	s.Require().NotContains(pollMessage.Text, optionTitle)
}

// open delete option page and click back button without deleting...
func (s *TelegramServiceSuite) TestBackButton() {
	userID := randomPositiveInt64()

	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	editOptionPage, optionTitle := s.createPollAndOptionAndGoToToEditOptionPage(userID, &sentMsg)

	deleteButton := s.findButtonByLowerText("delete option", editOptionPage.ReplyMarkup)

	s.sendCallback(deleteButton, userID)

	deleteConfirmation := s.asEditMessage(sentMsg)

	backButton := s.findButtonByLowerText("go back", deleteConfirmation.ReplyMarkup)

	s.sendCallback(backButton, userID)

	pollMessage := s.asEditMessage(sentMsg)

	// check that the option is not deleted
	s.Require().Contains(pollMessage.Text, optionTitle)
}

func (s *TelegramServiceSuite) createPollAndOptionAndGoToToEditOptionPage(
	userID int64, sentMsg *interface{},
) (tgbotapi.MessageConfig, string) {
	s.T().Helper()

	pollKeyboard := s.createPollAndGoToEditPollPage(userID, sentMsg).ReplyMarkup

	createOptionButton := s.findButtonByLowerText("add option", pollKeyboard)

	s.sendCallback(createOptionButton, userID)

	optionKeyboard := s.asEditMessage(*sentMsg).ReplyMarkup

	titleButton := s.findButtonByLowerText("title", optionKeyboard)

	s.sendCallback(titleButton, userID)

	// check that we are on the edit page for the button
	editPage := s.asEditMessage(*sentMsg)

	// reply to the edit page with some text
	optionTitle := randomSentence()
	reply := replyMessageUpdate(optionTitle, editPage.Text, userID)

	s.sendMessage(reply)

	editOptionPageNewMessage := s.asMessage(*sentMsg)

	s.Require().Contains(editOptionPageNewMessage.Text, optionTitle)

	return editOptionPageNewMessage, optionTitle
}
