package dbwrapper

import (
	"context"
	"errors"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/models/swagger"
)

type Wrapper struct {
	db models.DB
}

func New(db models.DB) Wrapper {
	return Wrapper{db: db}
}

func (w Wrapper) GetPollWithOptionsByID(ctx context.Context, pollID int32) (swagger.PollWithOptions, string, error) {
	p, err := w.db.GetPollWithOptionsByID(ctx, pollID)

	if err != nil {
		return swagger.PollWithOptions{},
			fmt.Sprintf("oops, can't find poll with id %d", pollID),
			fmt.Errorf("unable to get pollWithOptions by id: %s", err.Error())
	}

	return p, "", nil
}

func (w Wrapper) GetPollByID(ctx context.Context, pollID int32) (swagger.Poll, string, error) {
	p, err := w.db.GetPollByID(ctx, pollID)

	if err != nil {
		return swagger.Poll{},
			fmt.Sprintf("oops, can't find poll with id %d", pollID),
			fmt.Errorf("unable to get poll by id: %s", err.Error())
	}

	return p, "", nil
}

func (w Wrapper) GetUserVote(ctx context.Context, userID int64, pollID int32) (swagger.Vote, bool, error) {
	lastVote, err := w.db.GetUserVote(ctx, userID, pollID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return swagger.Vote{}, false,
				fmt.Errorf("unable to get last vote: %s", err.Error())
		}

		return swagger.Vote{}, false, nil
	}

	return lastVote, true, nil
}
