package editoption

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/deleteoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoptionfield"
	"github.com/babadro/forecaster/internal/helpers"
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

func (s *Service) NewRequest() (proto2.Message, *editoption.EditOption) {
	v := new(editoption.EditOption)

	return v, v
}

func (s *Service) RenderCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	args, err := parseCommandArgs(update.Message.ReplyToMessage.Text)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse command args: %s", err.Error())
	}

	if err = validateCommandArgs(args); err != nil {
		return nil, "", fmt.Errorf("invalid command args: %s", err.Error())
	}

	if args.optionID == 0 {

	}

	return s.editOptionDialog(ctx, "", &args.pollID, &args.optionID, &args.myPollsPage,
		update.Message.MessageID, update.Message.Chat.ID, update.Message.From.ID, false)
}

func validateCommandArgs(args commandArgs) error {
	if args.pollID == 0 {
		return fmt.Errorf("poll id should not be zero")
	}

	return nil
}

func getUpdateModel(fieldID editoptionfield.Field, input string) (swagger.UpdateOption, error) {
	res := swagger.UpdateOption{}

	switch fieldID {
	case editoptionfield.Field_TITLE:
		res.Title = &input
	case editoptionfield.Field_DESCRIPTION:
		res.Description = &input
	default:
		return swagger.UpdateOption{}, fmt.Errorf("unknown field %s", fieldID)
	}

	return res, nil
}

func getCreateModel(fieldID editoptionfield.Field, input string) (swagger.CreateOption, error) {
	res := swagger.CreateOption{}

	switch fieldID {
	case editoptionfield.Field_TITLE:
		res.Title = input
	case editoptionfield.Field_DESCRIPTION:
		res.Description = input
	default:
		return swagger.CreateOption{}, fmt.Errorf("unknown field %s", fieldID)
	}

	return res, nil
}

type commandArgs struct {
	pollID      int32
	optionID    int16
	field       editoptionfield.Field
	myPollsPage int32
}

const minCommandParts = 4

func parseCommandArgs(text string) (commandArgs, error) {
	newLineIDx := strings.Index(text, "\n")
	if newLineIDx == -1 {
		return commandArgs{}, fmt.Errorf("no new line found")
	}

	strArr := strings.Split(text[:newLineIDx], " ")
	if len(strArr) < minCommandParts {
		return commandArgs{}, fmt.Errorf("expected at least %d command parts, got %d", minCommandParts, len(strArr))
	}

	command, pollIDStr, optionIDStr, field := strArr[0], strArr[1], strArr[2], strArr[3]
	if command != models.EditOptionCommand {
		return commandArgs{}, fmt.Errorf("expected command %s, got %s", models.EditOptionCommand, command)
	}

	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		return commandArgs{}, fmt.Errorf("unable to parse poll id: %s", err.Error())
	}

	optionID, err := strconv.ParseInt(optionIDStr, 10, 16)
	if err != nil {
		return commandArgs{}, fmt.Errorf("unable to parse option id: %s", err.Error())
	}

	fieldID, ok := editoptionfield.Field_value[field]
	if !ok {
		return commandArgs{}, fmt.Errorf("unknown field %s", field)
	}

	var myPollsPage int64
	if len(strArr) > minCommandParts {
		myPollsPage, err = strconv.ParseInt(strArr[4], 10, 32)
		if err != nil {
			return commandArgs{}, fmt.Errorf("unable to parse my polls page: %s", err.Error())
		}
	}

	return commandArgs{
		pollID:      int32(pollID),
		optionID:    int16(optionID),
		field:       editoptionfield.Field(fieldID),
		myPollsPage: int32(myPollsPage),
	}, nil
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editoption.EditOption, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetOptionId() == 0 {
		return s.createOptionDialog("", req.PollId, req.ReferrerMyPollsPage,
			message.MessageID, chat.ID, true)
	}

	if req.GetPollId() == 0 {
		return nil, "", fmt.Errorf("can't edit poll: poll id %v is undefined", req.PollId)
	}

	return s.editOptionDialog(ctx, "", req.PollId, req.OptionId, req.ReferrerMyPollsPage,
		message.MessageID, chat.ID, upd.CallbackQuery.From.ID, true)
}

func (s *Service) editOptionDialog(
	ctx context.Context, validationErr string, pollID int32, optionID int16, myPollsPage int32,
	updateModel swagger.UpdateOption, doUpdate bool,
	messageID int, chatID, userID int64, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	op, idx := swagger.FindOptionByID(p.Options, optionID)
	if idx == -1 {
		return nil, "", fmt.Errorf("option %d not found in poll %d", optionID, pollID)
	}

	if doUpdate {
		updatedOption, err := s.db.UpdateOption(ctx, pollID, optionID, updateModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to update option: %s", err.Error())
		}

		p.Options[idx] = &updatedOption
	}

	keyboard, err := keyboardMarkup(pollID, int32(optionID), myPollsPage)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for options: %s", err.Error())
	}

	txt := editOptionTxt(validationErr, op)

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func editOptionTxt(validationErrMsg string, op *swagger.Option) string {
	var sb render.StringBuilder

	if validationErrMsg != "" {
		sb.Printf("<b>🚨🚨🚨\n%s\n🚨🚨🚨</b>\n\n", validationErrMsg)
	}

	sb.Printf("Title:\n<b>%s</b>\n", op.Title)
	sb.Printf("\nDescription\n:<b>%s</b>\n", op.Description)

	return sb.String()
}

func (s *Service) createOptionDialog(
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

func keyboardMarkup(pollID, optionID, myPollsPage int32) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []models.EditButton[editoptionfield.Field]{
		{"Title", editoptionfield.Field_TITLE},
		{"Description", editoptionfield.Field_DESCRIPTION},
	}

	buttonsCount := len(editButtons) + 2 // +2 for delete and back buttons

	keyboardBuilder := render.NewKeyboardBuilder(maxCountInRow, buttonsCount)

	pollIDPtr := helpers.NilIfZero(pollID)
	optionIDPtr := helpers.NilIfZero(optionID)
	myPollsPagePtr := helpers.NilIfZero(myPollsPage)

	for _, editButton := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditOptionFieldRoute, &editoptionfield.EditOptionField{
			PollId:              pollIDPtr,
			OptionId:            optionIDPtr,
			Field:               helpers.Ptr(editButton.Field),
			ReferrerMyPollsPage: myPollsPagePtr,
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

	if optionID != 0 {
		callbackData, err := proto.MarshalCallbackData(models.DeleteOptionRoute, &deleteoption.DeleteOption{
			PollId:              pollIDPtr,
			OptionId:            optionIDPtr,
			ReferrerMyPollsPage: myPollsPagePtr,
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
