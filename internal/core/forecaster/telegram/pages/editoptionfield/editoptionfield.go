package editoptionfield

import (
	"context"
	"errors"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoptionfield"
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

func (s *Service) NewRequest() (proto2.Message, *editoptionfield.EditOptionField) {
	v := new(editoptionfield.EditOptionField)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editoptionfield.EditOptionField, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	field := req.GetField()
	if field == editoptionfield.Field_UNDEFINED {
		return nil, "", errors.New("field is undefined")
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

	keyboard, err := keyboardMarkup(req.PollId, req.OptionId, req.ReferrerMyPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editpollfield page: %s", err.Error())
	}

	op := swagger.Option{PollID: pollID}

	if optionID := req.GetOptionId(); optionID != 0 {
		opPtr, idx := swagger.FindOptionByID(p.Options, int16(optionID))
		if idx == -1 {
			return nil, "", fmt.Errorf("option %d not found", optionID)
		}

		op = *opPtr
	} else if field != editoptionfield.Field_TITLE {
		errMessage := "First create Title, please, and then you can create other fields"

		return render.NewEditMessageTextWithKeyboard(chatID, messageID, errMessage, keyboard), "", nil
	}

	txt, err := txtMsg(op, field, req.GetReferrerMyPollsPage())
	if err != nil {
		return nil, "", fmt.Errorf("unable to create text for editpollfield page: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

func txtMsg(op swagger.Option, field editoptionfield.Field, referrerMyPollsPage int32) (string, error) {
	var sb render.StringBuilder

	sb.Printf("%s %d %d %s %d\n", models.EditOptionCommand, op.PollID, op.ID, field.String(), referrerMyPollsPage)

	sb.WriteString("\nEnter new value in reply to this message")

	sb.WriteString("\nCurrent value:\n")

	var fieldValue string

	switch field {
	case editoptionfield.Field_TITLE:
		fieldValue = op.Title
	case editoptionfield.Field_DESCRIPTION:
		fieldValue = op.Description
	case editoptionfield.Field_UNDEFINED:
		return "", errors.New("field is undefined")
	default:
		return "", fmt.Errorf("unknown field %s", field.String())
	}

	sb.WriteString(fieldValue)

	return sb.String(), nil
}

func keyboardMarkup(pollID, optionID, myPollsPage *int32) (tgbotapi.InlineKeyboardMarkup, error) {
	goBackData, err := proto.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
		PollId:              pollID,
		OptionId:            optionID,
		ReferrerMyPollsPage: myPollsPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for go back button: %s", err.Error())
	}

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{{{
			Text:         "Go back",
			CallbackData: goBackData,
		}}},
	}, nil
}
