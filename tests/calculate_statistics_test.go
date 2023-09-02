package polls_test

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APITestSuite) TestCalculateStatistics() {
	ctx := context.Background()
	// create control poll, it should not be affected by calculate statistics for other poll
	controlPoll := s.createPollWithOptions(2)
	_ = update[swagger.UpdateOption, swagger.Option](
		s.T(), swagger.UpdateOption{
			IsActualOutcome: helpers.Ptr[bool](true),
		}, optionURLWithIDs(s.apiAddr, controlPoll.ID, controlPoll.Options[0].ID),
	)

	controlVote, err := s.forecasterDB.CreateVote(ctx, swagger.CreateVote{
		OptionID: controlPoll.Options[0].ID,
		PollID:   controlPoll.ID,
		UserID:   gofakeit.Int64(),
	}, time.Time(controlPoll.Start).Add(time.Second).Unix())
	s.Require().NoError(err)

	// create poll that we will calculate statistics for
	p := s.createPollWithOptions(2)

	// set actual outcome for first option
	_ = update[swagger.UpdateOption, swagger.Option](
		s.T(), swagger.UpdateOption{
			IsActualOutcome: helpers.Ptr[bool](true),
		}, optionURLWithIDs(s.apiAddr, p.ID, p.Options[0].ID),
	)

	winOptionID := p.Options[0].ID

	type vote struct {
		timestamp int64
		userID    int64
		optionID  int16
	}

	votes := make([]vote, 10)
	for i := range votes {
		votes[i].timestamp = time.Time(p.Start).Unix() + int64(i)
		votes[i].userID = gofakeit.Int64()
	}

	gofakeit.ShuffleAnySlice(votes)

	for i := range votes {
		// 5 votes for first option and 5 votes for second option
		votes[i].optionID = p.Options[0].ID
		if i > 4 {
			votes[i].optionID = p.Options[1].ID
		}
	}

	for _, v := range votes {
		_, err := s.forecasterDB.CreateVote(ctx, swagger.CreateVote{
			OptionID: v.optionID,
			PollID:   p.ID,
			UserID:   v.userID,
		}, v.timestamp)
		s.Require().NoError(err)
	}

	// calculate statistics
	post(s.T(), urlWithID(s.apiAddr, "calculate-statistics", p.ID), http.StatusNoContent)

	// check statistics
	p = read[swagger.PollWithOptions](s.T(), urlWithID(s.apiAddr, "polls", p.ID))

	op, idx := swagger.FindOptionByID(p.Options, winOptionID)
	s.Require().NotEqual(-1, idx)

	s.Require().Equal(int32(5), op.TotalVotes)

	sort.Slice(votes, func(i, j int) bool {
		return votes[i].timestamp < votes[j].timestamp
	})

	// check votes positions
	position := int32(1)
	for _, v := range votes {
		dbVote, err := s.forecasterDB.GetUserVote(ctx, v.userID, p.ID)
		s.Require().NoError(err)

		// check only votes for win option
		if v.optionID == winOptionID {
			s.Require().Equal(position, dbVote.Position)
			position++
		} else {
			s.Require().Zero(dbVote.Position)
		}
	}

	// check statistics for control poll
	// control poll should not be affected by calculate statistics for other poll
	controlPoll = read[swagger.PollWithOptions](s.T(), urlWithID(s.apiAddr, "polls", controlPoll.ID))
	for _, opt := range controlPoll.Options {
		s.Require().Zero(opt.TotalVotes)
	}

	// since control poll statistics was not calculated, vote position should be zero
	dbVote, err := s.forecasterDB.GetUserVote(ctx, controlVote.UserID, controlVote.PollID)
	s.Require().NoError(err)
	s.Require().Zero(dbVote.Position)
}
