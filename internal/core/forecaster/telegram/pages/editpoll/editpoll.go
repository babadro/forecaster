package editpoll

import (
	"context"
	"errors"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
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

func (s *Service) NewRequest() (proto2.Message, *editpoll.EditPoll) {
	v := new(editpoll.EditPoll)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editpoll.EditPoll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	user := upd.CallbackQuery.From

	if user == nil {
		return nil, "", errors.New("user is nil")
	}

	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetCreatePoll() {
		return s.createPoll(req.GetReferrerMyPollsPage(), message.MessageID, chat.ID)
	}

	return nil, "", fmt.Errorf("edit poll is not implemented")
}

func (s *Service) createPoll(myPollsPage int32, messageID int, chatID int64) (tgbotapi.Chattable, string, error) {
	// todo text
	txt := "some text here about creating new poll"

	keyboard, err := createPollKeyboardMarkup(myPollsPage)
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

func createPollKeyboardMarkup(myPollsPage int32) (tgbotapi.InlineKeyboardMarkup, error) {
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
			PollId: helpers.Ptr[int32](0),
			Field:  helpers.Ptr(editButtons[i].Field),
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
