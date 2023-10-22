package deletepoll

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/deletepoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
	w  dbwrapper.Wrapper
}

func New(db models.DB) *Service {
	return &Service{
		db: db, w: dbwrapper.New(db),
	}
}

func (s *Service) NewRequest() (proto.Message, *deletepoll.DeletePoll) {
	v := new(deletepoll.DeletePoll)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *deletepoll.DeletePoll, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	pollID := req.GetPollId()
	if pollID == 0 {
		return nil, "", fmt.Errorf("poll id is undefined %v", req.PollId)
	}

	userID := upd.CallbackQuery.From.ID
	chatID := upd.CallbackQuery.Message.Chat.ID
	messageID := upd.CallbackQuery.Message.MessageID

	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	if req.GetNeedConfirmation() {
		return s.confirmation(p, req.ReferrerMyPollsPage, chatID, messageID)
	}

	if err = s.db.DeletePoll(ctx, pollID); err != nil {
		return nil, "", fmt.Errorf("unable to delete poll: %s", err.Error())
	}

	return successDeletion(p.Title, req.GetReferrerMyPollsPage(), chatID, messageID)
}

func successDeletion(
	pollTitle string, referrerMyPollsPage int32, chatID int64, messageID int,
) (tgbotapi.Chattable, string, error) {
	backData, err := proto2.MarshalCallbackData(models.MyPollsRoute, &mypolls.MyPolls{
		CurrentPage: helpers.OneIfZero(referrerMyPollsPage),
	})
	if err != nil {
		return nil, "", fmt.Errorf("unable to marshal go back callback data: %s", err.Error())
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Go back", CallbackData: backData},
		),
	)

	var sb render.StringBuilder

	sb.WriteString("<b>Poll was successfully deleted!</b>\n\n")

	sb.Printf("<b>%s</b>\n", pollTitle)

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, sb.String(), keyboard), "", nil
}

func (s *Service) confirmation(
	p swagger.PollWithOptions, referrerMyPollsPage *int32, chatID int64, messageID int,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := confirmationKeyboard(p.ID, referrerMyPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create confirmation keyboard: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, confirmationTxt(p), keyboard), "", nil
}

func confirmationTxt(p swagger.PollWithOptions) string {
	var sb render.StringBuilder

	sb.WriteString("<b>Are you sure you want to delete this poll?</b>\n\n")

	sb.Printf("<b>%s</b>\n", p.Title)

	sb.WriteString("<b>Options:</b>\n")

	for _, o := range p.Options {
		sb.WriteString(fmt.Sprintf("%s\n", o.Title))
	}

	sb.Printf("<i>Start Date: %s</i>\n", render.FormatTime(p.Start))
	sb.Printf("<i>End Date: %s</i>\n", render.FormatTime(p.Finish))

	return sb.String()
}

func confirmationKeyboard(pollID int32, referrerMyPollsPage *int32) (tgbotapi.InlineKeyboardMarkup, error) {
	pollIDPtr := helpers.Ptr(pollID)

	deleteData, err := proto2.MarshalCallbackData(models.DeletePollRoute, &deletepoll.DeletePoll{
		PollId:              pollIDPtr,
		ReferrerMyPollsPage: referrerMyPollsPage,
		NeedConfirmation:    helpers.Ptr(false),
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal delete callback data: %s", err.Error())
	}

	backData, err := proto2.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              pollIDPtr,
		ReferrerMyPollsPage: referrerMyPollsPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal go back callback data: %s", err.Error())
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "Delete", CallbackData: deleteData},
			tgbotapi.InlineKeyboardButton{Text: "Go back", CallbackData: backData},
		),
	), nil
}
