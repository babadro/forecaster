package editfield

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
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

func (s *Service) NewRequest() (proto2.Message, *editfield.EditField) {
	v := new(editfield.EditField)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editfield.EditField, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	field := req.GetField()
	if field == editfield.Field_UNDEFINED {
		return nil, "", fmt.Errorf("field is undefined")
	}

	var (
		p      swagger.PollWithOptions
		errMsg string
		err    error
	)

	userID := upd.CallbackQuery.From.ID
	chatID := upd.CallbackQuery.Message.Chat.ID
	messageID := upd.CallbackQuery.Message.MessageID

	keyboard, err := keyboardMarkup(req.PollId, req.ReferrerMyPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editField page: %s", err.Error())
	}

	if pollID := req.GetPollId(); pollID != 0 {
		p, errMsg, err = s.w.GetPollByID(ctx, pollID)
		if err != nil {
			return nil, errMsg, err
		}

		if p.TelegramUserID != userID {
			return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
		}
	} else if field != editfield.Field_TITLE {
		errMessage := "First create Title, please, and then you can create other fields"

		return render.NewEditMessageTextWithKeyboard(chatID, messageID, errMessage, keyboard), "", nil
	}

	txt, err := txtMsg(p, field, req.GetReferrerMyPollsPage())
	if err != nil {
		return nil, "", fmt.Errorf("unable to create text for editField page: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

func txtMsg(p swagger.PollWithOptions, field editfield.Field, referrerMyPollsPage int32) (string, error) {
	var sb render.StringBuilder

	sb.Printf("/editpoll %d %s %d\n", p.ID, field.String(), referrerMyPollsPage)

	sb.WriteString("\nEnter new value in reply to this message")

	sb.WriteString("\nCurrent value:\n")

	var fieldValue string

	switch field {
	case editfield.Field_TITLE:
		fieldValue = p.Title
	case editfield.Field_DESCRIPTION:
		fieldValue = p.Description
	case editfield.Field_START_DATE:
		fieldValue = p.Start.String()
	case editfield.Field_FINISH_DATE:
		fieldValue = p.Finish.String()
	case editfield.Field_UNDEFINED:
		return "", fmt.Errorf("field is undefined")
	default:
		return "", fmt.Errorf("unknown field %d", field)
	}

	sb.WriteString(fieldValue)

	return sb.String(), nil
}

func keyboardMarkup(pollID, myPollsPage *int32) (tgbotapi.InlineKeyboardMarkup, error) {
	goBackData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              pollID,
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
