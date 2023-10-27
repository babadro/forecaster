package models

import "github.com/babadro/forecaster/internal/models/swagger"

func PollsIDs(pollsArr []swagger.Poll) []int32 {
	ids := make([]int32, len(pollsArr))
	for i, p := range pollsArr {
		ids[i] = p.ID
	}

	return ids
}
