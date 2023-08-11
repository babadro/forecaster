package votepreview

import (
	"context"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
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
