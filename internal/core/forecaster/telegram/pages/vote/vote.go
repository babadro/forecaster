package vote

import (
	"context"
	"fmt"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	poll2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/vote"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RenderCallback(ctx context.Context, vote *vote.Vote, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	poll, err := s.db.GetPollByID(ctx, *vote.PollId)
	if err != nil {
		return nil,
			fmt.Sprintf("Oops, can't find poll with id %d", *vote.PollId),
			fmt.Errorf("vote: unable to get poll by id: %s", err.Error())
	}

	if time.Now().After(time.Time(poll.Finish)) {
		return nil,
			"Sorry, this poll is expired",
			fmt.Errorf("vote: poll is expired")
	}

	_, err = s.db.CreateVote(ctx, swagger.CreateVote{
		PollID:   *vote.PollId,
		OptionID: int16(*vote.OptionId),
		UserID:   upd.CallbackQuery.From.ID,
	}, time.Now().Unix())

	if err != nil {
		return nil,
			"Sorry, something went wrong, can't vote right now",
			fmt.Errorf("vote: unable to create vote: %s", err.Error())
	}

	callBackData, err := proto.MarshalCallbackData(models.PollRoute, &poll2.Poll{
		PollId: vote.PollId,
	})
	if err != nil {
		return nil,
			"Vote was successful, but I cant get you back to poll due to the error",
			fmt.Errorf("vote: unable to marshal callback data: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(
		upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID,
		"Success!",
		render.Keyboard(tgbotapi.InlineKeyboardButton{
			Text:         "Back to poll",
			CallbackData: callBackData,
		})), "", nil
}
