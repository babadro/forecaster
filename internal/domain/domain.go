package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("object not found")
)

type OptionWithOutcomeFlagAlreadyExistsError struct {
	PollID   int32
	OptionID int16
}

func (e OptionWithOutcomeFlagAlreadyExistsError) Error() string {
	return fmt.Sprintf("option with ActualOutcome=true already exists; pollID: %d, optionID: %d", e.PollID, e.OptionID)
}
