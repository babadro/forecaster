package telegram_test

import (
	"context"

	"github.com/babadro/forecaster/internal/domain"
)

// open create poll page and click back button...
func (s *TelegramServiceSuite) TestDeletePoll() {
	userID := randomPositiveInt64()

	var sentMsg interface{}

	s.mockTelegramSender(&sentMsg)

	editPoll := s.createPollAndGoToEditPollPage(userID, &sentMsg)

	deleteButton := s.findButtonByLowerText("delete poll", editPoll.ReplyMarkup)

	s.sendCallback(deleteButton, userID)

	pollsArr, _, err := s.db.GetPolls(context.Background(), 0, 1)
	s.Require().NoError(err)
	s.Require().Len(pollsArr, 1)
	p := pollsArr[0]

	deleteConfirmation := s.asEditMessage(sentMsg)

	s.Require().Contains(deleteConfirmation.Text, p.Title)

	// verify that poll was not deleted yet
	_, err = s.db.GetPollByID(context.Background(), p.ID)
	s.Require().NoError(err)

	deleteButton = s.findButtonByLowerText("delete", deleteConfirmation.ReplyMarkup)

	s.sendCallback(deleteButton, userID)

	// check that poll was deleted this time
	_, err = s.db.GetPollByID(context.Background(), p.ID)
	s.Require().ErrorIs(err, domain.ErrNotFound)
}
