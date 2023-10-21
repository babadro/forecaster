package telegram_test

import (
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *TelegramServiceSuite) TestEditOptionBackButton() {
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

	backButton := s.findButtonByLowerText("go back", editPage.ReplyMarkup)

	s.sendCallback(backButton, userID)

	// check that we are on the edit poll page
	pollMessage := s.asEditMessage(sentMsg)

	s.Require().Contains(pollMessage.Text, "Define your option title and description.")
}

func (s *TelegramServiceSuite) TestCreateOption() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	pollKeyboard := s.createPollAndGoToEditPollPage(userID, &sentMsg).ReplyMarkup

	createOptionButton := s.findButtonByLowerText("add option", pollKeyboard)

	s.sendCallback(createOptionButton, userID)

	// check that we are on edit option page
	var optionKeyboard any = s.asEditMessage(sentMsg).ReplyMarkup

	for _, buttonName := range []string{"title", "description"} {
		button := s.findButtonByLowerText(buttonName, optionKeyboard)

		s.sendCallback(button, userID)

		// check that we are on the edit page for the button
		editPage := s.asEditMessage(sentMsg)

		s.Require().Contains(editPage.Text, models.EditOptionCommand)

		// reply to the edit page with some text
		userInput := randomSentence()
		reply := replyMessageUpdate(userInput, editPage.Text, userID)

		s.sendMessage(reply)

		// check that we are on the edit option page and this page contains user input
		editOptionPageNewMessage := s.asMessage(sentMsg)

		s.Require().Contains(editOptionPageNewMessage.Text, userInput)

		optionKeyboard = editOptionPageNewMessage.ReplyMarkup

		if buttonName == "title" {
			// go back to the poll page and check that option was created
			backButton := s.findButtonByLowerText("go back", optionKeyboard)

			s.sendCallback(backButton, userID)

			pollMessage := s.asEditMessage(sentMsg)

			s.Require().Contains(pollMessage.Text, userInput)

			// return to the edit option page
			firstOptionButton := s.findButtonByLowerText("1", pollMessage.ReplyMarkup)

			s.sendCallback(firstOptionButton, userID)

			optionKeyboard = s.asEditMessage(sentMsg).ReplyMarkup
		}
	}
}

func (s *TelegramServiceSuite) createPollAndGoToEditPollPage(userID int64, sentMsg *interface{}) tgbotapi.MessageConfig {
	s.T().Helper()

	createPollKeyboard := s.goToCreatePollPage(userID, sentMsg).ReplyMarkup

	// create poll
	titleButton := s.findButtonByLowerText("title", createPollKeyboard)

	s.sendCallback(titleButton, userID)

	editPage := s.asEditMessage(*sentMsg)

	reply := replyMessageUpdate(randomSentence(), editPage.Text, userID)

	s.sendMessage(reply)

	return s.asMessage(*sentMsg)
}

func (s *TelegramServiceSuite) TestCreateOption_error_title_should_be_created_first() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	pollKeyboard := s.createPollAndGoToEditPollPage(userID, &sentMsg).ReplyMarkup

	createOptionButton := s.findButtonByLowerText("add option", pollKeyboard)

	s.sendCallback(createOptionButton, userID)

	// check that we are on edit option page
	var optionKeyboard any = s.asEditMessage(sentMsg).ReplyMarkup

	// click description button
	descriptionButton := s.findButtonByLowerText("description", optionKeyboard)

	s.sendCallback(descriptionButton, userID)

	// check that we are on the edit page for the button
	editPage := s.asEditMessage(sentMsg)

	s.Require().Contains(editPage.Text, "First create Title, please, and then you can create other fields")
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
