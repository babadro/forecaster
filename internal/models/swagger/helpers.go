package swagger

func FindOptionByID(options []*Option, id int16) (*Option, int) {
	for i, op := range options {
		if op.ID == id {
			return op, i
		}
	}

	return nil, -1
}

func GetOutcome(options []*Option) (Option, int) {
	for i, op := range options {
		if op.IsActualOutcome {
			return *op, i
		}
	}

	return Option{}, -1
}

func HasOutcome(options []*Option) bool {
	for _, op := range options {
		if op.IsActualOutcome {
			return true
		}
	}

	return false
}

func HasVotes(options []*Option) bool {
	for _, op := range options {
		if op.TotalVotes > 0 {
			return true
		}
	}

	return false
}

func TotalVotes(options []*Option) int32 {
	var totalVotes int32

	for _, op := range options {
		totalVotes += op.TotalVotes
	}

	return totalVotes
}

// todo unit test for fields (random)
func MergePolls(pollWithOptions PollWithOptions, p Poll) PollWithOptions {
	pollWithOptions.ID = p.ID
	pollWithOptions.SeriesID = p.SeriesID
	pollWithOptions.Title = p.Title
	pollWithOptions.Description = p.Description
	pollWithOptions.TelegramUserID = p.TelegramUserID
	pollWithOptions.Start = p.Start
	pollWithOptions.Finish = p.Finish
	pollWithOptions.CreatedAt = p.CreatedAt
	pollWithOptions.UpdatedAt = p.UpdatedAt

	return pollWithOptions
}
