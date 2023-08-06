package polls_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestOptions() {
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0

	poll := create[swagger.CreatePoll, swagger.Poll](
		s.T(), pollInput, s.url("polls"),
	)

	createInput := randomModel[swagger.CreateOption](s.T())
	createInput.PollID = poll.ID

	gotCreateResult := create[swagger.CreateOption, swagger.Option](s.T(), createInput, s.url("options"))

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
		s.T(), updateInput, optionURLWithIDs(s.apiAddr, poll.ID, optionID),
	)

	checkUpdateRes := func(t *testing.T, got swagger.Option) {
		require.Equal(t, optionID, got.ID)
		require.Equal(t, poll.ID, got.PollID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		timeRoundEqualNow(t, got.UpdatedAt)
	}

	checkUpdateRes(s.T(), gotUpdateResult)

	deleteOp(s.T(), optionURLWithIDs(s.apiAddr, poll.ID, optionID))
}

func (s *APITestSuite) TestOptions_pollDoesntExist() {
	createInput := randomModel[swagger.CreateOption](s.T())
	createInput.PollID = 999

	b, err := json.Marshal(createInput)
	s.Require().NoError(err)

	createResp, err := http.Post(
		s.url("options"),
		"application/json",
		bytes.NewReader(b))
	s.Require().NoError(err)

	defer func() { _ = createResp.Body.Close() }()

	s.Require().Equal(http.StatusBadRequest, createResp.StatusCode)
}
