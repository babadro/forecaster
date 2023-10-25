package telegram_test

/*
// open create poll page and click back button...
func (s *TelegramServiceSuite) TestDeleteOption() {

	s.createRandomPoll()

	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	pollKeyboard := s.createPollAndGoToEditPollPage(userID, &sentMsg).ReplyMarkup

	createOptionButton := s.findButtonByLowerText("add option", pollKeyboard)

	s.sendCallback(createOptionButton, userID)

	// check that we are on edit option page
	var optionKeyboard any = s.asEditMessage(sentMsg).ReplyMarkup

	// click title button
	titleButton := s.findButtonByLowerText("title", optionKeyboard)

	s.sendCallback(titleButton, userID)

	editPage := s.asEditMessage(sentMsg)

	userInput := randomSentence()
	reply := replyMessageUpdate(userInput, editPage.Text, userID)

	s.sendMessage(reply)
}

*/
