package telegram_test

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// create several polls, go to the last page, and then go back to the first page
// every page check that text contains expected polls and keyboard contains expected buttons...
func (s *TelegramServiceSuite) TestPolls_pagination() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createRandomPolls(24)
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	// send /start showpolls_1 command
	userID := int64(gofakeit.IntRange(1, math.MaxInt64))
	update := startShowPolls(1, userID)

	s.sendMessage(update)

	pollsPageStartCommand := s.asMessage(sentMsg)
	txt, buttons := pollsPageStartCommand.Text, s.buttonsFromInterface(pollsPageStartCommand.ReplyMarkup)

	// verify the first page
	s.verifyPollsPage(txt, buttons, polls, 1, 10, false, true)

	// send "Next" button
	nextButton := findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	pollsPage2 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage2.Text, s.buttonsFromInterface(pollsPage2.ReplyMarkup)

	// verify the second page
	s.verifyPollsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Next" button
	nextButton = findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Next")
	})
	s.sendCallback(nextButton, userID)

	pollsPage3 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage3.Text, s.buttonsFromInterface(pollsPage3.ReplyMarkup)

	// verify the third page
	s.verifyPollsPage(txt, buttons, polls, 21, 24, true, false)

	// send "Prev" button
	prevButton := findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Prev")
	})
	s.sendCallback(prevButton, userID)

	pollsPage2 = s.asEditMessage(sentMsg)
	txt, buttons = pollsPage2.Text, s.buttonsFromInterface(pollsPage2.ReplyMarkup)

	// verify the second page
	s.verifyPollsPage(txt, buttons, polls, 11, 20, true, true)

	// send "Prev" button
	prevButton = findItemByCriteria(s, buttons, func(button tgbotapi.InlineKeyboardButton) bool {
		return strings.Contains(button.Text, "Prev")
	})
	s.sendCallback(prevButton, userID)

	pollsPage1 := s.asEditMessage(sentMsg)
	txt, buttons = pollsPage1.Text, s.buttonsFromInterface(pollsPage1.ReplyMarkup)

	// verify the first page
	s.verifyPollsPage(txt, buttons, polls, 1, 10, false, true)
}

// check page contains expected polls and keyboard contains expected buttons...
func (s *TelegramServiceSuite) verifyPollsPage(
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

	expectedButtonsCount := pollsCount + 1 // +1 for Main Menu button

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

func (s *TelegramServiceSuite) buttonsContainsText(buttons []tgbotapi.InlineKeyboardButton, text string) {
	for _, b := range buttons {
		if strings.Contains(b.Text, text) {
			return
		}
	}

	s.Fail("buttons does not contain text: " + text)
}

// create several polls, chose the first poll, go to the poll page, and then go back to the polls page...
func (s *TelegramServiceSuite) TestPolls_chose_poll() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	polls := s.createRandomPolls(2)
	// polls should be sorted by created_at desc
	sort.Slice(polls, func(i, j int) bool {
		return time.Time(polls[i].CreatedAt).Unix() > (time.Time(polls[j].CreatedAt).Unix())
	})

	// send /start showpolls_1 command
	userID := int64(gofakeit.IntRange(1, math.MaxInt64))
	update := startShowPolls(1, userID)

	s.sendMessage(update)

	pollsPageStartCommand := s.asMessage(sentMsg)

	firstPollButton := findItemByCriteria(s,
		s.buttonsFromInterface(pollsPageStartCommand.ReplyMarkup),
		func(button tgbotapi.InlineKeyboardButton) bool {
			return button.Text == "1"
		})

	s.sendCallback(firstPollButton, userID)

	// verify the poll message
	pollMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(pollMsg.Text, polls[0].Title)

	// verify AllPolls button
	buttons := s.buttonsFromInterface(pollMsg.ReplyMarkup)
	allPollButtons := buttons[len(buttons)-1]
	s.Require().Contains(allPollButtons.Text, "All Polls")

	// send AllPolls button
	s.sendCallback(allPollButtons, userID)

	// verify the polls page
	pollsMessage := s.asEditMessage(sentMsg)
	s.verifyPollsPage(pollsMessage.Text, s.buttonsFromInterface(pollsMessage.ReplyMarkup), polls, 1, 2, false, false)
}
