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

	p, targetUserID, _ := s.setupForUserPollResultTest(false)

	update := startShowPoll(p.ID, targetUserID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)

	pollButtons := s.buttonsFromMarkup(pollMsg.ReplyMarkup)

	s.Require().Len(pollButtons, len(p.Options)+2) // +2 for "Show results" and navigation button
	showResultsButton := pollButtons[len(pollButtons)-2]
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

	p, targetUserID, _ := s.setupForUserPollResultTest(false)

	// send /start showuserres_<poll_id>_<user_id> command
	update := startShowUserRes(p.ID, targetUserID)
	s.sendMessage(update)

	userPollResultMsg := s.asMessage(sentMsg)
	s.checkStatisticsForUserPollResultTest(update.Message.From.UserName, userPollResultMsg.Text)

	s.Require().NotContains(userPollResultMsg.Text, "you") // in case of /start command the message is for 3rd person
}

func (s *TelegramServiceSuite) TestUserPollResult_user_voted_last() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	// 4 means that user voted last among right voters
	p, targetUserID, _ := s.setupForUserPollResultTest(true)

	// send /start showuserres_<poll_id>_<user_id> command
	update := startShowUserRes(p.ID, targetUserID)
	s.sendMessage(update)

	userPollResultMsg := s.asMessage(sentMsg)

	// if we calculated his position among right voters, then he should be on last position with 0% of votes behind
	// we should verify that we don't show 0% for him, it doesn't make sense
	s.Require().NotContains(userPollResultMsg.Text, "0%")
}

func (s *TelegramServiceSuite) TestUserPollResult_wrong_voted_user() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	p, _, votes := s.setupForUserPollResultTest(false)
	wonOption, idx := swagger.GetOutcome(p.Options)
	s.Require().NotEqual(-1, idx)

	targetUserID := findItemByCriteria(s, votes, func(v swagger.Vote) bool {
		return v.OptionID != wonOption.ID
	}).UserID

	// send /start showuserres_<poll_id>_<user_id> command
	update := startShowPoll(p.ID, targetUserID)
	s.sendMessage(update)

	pollMsg := s.asMessage(sentMsg)
	pollButtons := s.buttonsFromMarkup(pollMsg.ReplyMarkup)
	showResultsButton := pollButtons[len(pollButtons)-2]
	_ = s.sendCallback(showResultsButton, targetUserID)

	userPollResultMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(userPollResultMsg.Text, "didn't quite pan out this time")
}

func (s *TelegramServiceSuite) TestBackToPollButton() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	p, targetUserID, _ := s.setupForUserPollResultTest(false)

	// send /start showuserres_<poll_id>_<user_id> command
	update := startShowUserRes(p.ID, targetUserID)
	s.sendMessage(update)

	userPollResultMsg := s.asMessage(sentMsg)

	// verify the "Back to poll" button
	userPollResultButtons := s.buttonsFromMarkup(userPollResultMsg.ReplyMarkup)
	s.Require().NotEmpty(userPollResultButtons)

	backToPollButton := userPollResultButtons[len(userPollResultButtons)-1]
	s.Require().Contains(backToPollButton.Text, "Back")
	_ = s.sendCallback(backToPollButton, targetUserID)

	// verify that we are back to the same poll
	pollMsg := s.asEditMessage(sentMsg)
	s.Require().Contains(pollMsg.Text, p.Title)
}

func (s *TelegramServiceSuite) setupForUserPollResultTest(
	setOnLastPosition bool,
) (swagger.PollWithOptions, int64, []swagger.Vote) {
	p := s.createRandomPoll(time.Now())
	wonOption := p.Options[0]

	ctx := context.Background()
	now := time.Now()

	// create votes
	var counter int64

	targetUserID := gofakeit.Int64()

	var votes []swagger.Vote

	for _, op := range p.Options {
		for i := 0; i < 5; i++ {
			userID := gofakeit.Int64()
			// set target user id for won option
			votePosition := 2
			if setOnLastPosition {
				votePosition = 4
			}

			if op.ID == wonOption.ID && i == votePosition {
				userID = targetUserID
			}

			v, err := s.db.CreateVote(ctx, swagger.CreateVote{
				OptionID: op.ID,
				PollID:   p.ID,
				UserID:   userID,
			}, now.Unix()+counter)
			s.Require().NoError(err)

			votes = append(votes, v)

			counter++
		}
	}

	// set actual outcome for first option
	_, err := s.db.UpdateOption(ctx, p.ID, wonOption.ID, swagger.UpdateOption{
		IsActualOutcome: helpers.Ptr[bool](true),
	}, now)
	s.Require().NoError(err)

	s.Require().NoError(s.db.CalculateStatistics(ctx, p.ID))

	p, err = s.db.GetPollByID(ctx, p.ID)
	s.Require().NoError(err)

	return p, targetUserID, votes
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
