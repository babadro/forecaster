package telegram_test

import (
	"context"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *TelegramServiceSuite) TestUserPollResult_callback_happy_path() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	p, targetUserID := s.setupForUserPollResultTest()

	update := startShowPoll(p.ID, targetUserID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromInterface(pollMsg.ReplyMarkup)

	s.Require().Len(pollButtons, len(p.Options)+1) // +1 for "Show results" button

	showResultsButton := pollButtons[len(pollButtons)-1]
	sentCallback := s.sendCallback(showResultsButton, targetUserID)
	userName := sentCallback.CallbackQuery.From.UserName

	// verify the user poll result message
	userPollResultMsg := s.asEditMessage(sentMsg)

	s.checkStatisticsForUserPollResultTest(userName, userPollResultMsg.Text)

	s.Require().Contains(userPollResultMsg.Text, "you") // in case of callback the message is for 2nd person
}

func (s *TelegramServiceSuite) TestUserPollResult_command_happy_path() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	p, targetUserID := s.setupForUserPollResultTest()

	// send /start showuserres_<poll_id>_<user_id> command
	update := startShowUserRes(p.ID, targetUserID)
	s.sendMessage(update)

	userPollResultMsg := s.asMessage(sentMsg)
	s.checkStatisticsForUserPollResultTest(update.Message.From.UserName, userPollResultMsg.Text)

	s.Require().NotContains(userPollResultMsg.Text, "you") // in case of /start command the message is for 3rd person
}

func (s *TelegramServiceSuite) setupForUserPollResultTest() (swagger.PollWithOptions, int64) {
	p := s.createRandomPoll()
	wonOption := p.Options[0]

	ctx := context.Background()
	now := time.Now()

	// create votes
	var counter int64
	targetUserID := gofakeit.Int64()
	for _, op := range p.Options {
		for i := 0; i < 5; i++ {
			userID := gofakeit.Int64()
			// set target user id for won option
			if op.ID == wonOption.ID && i == 2 {
				userID = targetUserID
			}

			_, err := s.db.CreateVote(ctx, swagger.CreateVote{
				OptionID: op.ID,
				PollID:   p.ID,
				UserID:   userID,
			}, now.Unix()+counter)
			s.Require().NoError(err)

			counter++
		}
	}

	// set actual outcome for first option
	_, err := s.db.UpdateOption(ctx, p.ID, wonOption.ID, swagger.UpdateOption{
		IsActualOutcome: helpers.Ptr[bool](true),
	}, now)
	s.Require().NoError(err)

	s.Require().NoError(s.db.CalculateStatistics(ctx, p.ID))

	return p, targetUserID
}

func (s *TelegramServiceSuite) checkStatisticsForUserPollResultTest(userName, text string) {
	// statistics should be:
	// totalVotes = 15
	// votesForVonOption = 5
	// prozentOfWonVotesBehind = 2/5 = 40%
	// prozentOfAllVotesBehind = (2+10)/15 = 12/15 = 80%
	for _, substring := range []string{
		userName,
		"predicted",
		"Out of 15",
		"only 5",
		"40%",
		"80%",
	} {
		s.Require().Contains(text, substring)
	}
}
