package polls_tests

import (
	"testing"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestPolls() {
	createInput := swagger.CreatePoll{
		Description: "test desc",
		Title:       "test title",
		Start:       strfmt.DateTime(time.Now().Add(time.Hour)),
		Finish:      strfmt.DateTime(time.Now().Add(time.Hour * 2)),
	}

	checkReadRes := func(t *testing.T, got swagger.Poll) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, createInput.Start, got.Start)
		timeRoundEqual(t, createInput.Finish, got.Finish)
	}

	updateInput := swagger.UpdatePoll{
		Description: helpers.Ptr("updated desc"),
		Title:       helpers.Ptr("updated title"),
		Start:       helpers.Ptr(strfmt.DateTime(time.Now().Add(time.Hour * 3))),
		Finish:      helpers.Ptr(strfmt.DateTime(time.Now().Add(time.Hour * 4))),
	}

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Poll) {
		require.Equal(t, id, got.ID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		require.NotZero(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)

		timeRoundEqual(t, *updateInput.Start, got.Start)
		timeRoundEqual(t, *updateInput.Finish, got.Finish)
	}

	testInput := crudEndpointTestInput[swagger.CreatePoll, swagger.Poll, swagger.UpdatePoll]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkReadRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "polls",
	}

	testCRUDEndpoints[swagger.CreatePoll, swagger.Poll](s.T(), testInput)
}
