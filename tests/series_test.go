package polls_tests

import (
	"testing"

	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TfestSeries() {
	createInput := swagger.CreateSeries{
		Description: "test desc",
		Title:       "test title",
	}

	checkReadRes := func(t *testing.T, got swagger.Series) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)
	}

	updateInput := swagger.UpdateSeries{
		Description: helpers.Ptr("updated desc"),
		Title:       helpers.Ptr("updated title"),
	}

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Series) {
		require.Equal(t, id, got.ID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		require.NotZero(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)
	}

	testInput := crudEndpointTestInput[swagger.CreateSeries, swagger.Series, swagger.UpdateSeries]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkReadRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "series",
	}

	testCRUDEndpoints[swagger.CreateSeries, swagger.Series](s.T(), testInput)
}
