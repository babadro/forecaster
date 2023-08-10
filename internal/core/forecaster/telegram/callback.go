package telegram

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
)

type callbackHandlerFunc func(ctx context.Context, callbackData string) (tgbotapi.Chattable, string, error)

type pageService[T proto.Message] interface {
	Render(ctx context.Context, request T) (tgbotapi.Chattable, string, error)
}

func NewCallbackHandlers(db models.DB, bot models.TgBot) [256]callbackHandlerFunc {
	var handlers [256]callbackHandlerFunc

	defaultHandler := func(ctx context.Context, callbackData string) (tgbotapi.Chattable, string, error) {
		return nil, "", fmt.Errorf("handler for route %d is not implemented", callbackData[0])
	}

	for i := range handlers {
		handlers[i] = defaultHandler
	}

	handlers[models.VotePreviewRoute] = unmarshalMiddleware(votepreview.VotePreview(db, bot))

}

func unmarshalMiddleware[T proto.Message](next pageService[T]) callbackHandlerFunc {
	return func(ctx context.Context, callbackData string) (tgbotapi.Chattable, string, error) {
		var req T
		if err := helpers.UnmarshalCallbackData(callbackData, req); err != nil {
			return nil, "", fmt.Errorf("can't unmarshal callback data: %w", err)
		}

		return next.Render(ctx, req)
	}
}
