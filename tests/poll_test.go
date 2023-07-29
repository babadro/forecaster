package polls_tests

import (
	"testing"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"

	. "github.com/onsi/gomega"
)

func (s *APITestSuite) TestPolls() {
	createInput := randomModel[swagger.CreatePoll](s.T())
	createInput.SeriesID = 0

	checkCreateRes := func(t *testing.T, got swagger.Poll) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, createInput.Start, got.Start)
		timeRoundEqual(t, createInput.Finish, got.Finish)
	}

	checkReadRes := func(t *testing.T, got swagger.PollWithOptions) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, createInput.Start, got.Start)
		timeRoundEqual(t, createInput.Finish, got.Finish)
	}

	updateInput := randomModel[swagger.UpdatePoll](s.T())
	updateInput.SeriesID = helpers.Ptr[int32](0)

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Poll) {
		require.Equal(t, id, got.ID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		require.NotZero(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, *updateInput.Start, got.Start)
		timeRoundEqual(t, *updateInput.Finish, got.Finish)
	}

	testInput := crudEndpointTestInput[swagger.CreatePoll, swagger.Poll, swagger.PollWithOptions, swagger.UpdatePoll, swagger.Poll]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkCreateRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "polls",
	}

	testCRUDEndpoints[swagger.CreatePoll, swagger.Poll](s.T(), testInput)
}

func (s *APITestSuite) TestPolls_Options() {
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0

	poll := create[swagger.CreatePoll, swagger.Poll](s.T(), pollInput, "polls")

	createdOptions := make([]*swagger.Option, 3)
	for i := range createdOptions {
		optionInput := randomModel[swagger.CreateOption](s.T())
		optionInput.PollID = poll.ID

		createdOption := create[swagger.CreateOption, swagger.Option](s.T(), optionInput, "options")
		createdOptions[i] = &createdOption
	}

	gotPollOptions := read[swagger.PollWithOptions](s.T(), "polls", poll.ID).Options

	NewGomegaWithT(s.T()).Expect(gotPollOptions).To(ConsistOf(createdOptions))
}
