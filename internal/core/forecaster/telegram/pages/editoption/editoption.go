package editoption

import (
	"context"
	"fmt"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/deleteoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoptionfield"
	"github.com/babadro/forecaster/internal/helpers"
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

func (s *Service) NewRequest() (proto2.Message, *editoption.EditOption) {
	v := new(editoption.EditOption)

	return v, v
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editoption.EditOption, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetOptionId() == 0 {
		return s.createOption("", req.PollId, req.ReferrerMyPollsPage,
			message.MessageID, chat.ID, true)
	}

	if req.GetPollId() == 0 {
		return nil, "", fmt.Errorf("can't edit poll: poll id is undefined")
	}

	return nil, "", nil
}

func (s *Service) editOption(
	ctx context.Context, pollID, optionID, myPollsPage *int32, messageID int, chatID, userID int64,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, *pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	return nil, "", nil
}

func (s *Service) createOption(
	validationErrMsg string, pollID, myPollsPage *int32, messageID int, chatID int64, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := keyboardMarkup(pollID, nil, myPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for options: %s", err.Error())
	}

	txt := createOptionTxt(validationErrMsg)

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func createOptionTxt(validationErrMsg string) string {
	txt := `Define your option title and description.`

	var sb render.StringBuilder
	if validationErrMsg != "" {
		sb.Printf("\n\n<b>%s</b>", validationErrMsg)
	}

	sb.WriteString(txt)

	return sb.String()
}

const (
	maxCountInRow = 2
)

func keyboardMarkup(pollID, optionID, myPollsPage *int32) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []models.EditButton[editoptionfield.Field]{
		{"Title", editoptionfield.Field_TITLE},
		{"Description", editoptionfield.Field_DESCRIPTION},
	}

	buttonsCount := len(editButtons) + 2 // +2 for delete and back buttons

	keyboardBuilder := render.NewKeyboardBuilder(maxCountInRow, buttonsCount)

	for _, editButton := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditOptionFieldRoute, &editoptionfield.EditOptionField{
			PollId:              pollID,
			OptionId:            optionID,
			Field:               helpers.Ptr(editButton.Field),
			ReferrerMyPollsPage: myPollsPage,
		})

		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal edit option field callback data: %s", err.Error())
		}

		keyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         editButton.Text,
			CallbackData: callbackData,
		})
	}

	if !helpers.IsZero(optionID) {
		callbackData, err := proto.MarshalCallbackData(models.DeleteOptionRoute, &deleteoption.DeleteOption{
			PollId:              pollID,
			OptionId:            optionID,
			ReferrerMyPollsPage: myPollsPage,
		})

		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal delete option callback data: %s", err.Error())
		}

		keyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         "Delete",
			CallbackData: callbackData,
		})
	}

	return tgbotapi.InlineKeyboardMarkup{}, nil
}
