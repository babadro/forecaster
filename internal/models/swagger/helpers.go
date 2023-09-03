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
