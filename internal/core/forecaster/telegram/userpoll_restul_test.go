package telegram_test

import (
	"context"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *TelegramServiceSuite) TestUserPollResult_callback_happy_path() {
	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	p := s.createRandomPoll()
	wonOption := p.Options[0]

	ctx := context.Background()
	now := time.Now()

	// create votes
	var counter int64
	targetUserID := gofakeit.Int64()
	for _, op := range p.Options {
		for i := 0; i < 5; i++ {
			userID := gofakeit.Int64()
			// set target user id for won option
			if op.ID == wonOption.ID && i == 2 {
				userID = targetUserID
			}

			_, err := s.db.CreateVote(ctx, swagger.CreateVote{
				OptionID: op.ID,
				PollID:   p.ID,
				UserID:   userID,
			}, now.Unix()+counter)
			s.Require().NoError(err)
		}

		counter++
	}

	// set actual outcome for first option
	_, err := s.db.UpdateOption(ctx, p.ID, wonOption.ID, swagger.UpdateOption{}, now)
	s.Require().NoError(err)

}
