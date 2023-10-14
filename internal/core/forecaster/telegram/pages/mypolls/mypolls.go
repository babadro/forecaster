package mypolls

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
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

func (s *Service) NewRequest() (proto2.Message, *mypolls.MyPolls) {
	v := new(mypolls.MyPolls)

	return v, v
}

func (s *Service) RenderCallback(
	_ context.Context, _ *mypolls.MyPolls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	txtMsg := "Getting polls are not implemented yet"

	keyboard, err := keyboardMarkup()
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	return render.NewEditMessageTextWithKeyboard(
		upd.CallbackQuery.Message.Chat.ID, upd.CallbackQuery.Message.MessageID, txtMsg, keyboard), "", nil
}

func keyboardMarkup() (tgbotapi.InlineKeyboardMarkup, error) {
	editPollData, err := proto.MarshalCallbackData(models.EditPollRoute, &editpoll.EditPoll{})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal editPoll callback data: %s", err.Error())
	}

	mainMenuButton, err := proto.MarshalCallbackData(models.MainPageRoute, &mainpage.MainPage{})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal mainPage callback data: %s", err.Error())
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:         "Create poll",
				CallbackData: editPollData,
			},
			tgbotapi.InlineKeyboardButton{
				Text:         "Main Menu",
				CallbackData: mainMenuButton,
			},
		),
	), nil
}
