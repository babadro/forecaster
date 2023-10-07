package editpoll

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db            models.DB
	w             dbwrapper.Wrapper
	fieldIDtoName map[editpoll.Field]string
}

func New(db models.DB) *Service {
	return &Service{
		db: db, w: dbwrapper.New(db),
		fieldIDtoName: map[editpoll.Field]string{
			editpoll.Field_TITLE:       "title",
			editpoll.Field_DESCRIPTION: "description",
			editpoll.Field_START_DATE:  "start date",
			editpoll.Field_FINISH_DATE: "finish date",
		},
	}
}

func (s *Service) NewRequest() (proto2.Message, *editpoll.EditPoll) {
	v := new(editpoll.EditPoll)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editpoll.EditPoll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetCreatePoll() {
		return s.createPoll(req.GetReferrerMyPollsPage(), message.MessageID, chat.ID)
	}

	return nil, "", fmt.Errorf("edit poll is not implemented")
}

func (s *Service) editField(ctx context.Context, pollID int32, myPollsPage int32, field editpoll.Field, chatID, userID int64) (tgbotapi.Chattable, string, error) {
	var (
		p      swagger.PollWithOptions
		errMsg string
		err    error
	)

	if pollID != 0 {
		p, errMsg, err = s.w.GetPollByID(ctx, pollID)
		if err != nil {
			return nil, errMsg, err
		}

		if p.TelegramUserID != userID {
			return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
		}
	}

	txt, err := s.editFieldTxt(p, field)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create text for editField page: %s", err.Error())
	}

	keyboard, err := editFieldKeyboardMarkup(pollID, myPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editField page: %s", err.Error())
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func (s *Service) editFieldTxt(p swagger.PollWithOptions, field editpoll.Field) (string, error) {
	var sb render.StringBuilder
	sb.Printf("/editpoll %d %s\n", p.ID, field.String())

	sb.WriteString("\nEnter new value in reply to this message")

	sb.WriteString("\nCurrent value:\n")

	fieldValue := ""
	switch field {
	case editpoll.Field_TITLE:
		fieldValue = p.Title
	case editpoll.Field_DESCRIPTION:
		fieldValue = p.Description
	case editpoll.Field_START_DATE:
		fieldValue = p.Start.String()
	case editpoll.Field_FINISH_DATE:
		fieldValue = p.Finish.String()
	default:
		return "", fmt.Errorf("unknown field %d", field)
	}

	sb.WriteString(fieldValue)

	return sb.String(), nil
}

func editFieldKeyboardMarkup(pollID, myPollsPage int32) (tgbotapi.InlineKeyboardMarkup, error) {
	goBackData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
		PollId:              pollID,
		CreatePoll:          nil,
		Field:               nil,
		ReferrerMyPollsPage: nil,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for go back button: %s", err.Error())
	}

	keyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Go back",
		CallbackData: goBackData,
	})
	return tgbotapi.InlineKeyboardMarkup{}, nil
}

func (s *Service) editPoll(ctx context.Context, pollID, myPollsPage int32, messageID int, chatID, userID int64) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	keyboard, err := pollKeyboardMarkup(pollID, myPollsPage, nil)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editPoll page: %s", err.Error())
	}

	txt := editPollTxt(p)

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

func editPollTxt(p swagger.PollWithOptions) string {
	// todo text
	return "some text here about editing poll"
}

func (s *Service) createPoll(myPollsPage int32, messageID int, chatID int64) (tgbotapi.Chattable, string, error) {
	// todo text
	txt := "some text here about creating new poll"

	keyboard, err := pollKeyboardMarkup(0, myPollsPage, helpers.Ptr(true))
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for createPoll page: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

type editButton struct {
	text  string
	Field editpoll.Field
}

const maxCountInRow = 3

func pollKeyboardMarkup(pollID, myPollsPage int32, createPoll *bool) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []editButton{
		{"Title", editpoll.Field_TITLE},
		{"Description", editpoll.Field_DESCRIPTION},
		{"Start date", editpoll.Field_START_DATE},
		{"Finish date", editpoll.Field_FINISH_DATE},
	}

	buttonsCount := len(editButtons) + 1 // +1 for Go back button

	keyboardBuilder := render.NewKeyboardBuilder(maxCountInRow, buttonsCount)

	for i := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{
			PollId:              helpers.Ptr(pollID),
			Field:               helpers.Ptr(editButtons[i].Field),
			CreatePoll:          createPoll,
			ReferrerMyPollsPage: helpers.Ptr(myPollsPage),
		})

		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal callback data for %s button: %s", editButtons[i].text, err.Error())
		}

		keyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         editButtons[i].text,
			CallbackData: callbackData,
		})
	}

	goBackData, err := proto.MarshalCallbackData(models.MyPollsRoute, &mypolls.MyPolls{
		CurrentPage: helpers.Ptr(myPollsPage),
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for go back button: %s", err.Error())
	}

	keyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Go back",
		CallbackData: goBackData,
	})

	return keyboardBuilder.MarkUp(), nil
}
