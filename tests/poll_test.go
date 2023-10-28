package polls_test

import (
	"context"
	"github.com/babadro/forecaster/internal/models"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"testing"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
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
		require.Equal(t, createInput.SeriesID, got.SeriesID)
		require.Equal(t, createInput.TelegramUserID, got.TelegramUserID)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, createInput.Start, got.Start)
		timeRoundEqual(t, createInput.Finish, got.Finish)
	}

	checkReadRes := func(t *testing.T, got swagger.PollWithOptions) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)
		require.Equal(t, createInput.SeriesID, got.SeriesID)
		require.Equal(t, createInput.TelegramUserID, got.TelegramUserID)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, createInput.Start, got.Start)
		timeRoundEqual(t, createInput.Finish, got.Finish)
	}

	updateInput := randomModel[swagger.UpdatePoll](s.T())
	updateInput.SeriesID = helpers.Ptr[int32](0)
	updateInput.Status = randomPollStatus()

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Poll) {
		require.Equal(t, id, got.ID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)
		require.Equal(t, *updateInput.SeriesID, got.SeriesID)
		require.Equal(t, *updateInput.TelegramUserID, got.TelegramUserID)

		require.NotZero(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, *updateInput.Start, got.Start)
		timeRoundEqual(t, *updateInput.Finish, got.Finish)
	}

	testInput := crudEndpointTestInput[
		swagger.CreatePoll, swagger.Poll, swagger.PollWithOptions, swagger.UpdatePoll, swagger.Poll,
	]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkCreateRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "polls",
	}

	testCRUDEndpoints[swagger.CreatePoll, swagger.Poll](s.T(), testInput, s.apiAddr)
}

func (s *APITestSuite) TestPolls_Options() {
	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = 0
	pollInput.Start = strfmt.DateTime(time.Now())
	pollInput.Finish = strfmt.DateTime(time.Now().Add(time.Hour))

	poll := create[swagger.CreatePoll, swagger.Poll](s.T(), pollInput, s.url("polls"))

	createdOptions := make([]*swagger.Option, 3)
	for i := range createdOptions {
		optionInput := randomModel[swagger.CreateOption](s.T())
		optionInput.PollID = poll.ID

		createdOption := create[swagger.CreateOption, swagger.Option](s.T(), optionInput, s.url("options"))
		createdOptions[i] = &createdOption
	}

	// set actual outcome to the first option
	updateInput := swagger.UpdateOption{
		IsActualOutcome: helpers.Ptr[bool](true),
	}
	updatedOption := update[swagger.UpdateOption, swagger.Option](
		s.T(), updateInput, optionURLWithIDs(s.apiAddr, poll.ID, createdOptions[0].ID),
	)
	createdOptions[0] = &updatedOption

	gotPollOptions := read[swagger.PollWithOptions](s.T(), urlWithID(s.apiAddr, "polls", poll.ID)).Options

	NewGomegaWithT(s.T()).Expect(gotPollOptions).To(ConsistOf(createdOptions))
}

func (s *APITestSuite) TestPolls_Series() {
	series := create[swagger.CreateSeries, swagger.Series](
		s.T(), randomModel[swagger.CreateSeries](s.T()), s.url("series"),
	)

	pollInput := randomModel[swagger.CreatePoll](s.T())
	pollInput.SeriesID = series.ID

	poll := create[swagger.CreatePoll, swagger.Poll](s.T(), pollInput, s.url("polls"))

	require.Equal(s.T(), series.ID, poll.SeriesID)
}

// check that popularity is in returning poll and in updated poll model...
func (s *APITestSuite) TestPolls_popularity() {
	createInput := randomModel[swagger.CreatePoll](s.T())
	createInput.SeriesID = 0

	poll := create[swagger.CreatePoll, swagger.Poll](
		s.T(), createInput, s.url("polls"),
	)

	ctx := context.Background()

	popularity := rand.Int31()

	_, err := s.testDB.DB.Exec(ctx, "UPDATE forecaster.polls SET popularity = $1 WHERE id = $2", popularity, poll.ID)
	require.NoError(s.T(), err)

	gotPoll := read[swagger.PollWithOptions](s.T(), urlWithID(s.apiAddr, "polls", poll.ID))

	s.Require().Equal(popularity, gotPoll.Popularity)

	updateInput := randomModel[swagger.UpdatePoll](s.T())
	updateInput.SeriesID = helpers.Ptr[int32](0)
	updateInput.Status = randomPollStatus()

	gotUpdateResult := update[swagger.UpdatePoll, swagger.Poll](
		s.T(), updateInput, urlWithID(s.apiAddr, "polls", poll.ID),
	)

	s.Require().Equal(popularity, gotUpdateResult.Popularity)
}

func randomPollStatus() swagger.PollStatus {
	return swagger.PollStatus(gofakeit.RandomInt([]int{
		int(models.UnknownPollStatus), int(models.DraftPollStatus),
		int(models.ActivePollStatus), int(models.FinishedPollStatus),
	}))
}
