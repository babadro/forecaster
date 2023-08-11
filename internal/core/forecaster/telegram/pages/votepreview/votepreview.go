package votepreview

import (
	"context"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RenderCallback(ctx context.Context, callbackData *votepreview2.VotePreview) (tgbotapi.Chattable, string, error) {

	return nil, "", nil
}

func txtMsg(voted, expired bool, option swagger.Option, id int) string {
	var sb strings.Builder
	if voted {
		sb.WriteString("You have already voted.\n")
	}

	if expired {
		sb.WriteString("This poll is expired!\n")
	}

	if !(voted || expired) {
		sb.WriteString(fmt.Sprintf("Are you sure? Your choice is:\n %d %s", id+1, option.String()))
	} else {
		sb.WriteString(option.String())
	}

	return sb.String()
}
