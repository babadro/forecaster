package polls

import (
	"context"
	"errors"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto2.Message, *polls.Polls) {
	v := new(polls.Polls)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {

}

func (s *Service) RenderCallback(
	ctx context.Context, req *polls.Polls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {

}

const pageSize = 10

func (s *Service) render(
	ctx context.Context, currentPage int32, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	polls, totalCount, err := s.db.GetPolls(ctx, currentPage, pageSize)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, "", nil // TODO: render not found result
		}

		return nil, "", fmt.Errorf("unable to get polls: %s", err.Error())
	}

	return nil, "", nil
}

func txtMsg(p []swagger.Poll) string {
	var sb render.StringBuilder
	for i, poll := range p {
		sb.Printf("%d. %s\n", i+1, poll.Title)
	}

	return sb.String()
}
