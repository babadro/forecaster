package telegram_test

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *TelegramServiceSuite) TestMyPolls_pagination() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	userID := randomPositiveInt64()

	polls := s.createRandomPolls(24, withTelegramUserID(userID))
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	myPollsPage := s.goToMyPollsPage(userID, &sentMsg)
	txt, buttons := myPollsPage.Text, s.buttonsFromMarkup(myPollsPage.ReplyMarkup)

	// verify the first page
	s.verifyMyPollsPage(txt, buttons, polls, 1, 10, false, true)

	// send "Next" button
	nextButton := findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	pollsPage2 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage2.Text, s.buttonsFromMarkup(pollsPage2.ReplyMarkup)

	// verify the second page
	s.verifyMyPollsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Next" button
	nextButton = findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	pollsPage3 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage3.Text, s.buttonsFromMarkup(pollsPage3.ReplyMarkup)

	// verify the third page
	s.verifyMyPollsPage(txt, buttons, polls, 21, 24, true, false)
}

// check page contains expected polls and keyboard contains expected buttons...
func (s *TelegramServiceSuite) verifyMyPollsPage(
	txt string, buttons []tgbotapi.InlineKeyboardButton, allPolls []swagger.PollWithOptions,
	firstPoll, lastPoll int, prevButton, nextButton bool,
) {
	s.T().Helper()

	for i, poll := range allPolls {
		idx := i + 1
		if idx >= firstPoll && idx <= lastPoll {
			s.Require().Contains(txt, poll.Title)
		} else {
			s.Require().NotContains(txt, poll.Title)
		}
	}

	pollsCount := lastPoll - firstPoll + 1

	expectedButtonsCount := pollsCount + 2 // +2 for Main Menu button and Create Poll button

	s.buttonsContainsText(buttons, "Create poll")

	if prevButton {
		expectedButtonsCount++

		s.buttonsContainsText(buttons, "Prev")
	}

	if nextButton {
		expectedButtonsCount++

		s.buttonsContainsText(buttons, "Next")
	}

	s.Require().Len(buttons, expectedButtonsCount)

	for i := 0; i < pollsCount; i++ {
		s.buttonsContainsText(buttons, strconv.Itoa(i+1))
	}
}

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

func (s *TelegramServiceSuite) goToMyPollsPage(
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

	return myPollsMessage
}
