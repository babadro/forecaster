package polls_tests

import (
	"testing"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestOptions() {
	poll := create[swagger.CreatePoll, swagger.Poll](
		s.T(), randomModel[swagger.CreatePoll](s.T()), "polls",
	)

	createInput := randomModel[swagger.CreateOption](s.T())
	createInput.PollID = poll.ID

	gotCreateResult := create[swagger.CreateOption, swagger.Option](s.T(), createInput, "options")

	checkCreateRes := func(t *testing.T, got swagger.Option) {
		require.NotZero(t, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)
	}

	checkCreateRes(s.T(), gotCreateResult)

	updateInput := randomModel[swagger.UpdateOption](s.T())

	gotUpdateResult := update[swagger.UpdateOption, swagger.Option](
		s.T(), updateInput, "options", gotCreateResult.ID,
	)

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Option) {
		require.Equal(t, id, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)
	}

	checkUpdateRes(s.T(), gotCreateResult.ID, gotUpdateResult)

	deleteOp(s.T(), "options", gotCreateResult.ID)
}
