package telegram_test

import (
	"context"
)

// open create poll page and click back button...
func (s *TelegramServiceSuite) TestDeleteOption() {
	userID := randomPositiveInt64()

	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	pollKeyboard := s.createPollAndGoToEditPollPage(userID, &sentMsg).ReplyMarkup

	createOptionButton := s.findButtonByLowerText("add option", pollKeyboard)

	s.sendCallback(createOptionButton, userID)

	var optionKeyboard any = s.asEditMessage(sentMsg).ReplyMarkup

	titleButton := s.findButtonByLowerText("title", optionKeyboard)

	s.sendCallback(titleButton, userID)

	// check that we are on the edit page for the button
	editPage := s.asEditMessage(sentMsg)

	// reply to the edit page with some text
	userInput := randomSentence()
	reply := replyMessageUpdate(userInput, editPage.Text, userID)

	s.sendMessage(reply)

	editOptionPageNewMessage := s.asMessage(sentMsg)

	optionKeyboard = editOptionPageNewMessage.ReplyMarkup

	deleteButton := s.findButtonByLowerText("delete option", optionKeyboard)

	s.sendCallback(deleteButton, userID)

	deleteConfirmation := s.asEditMessage(sentMsg)

	s.Require().Contains(deleteConfirmation.Text, "Are you sure you want to delete this option?")

	// verify that option was not deleted yet
	pollsArr, _, err := s.db.GetPolls(context.Background(), 0, 1)
	s.Require().NoError(err)
	s.Require().Len(pollsArr, 1)

	p := pollsArr[0]

	pollWithOptions, err := s.db.GetPollByID(context.Background(), p.ID)
	s.Require().NoError(err)

	s.Require().Len(pollWithOptions.Options, 1)

}
