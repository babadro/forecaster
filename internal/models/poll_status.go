package models

import (
	"database/sql/driver"
	"fmt"
)

type PollStatus int32

const (
	UnknownPollStatus PollStatus = iota
	DraftPollStatus
	ActivePollStatus
	FinishedPollStatus
)

const unknown = "unknown"

func (u *PollStatus) String() string {
	switch *u {
	case UnknownPollStatus:
		return unknown
	case DraftPollStatus:
		return "draft"
	case ActivePollStatus:
		return "active"
	case FinishedPollStatus:
		return "finished"
	default:
		return unknown
	}
}

func (u *PollStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to unmarshal PollStatus value: %v", value)
	}

	switch str {
	case unknown:
		*u = UnknownPollStatus
	case "draft":
		*u = DraftPollStatus
	case "active":
		*u = ActivePollStatus
	case "finished":
		*u = FinishedPollStatus
	default:
		return fmt.Errorf("unknown PollStatus value: %v", value)
	}

	return nil
}

func (u *PollStatus) Value() (driver.Value, error) {
	return u.String(), nil
}
