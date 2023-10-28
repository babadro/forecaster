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

func (u *PollStatus) String() string {
	switch *u {
	case UnknownPollStatus:
		return "unknown"
	case DraftPollStatus:
		return "draft"
	case ActivePollStatus:
		return "active"
	case FinishedPollStatus:
		return "finished"
	default:
		return "unknown"
	}
}

// Implement the sql.Scanner interface for PollStatus
func (u *PollStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("Failed to unmarshal PollStatus value: %v", value)
	}

	switch str {
	case "unknown":
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

// Implement the driver.Valuer interface for PollStatus
func (u *PollStatus) Value() (driver.Value, error) {
	return u.String(), nil
}
