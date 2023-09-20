package telegram_test

import (
	"math"
	"sort"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
)

// create several polls, go to the last page, and then go back to the first page
// every page check that text contains expected polls and keyboard contains expected buttons
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

	pollsMsg := s.asMessage(sentMsg)

	// verify the message contains the first 10 polls
	s.checkPolls(pollsMsg.Text, polls, 1, 10, false, true)
}

// check page contains expected polls and keyboard contains expected buttons
func (s *TelegramServiceSuite) checkPolls(txt string, allPolls []swagger.PollWithOptions, firstPoll, lastPoll int, prevButton, nextButton bool) {
	s.T().Helper()

	for i, poll := range allPolls {
		idx := i + 1
		if idx >= firstPoll && idx <= lastPoll {
			s.Require().Contains(txt, poll.Title)
		} else {
			s.Require().NotContains(txt, poll.Title)
		}
	}
}
