package editpoll

import (
	"context"
	"fmt"
	"strings"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
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

func (s *Service) NewRequest() (proto2.Message, *editpoll.EditPoll) {
	v := new(editpoll.EditPoll)

	return v, v
}

func (s *Service) RenderCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	if update.Message.ReplyToMessage == nil {
		return nil, "", fmt.Errorf("parsing command directly from message is not implemented yet")
	}

	parentText := update.Message.ReplyToMessage.Text

	newLineIdx := strings.Index(parentText, "\n")
	if newLineIdx == -1 {
		return nil, "", fmt.Errorf("unable to parse parent message text")
	}

	strings.Split(parentText[:newLineIdx], " ")

	// todo

	return nil, "", fmt.Errorf("not implemented")
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
