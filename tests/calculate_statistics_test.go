package polls_test

import (
	"context"
	"net/http"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APITestSuite) TestCalculateStatistics() {
	p := s.createPollWithOptions(2)

	// set actual outcome for first option
	_ = update[swagger.UpdateOption, swagger.Option](
		s.T(), swagger.UpdateOption{
			IsActualOutcome: helpers.Ptr[bool](true),
		}, optionURLWithIDs(s.apiAddr, p.ID, p.Options[0].ID),
	)

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

	ctx := context.Background()

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
}
