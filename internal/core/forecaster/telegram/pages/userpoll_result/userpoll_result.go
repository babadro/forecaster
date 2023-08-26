package userpollresult

import (
	"context"
	"fmt"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
	w  dbwrapper.Wrapper
}

func New(db models.DB) *Service {
	return &Service{db: db, w: dbwrapper.New(db)}
}

func (s *Service) NewRequest() (proto2.Message, *userpollresult.UserPollResult) {
	v := new(userpollresult.UserPollResult)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *userpollresult.UserPollResult, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, req.GetPollId())
	if err != nil {
		return nil, errMsg, err
	}

	outcome, ok := swagger.GetOutcome(p.Options)
	if !ok {
		return nil, "", fmt.Errorf("userpoll result: can't get outcome for pollID: %d", p.ID)
	}

	user := upd.CallbackQuery.From

	lastVote, found, err := s.w.GetLastVote(ctx, user.ID, p.ID)
	if err != nil {
		return nil, "", err
	}

	if !found {
		return nil, "", fmt.Errorf("userpoll result: can't find last user's vote for pollID: %d", p.ID)
	}

	if lastVote.OptionID != outcome.ID {
		return nil, "", fmt.Errorf("userpoll result: last user's vote is not outcome for pollID: %d", p.ID)
	}

	msg := txtMsg(user.UserName, outcome.Title, time.Time(p.Finish), lastVote.EpochUnixTimestamp)

	markup, err := keyboardMarkup()
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to create keyboard markup: %s", err.Error())
	}

	origMsg := upd.CallbackQuery.Message

	return render.NewEditMessageTextWithKeyboard(origMsg.Chat.ID, origMsg.MessageID, msg, markup), "", nil
}

func txtMsg(userName, optionTitle string, finishPoll time.Time, voteUnixTime int64) string {
	var sb render.StringBuilder

	advanceTimeNumber, advanceTimeUnit := render.GetHighestTimeUnit(finishPoll.Sub(time.Unix(voteUnixTime, 0)))

	sb.Printf("<b>%s</b> you predicted that %s %d %s before!", userName, optionTitle, advanceTimeNumber, advanceTimeUnit)

	return sb.String()
}

func keyboardMarkup() (tgbotapi.InlineKeyboardMarkup, error) {
	return render.Keyboard(
		tgbotapi.InlineKeyboardButton{
			// todo
		},
	), nil
}
