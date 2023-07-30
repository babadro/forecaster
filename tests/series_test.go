package pollstests

import (
	"testing"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/stretchr/testify/require"
)

func (s *APITestSuite) TestSeries() {
	createInput := randomModel[swagger.CreateSeries](s.T())

	checkReadRes := func(t *testing.T, got swagger.Series) {
		require.NotZero(t, got.ID)

		require.Equal(t, createInput.Description, got.Description)
		require.Equal(t, createInput.Title, got.Title)

		timeRoundEqualNow(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)
	}

	updateInput := randomModel[swagger.UpdateSeries](s.T())

	checkUpdateRes := func(t *testing.T, id int32, got swagger.Series) {
		require.Equal(t, id, got.ID)

		require.Equal(t, *updateInput.Description, got.Description)
		require.Equal(t, *updateInput.Title, got.Title)

		require.NotZero(t, got.CreatedAt)
		timeRoundEqualNow(t, got.UpdatedAt)
	}

	testInput := crudEndpointTestInput[
		swagger.CreateSeries, swagger.Series, swagger.Series, swagger.UpdateSeries, swagger.Series,
	]{
		createInput:    createInput,
		updateInput:    updateInput,
		checkCreateRes: checkReadRes,
		checkReadRes:   checkReadRes,
		checkUpdateRes: checkUpdateRes,
		path:           "series",
	}

	testCRUDEndpoints[swagger.CreateSeries, swagger.Series](s.T(), testInput, s.apiAddr)
}
