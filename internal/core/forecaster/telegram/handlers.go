package telegram

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecast"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecasts"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"

	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/vote"
	votepreviewproto "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

type handlerFunc func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error)

type pageService[T proto.Message] interface {
	RenderCallback(ctx context.Context, req T, upd tgbotapi.Update) (tgbotapi.Chattable, string, error)
	// NewRequest returns proto message and request for RenderCallback
	// Under the hood both of returned values are the same pointer to the same struct
	NewRequest() (proto.Message, T)
}

func newCallbackHandlers(svc pageServices) [256]handlerFunc {
	var handlers [256]handlerFunc

	handlers[models.MainPageRoute] = unmarshalMiddleware[*mainpage.MainPage](svc.main)
	handlers[models.VotePreviewRoute] = unmarshalMiddleware[*votepreviewproto.VotePreview](svc.votePreview)
	handlers[models.VoteRoute] = unmarshalMiddleware[*vote.Vote](svc.vote)
	handlers[models.PollRoute] = unmarshalMiddleware[*poll.Poll](svc.poll)
	handlers[models.UserPollResultRoute] = unmarshalMiddleware[*userpollresult.UserPollResult](svc.userPollResult)
	handlers[models.PollsRoute] = unmarshalMiddleware[*polls.Polls](svc.polls)
	handlers[models.ForecastsRoute] = unmarshalMiddleware[*forecasts.Forecasts](svc.forecasts)
	handlers[models.ForecastRoute] = unmarshalMiddleware[*forecast.Forecast](svc.forecast)
	handlers[models.EditPollRoute] = unmarshalMiddleware[*editpoll.EditPoll](svc.editPoll)
	handlers[models.EditFieldRoute] = unmarshalMiddleware[*editfield.EditField](svc.editField)
	handlers[models.MyPollsRoute] = unmarshalMiddleware[*mypolls.MyPolls](svc.myPolls)

	defaultHandler := func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		return nil, "", fmt.Errorf("handler for route %d is not implemented", upd.CallbackQuery.Data[0])
	}

	for i := range handlers {
		if handlers[i] != nil {
			handlers[i] = chainMiddlewares(handlers[i],
				validateCallbackInput,
			)
		} else {
			handlers[i] = defaultHandler
		}
	}

	return handlers
}

func unmarshalMiddleware[T proto.Message](next pageService[T]) handlerFunc {
	return func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		requestAsProtoMessage, requestAsStruct := next.NewRequest()

		if err := proto2.UnmarshalCallbackData(upd.CallbackQuery.Data, requestAsProtoMessage); err != nil {
			return nil, "", fmt.Errorf("can't unmarshal callback data: %w", err)
		}

		return next.RenderCallback(ctx, requestAsStruct, upd)
	}
}

type middleware func(next handlerFunc) handlerFunc

func chainMiddlewares(mainHandler handlerFunc, middlewares ...middleware) handlerFunc {
	h := mainHandler

	for i := range middlewares {
		h = middlewares[len(middlewares)-1-i](h)
	}

	return h
}

func validateCallbackInput(next handlerFunc) handlerFunc {
	return func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		if upd.CallbackQuery == nil {
			return nil, "", fmt.Errorf("callback query is nil")
		}

		if upd.CallbackQuery.Message == nil {
			return nil, "", fmt.Errorf("callbackQuery.message is nil")
		}

		if upd.CallbackQuery.Message.Chat == nil {
			return nil, "", fmt.Errorf("callbackQuery.chat is nil")
		}

		if upd.CallbackQuery.From == nil {
			return nil, "", fmt.Errorf("callbackQuery.from is nil")
		}

		return next(ctx, upd)
	}
}

func validateCommandInput(next handlerFunc) handlerFunc {
	return func(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
		if upd.Message.Chat == nil {
			return nil, "", fmt.Errorf("chat is nil")
		}

		return next(ctx, upd)
	}
}
