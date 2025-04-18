package votepreview

import (
	"context"
	"fmt"
	"time"

	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/vote"
	votepreview2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto.Message, *votepreview2.VotePreview) {
	v := new(votepreview2.VotePreview)
	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, votePreview *votepreview2.VotePreview, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	p, err := s.db.GetPollByID(ctx, *votePreview.PollId)
	if err != nil {
		return nil,
			fmt.Sprintf("oops, can't find poll with id %d", *votePreview.PollId),
			fmt.Errorf("unable to get poll by id: %s", err.Error())
	}

	op, idx := swagger.FindOptionByID(p.Options, int16(*votePreview.OptionId))
	if idx == -1 {
		return nil,
			"Sorry, something went wrong, I can't show this option right now",
			fmt.Errorf("votepreview: unable to find option with id %d", *votePreview.OptionId)
	}

	expired := time.Now().After(time.Time(p.Finish))

	msg := txtMsg(expired, *op)

	markup, err := keyboardMarkup(p.ID, votePreview.GetReferrerForecastsPage(), op.ID, expired)
	if err != nil {
		return nil,
			"Sorry, something went wrong, I can't show this option right now",
			fmt.Errorf("votepreview: unable to create keyboard markup: %s", err.Error())
	}

	origMsg := upd.CallbackQuery.Message

	return render.NewEditMessageTextWithKeyboard(origMsg.Chat.ID, origMsg.MessageID, msg, markup), "", nil
}

func txtMsg(expired bool, option swagger.Option) string {
	var sb render.StringBuilder

	if expired {
		sb.WriteStringLn("This poll is expired!")
	} else {
		sb.WriteStringLn("Vote for this option?")
	}

	sb.WriteStringLn(option.Title)
	sb.WriteString(option.Description)

	return sb.String()
}

func keyboardMarkup(
	pollID, referrerForecastsPage int32, optionID int16, voteNotAllowed bool,
) (tgbotapi.InlineKeyboardMarkup, error) {
	pollMsg := &poll.Poll{PollId: helpers.Ptr(pollID)}
	if referrerForecastsPage > 0 {
		pollMsg.ReferrerForecastsPage = helpers.Ptr(referrerForecastsPage)
	}

	backData, err := proto2.MarshalCallbackData(models.PollRoute, pollMsg)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall poll callback data: %s", err.Error())
	}

	backBtn := tgbotapi.InlineKeyboardButton{Text: "Back", CallbackData: backData}

	if voteNotAllowed {
		return render.Keyboard(backBtn), nil
	}

	voteMsg := &vote.Vote{
		PollId:   helpers.Ptr(pollID),
		OptionId: helpers.Ptr(int32(optionID)),
	}
	if referrerForecastsPage > 0 {
		voteMsg.ReferrerForecastsPage = helpers.Ptr(referrerForecastsPage)
	}

	data, err := proto2.MarshalCallbackData(models.VoteRoute, voteMsg)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall vote callback data: %s", err.Error())
	}

	return render.Keyboard(
		tgbotapi.InlineKeyboardButton{Text: "Yes", CallbackData: data},
		backBtn,
	), nil
}
