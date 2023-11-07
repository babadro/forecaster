package editstatus

import (
	"context"
	"fmt"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editstatus"
	"github.com/babadro/forecaster/internal/helpers"
	models2 "github.com/babadro/forecaster/internal/models"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
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

func (s *Service) NewRequest() (proto2.Message, *editstatus.EditStatus) {
	v := new(editstatus.EditStatus)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editstatus.EditStatus, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	newStatus := req.GetStatus()
	if newStatus != editstatus.Status_ACTIVE && newStatus != editstatus.Status_FINISHED {
		return nil, "", fmt.Errorf("only active and finished statuses are allowed")
	}

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

	status, err := proto.PollStatusFromProto(newStatus)
	if err != nil {
		return nil, "", fmt.Errorf("unable to convert proto status to status: %s", err.Error())
	}

	if req.GetNeedConfirmation() {
		return s.confirmation(p, req.ReferrerMyPollsPage, chatID, messageID, newStatus)
	}

	if _, err = s.db.UpdatePoll(ctx, pollID, swagger.UpdatePoll{
		Status: swagger.PollStatus(status),
	}, time.Now()); err != nil {
		return nil, "", fmt.Errorf("unable to update poll: %s", err.Error())
	}

	return successUpdateStatus(p, req.ReferrerMyPollsPage, chatID, messageID, status)
}

func successUpdateStatus(
	p swagger.Poll, referrerMyPollsPage *int32, chatID int64, messageID int, newStatus models2.PollStatus,
) (tgbotapi.Chattable, string, error) {
	backData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              helpers.Ptr(p.ID),
		ReferrerMyPollsPage: referrerMyPollsPage,
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

	switch newStatus {
	case models2.ActivePollStatus:
		sb.Printf("Poll <b>%s</b> is active now", p.Title)
	case models2.FinishedPollStatus:
		sb.Printf("Poll <b>%s</b> is finished now", p.Title)
	case models2.UnknownPollStatus, models2.DraftPollStatus:
		return nil, "", fmt.Errorf("not allowed status %s", newStatus.String())
	default:
		return nil, "", fmt.Errorf("unknown status %d", newStatus)
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, sb.String(), keyboard), "", nil
}

func (s *Service) confirmation(
	p swagger.Poll, referrerMyPollsPage *int32, chatID int64, messageID int, newStatus editstatus.Status,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := confirmationKeyboard(p.ID, referrerMyPollsPage, newStatus)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create confirmation keyboard: %s", err.Error())
	}

	txt, err := confirmationTxt(p, models2.ActivePollStatus)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create confirmation text: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

func confirmationTxt(p swagger.Poll, newStatus models2.PollStatus) (string, error) {
	var sb render.StringBuilder

	switch newStatus {
	case models2.ActivePollStatus:
		sb.WriteString("<b>Are you sure you want to activate this poll?</b>\n\n")
	case models2.FinishedPollStatus:
		sb.WriteString("<b>Are you sure you want to finish this poll?</b>\n\n")
	case models2.UnknownPollStatus, models2.DraftPollStatus:
		return "", fmt.Errorf("not allowed status %s", newStatus.String())
	default:
		return "", fmt.Errorf("unknown status %d", newStatus)
	}

	sb.Printf("<b>%s</b>\n", p.Title)

	sb.Printf("<i>Start Date: %s</i>\n", render.FormatTime(p.Start))
	sb.Printf("<i>End Date: %s</i>\n", render.FormatTime(p.Finish))

	return sb.String(), nil
}

func confirmationKeyboard(pollID int32, referrerMyPollsPage *int32, newStatus editstatus.Status) (tgbotapi.InlineKeyboardMarkup, error) {
	pollIDPtr := helpers.Ptr(pollID)

	editStatusData, err := proto.MarshalCallbackData(models.EditStatusRoute, &editstatus.EditStatus{
		PollId:              pollIDPtr,
		Status:              helpers.Ptr(newStatus),
		ReferrerMyPollsPage: referrerMyPollsPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal delete callback data: %s", err.Error())
	}

	backData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              pollIDPtr,
		ReferrerMyPollsPage: referrerMyPollsPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal go back callback data: %s", err.Error())
	}

	var statusButtonText string

	switch newStatus {
	case editstatus.Status_ACTIVE:
		statusButtonText = "Activate poll"
	case editstatus.Status_FINISHED:
		statusButtonText = "Finish poll"
	case editstatus.Status_UNKNOWN, editstatus.Status_DRAFT:
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("not allowed status %s", newStatus.String())
	default:
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unknown status %d", newStatus)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: statusButtonText, CallbackData: editStatusData},
			tgbotapi.InlineKeyboardButton{Text: "Go back", CallbackData: backData},
		),
	), nil
}
