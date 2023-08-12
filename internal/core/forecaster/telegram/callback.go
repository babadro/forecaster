package telegram

import (
	"context"
	"fmt"

	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/vote"
	votepreviewproto "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
)

type callbackHandlerFunc func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error)

type pageService[T proto.Message] interface {
	RenderCallback(ctx context.Context, req T, upd tgbotapi.Update) (tgbotapi.Chattable, string, error)
}

func newCallbackHandlers(svc pageServices) [256]callbackHandlerFunc {
	var handlers [256]callbackHandlerFunc

	defaultHandler := func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		return nil, "", fmt.Errorf("handler for route %d is not implemented", upd.CallbackQuery.Data[0])
	}

	for i := range handlers {
		handlers[i] = defaultHandler
	}

	handlers[models.VotePreviewRoute] = unmarshalMiddleware[*votepreviewproto.VotePreview](svc.votePreview)
	handlers[models.VoteRoute] = unmarshalMiddleware[*vote.Vote](svc.vote)

	return handlers
}

func unmarshalMiddleware[T proto.Message](next pageService[T]) callbackHandlerFunc {
	return func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		var req T
		if err := proto2.UnmarshalCallbackData(upd.CallbackQuery.Data, req); err != nil {
			return nil, "", fmt.Errorf("can't unmarshal callback data: %w", err)
		}

		return next.RenderCallback(ctx, req, upd)
	}
}
