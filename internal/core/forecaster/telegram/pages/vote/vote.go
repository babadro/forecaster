package vote

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	poll2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/helpers"
	proto2 "google.golang.org/protobuf/proto"

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

func (s *Service) NewRequest() (proto2.Message, *vote.Vote) {
	v := new(vote.Vote)
	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, vote *vote.Vote, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
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

	chatID, messageID := upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID
	if err != nil {
		// You have already voted for this option
		if errors.Is(err, domain.ErrVoteWithSameOptionAlreadyExists) {
			return tryCreateMessage(chatID, messageID, vote.PollId, vote.GetReferrerForecastsPage(),
				"You have already voted for this option", "",
			)
		}

		return nil,
			"Sorry, something went wrong, can't vote right now",
			fmt.Errorf("vote: unable to create vote: %s", err.Error())
	}

	return tryCreateMessage(chatID, messageID, vote.PollId, vote.GetReferrerForecastsPage(),
		"Success!", "Vote was successful, but I cant get you back to poll due to the error")
}

func tryCreateMessage(
	chatID int64, msgID int, pollID *int32, referrerForecastsPage int32, successText, failText string,
) (tgbotapi.Chattable, string, error) {
	pollMsg := &poll2.Poll{PollId: pollID}
	if referrerForecastsPage > 0 {
		pollMsg.ReferrerForecastsPage = helpers.Ptr(referrerForecastsPage)
	}

	callBackData, err := proto.MarshalCallbackData(models.PollRoute, pollMsg)
	if err != nil {
		return nil, failText,
			fmt.Errorf("vote: unable to marshal callback data: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(
		chatID, msgID, successText,
		render.Keyboard(tgbotapi.InlineKeyboardButton{
			Text:         "Back to poll",
			CallbackData: callBackData,
		}),
	), "", nil
}
