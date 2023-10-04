package editpoll

import (
	"context"
	"errors"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
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

	pollID := req.GetPollId()

	if pollID == 0 {
		return s.createPoll(ctx, upd, user, chat, message)
	}

	return nil, "", fmt.Errorf("edit poll is not implemented")
}

func (s *Service) createPoll(ctx context.Context, upd tgbotapi.Update, user *tgbotapi.User, chat *tgbotapi.Chat, message *tgbotapi.Message) (tgbotapi.Chattable, string, error) {

	txt := "some text here about creating new poll"

	return s.render(ctx, int32(pollID), 0, user.ID, chat.ID, message.MessageID, false)
}

func createPollKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	res := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0, 4),
	}

	for _, field := range []struct {
		text        string
		fieldToEdit editpoll.FieldToEdit
	}{
		{"Title", editpoll.FieldToEdit_TITLE},
		{"Description", editpoll.FieldToEdit_DESCRIPTION},
		{"Start date", editpoll.FieldToEdit_START_DATE},
		{"Finish date", editpoll.FieldToEdit_FINISH_DATE},
	} {

	}

	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{Text: "Create new poll", CallbackData: nil},
			},
		},
	}
}
