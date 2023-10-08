package editpoll

import (
	"context"
	"fmt"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"time"
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

func (s *Service) NewRequest() (proto2.Message, *editpoll.EditPoll) {
	v := new(editpoll.EditPoll)

	return v, v
}

func (s *Service) RenderCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	if update.Message.ReplyToMessage == nil {
		return nil, "", fmt.Errorf("parsing command directly from message is not implemented yet")
	}

	pollID, fieldID, err := parseCommandArgs(update.Message.ReplyToMessage.Text)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse command args: %s", err.Error())
	}

	createModel, updateModel, err := getDBModel(pollID, fieldID, update.Message.Text)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get db model: %s", err.Error())
	}

	var (
		p   swagger.Poll
		err error
	)

	if pollID == 0 {
		p, err = s.db.CreatePoll(ctx, createModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to create poll: %s", err.Error())
		}
	} else {
		p, err = s.db.UpdatePoll(ctx, pollID, updateModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to update poll: %s", err.Error())
		}
	}

	return nil, "", fmt.Errorf("not implemented")
}

func getDBModel(pollID int32, fieldID editfield.Field, input string) (swagger.CreatePoll, swagger.UpdatePoll, error) {
	create, update := swagger.CreatePoll{}, swagger.UpdatePoll{}

	switch fieldID {
	case editfield.Field_TITLE:
		create.Title, update.Title = input, &input
	case editfield.Field_DESCRIPTION:
		create.Description, update.Description = input, &input
	case editfield.Field_START_DATE:
		date, err := parseDate(input)
		if err != nil {
			return swagger.CreatePoll{}, swagger.UpdatePoll{},
				fmt.Errorf("unable to parse start date %s date: %s", input, err.Error())
		}

		create.Start, update.Start = date, &date
	case editfield.Field_FINISH_DATE:
		date, err := parseDate(input)
		if err != nil {
			return swagger.CreatePoll{}, swagger.UpdatePoll{},
				fmt.Errorf("unable to parse finish date %s date: %s", input, err.Error())
		}

		create.Finish, update.Finish = date, &date
	default:
		return swagger.CreatePoll{}, swagger.UpdatePoll{},
			fmt.Errorf("unknown field %d", fieldID)
	}

	return create, update, nil
}

func parseDate(input string) (strfmt.DateTime, error) {
	return strfmt.DateTime(time.Now()), nil
}

func parseCommandArgs(text string) (int32, editfield.Field, error) {
	newLineIdx := strings.Index(text, "\n")
	if newLineIdx == -1 {
		return 0, 0, fmt.Errorf("no new line found")
	}

	strArr := strings.Split(text[:newLineIdx], " ")
	if len(strArr) != 3 {
		return 0, 0, fmt.Errorf("expected 3 words, got %d", len(strArr))
	}

	command, pollIDStr, field := strArr[0], strArr[1], strArr[2]
	if command != models.EditPollCommand {
		return 0, 0, fmt.Errorf("expected %s command, got %s", models.EditPollCommand, command)
	}

	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to parse pollID: %s", err.Error())
	}

	var fieldID editfield.Field

	switch field {
	case "title":
		fieldID = editfield.Field_TITLE
	case "description":
		fieldID = editfield.Field_DESCRIPTION
	case "start":
		fieldID = editfield.Field_START_DATE
	case "finish":
		fieldID = editfield.Field_FINISH_DATE
	default:
		return 0, 0, fmt.Errorf("unknown field %s", field)
	}

	return int32(pollID), fieldID, nil
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editpoll.EditPoll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetPollId() == 0 {
		return s.createPoll(req.ReferrerMyPollsPage, message.MessageID, chat.ID)
	}

	return s.editPoll(ctx, req.PollId, req.ReferrerMyPollsPage, message.MessageID, chat.ID, upd.CallbackQuery.From.ID)
}

func (s *Service) editPoll(ctx context.Context, pollID, myPollsPage *int32, messageID int, chatID, userID int64) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, *pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	keyboard, err := pollKeyboardMarkup(pollID, myPollsPage)
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

func (s *Service) createPoll(myPollsPage *int32, messageID int, chatID int64) (tgbotapi.Chattable, string, error) {
	// todo text
	txt := "some text here about creating new poll"

	keyboard, err := pollKeyboardMarkup(nil, myPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for createPoll page: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
}

type editButton struct {
	text  string
	Field editfield.Field
}

const maxCountInRow = 3

func pollKeyboardMarkup(pollID, myPollsPage *int32) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []editButton{
		{"Title", editfield.Field_TITLE},
		{"Description", editfield.Field_DESCRIPTION},
		{"Start date", editfield.Field_START_DATE},
		{"Finish date", editfield.Field_FINISH_DATE},
	}

	buttonsCount := len(editButtons) + 1 // +1 for Go back button

	keyboardBuilder := render.NewKeyboardBuilder(maxCountInRow, buttonsCount)

	for i := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditFieldRoute, &editpoll.EditPoll{
			PollId:              pollID,
			ReferrerMyPollsPage: myPollsPage,
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
		CurrentPage: myPollsPage,
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
