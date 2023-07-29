package polls_tests

import (
	"testing"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestOptions() {
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0

	poll := create[swagger.CreatePoll, swagger.Poll](
		s.T(), pollInput, "polls",
	)

	createInput := randomModel[swagger.CreateOption](s.T())
	createInput.PollID = poll.ID

	gotCreateResult := create[swagger.CreateOption, swagger.Option](s.T(), createInput, "options")

	optionID := gotCreateResult.ID

	checkCreateRes := func(t *testing.T, got swagger.Option) {
		require.NotZero(t, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.UpdatedAt)
	}

	checkCreateRes(s.T(), gotCreateResult)

	updateInput := randomModel[swagger.UpdateOption](s.T())

	gotUpdateResult := update[swagger.UpdateOption, swagger.Option](
		s.T(), updateInput, "options", optionID,
	)

	checkUpdateRes := func(t *testing.T, got swagger.Option) {
		require.Equal(t, optionID, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		timeRoundEqualNow(t, got.UpdatedAt)
	}

	checkUpdateRes(s.T(), gotUpdateResult)

	deleteOp(s.T(), "options", gotCreateResult.ID)
}
