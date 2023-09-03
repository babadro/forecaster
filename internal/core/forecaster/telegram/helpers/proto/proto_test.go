package proto

import (
	"testing"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalCallbackData(t *testing.T) {
	showMyResultsData, err := MarshalCallbackData(models.UserPollResultRoute, &userpollresult.UserPollResult{
		UserId: helpers.Ptr[int64](999),
		PollId: helpers.Ptr[int32](475903447),
	})

	require.NoError(t, err)

	var unmarshalledData userpollresult.UserPollResult
	err = UnmarshalCallbackData(*showMyResultsData, &unmarshalledData)
	require.NoError(t, err)
}
