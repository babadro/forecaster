package fcasterbot

import "time"

type Poll struct {
	ID            int32
	Title         string
	Description   string
	Start, Finish time.Time
	Options       []Option
	UpdatedAt     time.Time
}

type Option struct {
	ID          int32
	PollID      int32
	Title       string
	Description string
}
