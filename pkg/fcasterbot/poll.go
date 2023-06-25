package fcasterbot

import "time"

type Poll struct {
	ID          int32
	Name        string
	Description string
	Start, End  time.Time
	Options     []Option
	UpdatedAt   time.Time
}

type Option struct {
	ID          int32
	Name        string
	Description string
}
