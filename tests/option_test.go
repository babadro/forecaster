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

	checkReadRes := func(t *testing.T, got swagger.Option) {
		require.NotZero(t, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)
	}

	updateInput := randomModel[swagger.UpdateOption](s.T())

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Option) {
		require.Equal(t, id, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)
	}

	testInput := crudEndpointTestInput[swagger.CreateOption, swagger.Option, swagger.UpdateOption]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkReadRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "options",
	}

	testCRUDEndpoints[swagger.CreateOption, swagger.Option](s.T(), testInput)
}
