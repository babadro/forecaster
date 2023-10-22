package deleteoption

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/deleteoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
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

func (s *Service) NewRequest() (proto.Message, *deleteoption.DeleteOption) {
	v := new(deleteoption.DeleteOption)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *deleteoption.DeleteOption, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	pollID := req.GetPollId()
	if pollID == 0 {
		return nil, "", fmt.Errorf("poll id is undefined %v", req.PollId)
	}

	optionID := req.GetOptionId()
	if optionID == 0 {
		return nil, "", fmt.Errorf("option id is undefined %v", req.OptionId)
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

	op, idx := swagger.FindOptionByID(p.Options, int16(optionID))
	if idx == -1 {
		return nil, "", fmt.Errorf("option %d not found in poll %d", optionID, pollID)
	}

	if req.GetNeedConfirmation() {
		return s.confirmation(*op, req.ReferrerMyPollsPage, chatID, messageID)
	}

	if err = s.db.DeleteOption(ctx, pollID, int16(optionID)); err != nil {
		return nil, "", fmt.Errorf("unable to delete option: %s", err.Error())
	}

	return successDeletion(op.Title, p.ID, req.GetReferrerMyPollsPage(), chatID, messageID)
}

func successDeletion(
	optionTitle string, pollID, referrerMyPollsPage int32, chatID int64, messageID int,
) (tgbotapi.Chattable, string, error) {
	backData, err := proto2.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              helpers.Ptr(pollID),
		ReferrerMyPollsPage: helpers.Ptr(referrerMyPollsPage),
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

	sb.WriteString("<b>Option was successfully deleted!</b>\n\n")

	sb.Printf("<b>%s</b>\n", optionTitle)

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, sb.String(), keyboard), "", nil
}

func (s *Service) confirmation(
	op swagger.Option, referrerMyPollsPage *int32, chatID int64, messageID int,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := confirmationKeyboard(op.PollID, op.ID, referrerMyPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create confirmation keyboard: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, confirmationTxt(op), keyboard), "", nil
}

func confirmationTxt(op swagger.Option) string {
	var sb render.StringBuilder

	sb.WriteString("<b>Are you sure you want to delete this option?</b>\n\n")

	sb.Printf("<b>%s</b>\n", op.Title)

	return sb.String()
}

func confirmationKeyboard(
	pollID int32, optionID int16, referrerMyPollsPage *int32,
) (tgbotapi.InlineKeyboardMarkup, error) {
	pollIDPtr, optionIDPtr := helpers.Ptr(pollID), helpers.Ptr(int32(optionID))

	deleteData, err := proto2.MarshalCallbackData(models.DeleteOptionRoute, &deleteoption.DeleteOption{
		PollId:           pollIDPtr,
		NeedConfirmation: helpers.Ptr(true),
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal delete callback data: %s", err.Error())
	}

	backData, err := proto2.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
		PollId:              pollIDPtr,
		OptionId:            optionIDPtr,
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
