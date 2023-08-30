package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound                        = errors.New("object not found")
	ErrVoteWithSameOptionAlreadyExists = errors.New("vote with the same option already exists")
)

type OptionWithOutcomeFlagAlreadyExistsError struct {
	PollID   int32
	OptionID int16
}

func (e OptionWithOutcomeFlagAlreadyExistsError) Error() string {
	return fmt.Sprintf("option with IsActualOutcome=true already exists; pollID: %d, optionID: %d", e.PollID, e.OptionID)
}
