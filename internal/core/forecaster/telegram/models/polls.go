package models

import "github.com/babadro/forecaster/internal/models"

type PollFilter struct {
	TelegramUserID Nullable[int64]
	Status         Nullable[models.PollStatus]
}

func NewPollFilter() PollFilter {
	return PollFilter{}
}

func (f PollFilter) WithTelegramUserID(id int64) PollFilter {
	f.TelegramUserID = NewNullable(id)

	return f
}

func (f PollFilter) WithStatus(status models.PollStatus) PollFilter {
	f.Status = NewNullable(status)

	return f
}

type PollSortType byte

const (
	DefaultPollSort PollSortType = iota
	CreatedAtPollSort
	PopularityPollSort
)

type PollSort struct {
	By  PollSortType
	Asc bool
}

type PollsFlags int32

const (
	FilterFinishedStatus PollsFlags = 1 << iota
)

func (f PollsFlags) IsSet(flag PollsFlags) bool {
	return f&flag != 0
}
