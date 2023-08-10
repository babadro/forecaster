package votepreview

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func VotePreview(ctx context.Context, callbackData string, scope models.Scope) (tgbotapi.Chattable, string, error) {
	var req votepreview2.VotePreview
	if err := helpers.UnmarshalCallbackData(callbackData, &req); err != nil {
		return nil, "", fmt.Errorf("can't unmarshal votePreview callback data: %w", err)
	}
}
