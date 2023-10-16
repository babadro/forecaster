package editpoll

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
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editfield"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
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

	args, err := parseCommandArgs(update.Message.ReplyToMessage.Text)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse command args: %s", err.Error())
	}

	createModel, updateModel, validationErr, err := getDBModel(
		ctx, args.field, update.Message.Text, update.Message.From.ID,
	)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get db model: %s", err.Error())
	}

	if validationErr != "" {
		if args.pollID == 0 {
			return s.createPoll(validationErr, helpers.NilIfZero(args.myPollsPage),
				update.Message.MessageID, update.Message.Chat.ID, false)
		}

		return s.editPoll(ctx, validationErr, &args.pollID, helpers.NilIfZero(args.myPollsPage),
			update.Message.MessageID, update.Message.Chat.ID, update.Message.From.ID, false)
	}

	var p swagger.Poll

	if args.pollID == 0 {
		p, err = s.db.CreatePoll(ctx, createModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to create poll: %s", err.Error())
		}
	} else {
		p, err = s.db.UpdatePoll(ctx, args.pollID, updateModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to update poll: %s", err.Error())
		}
	}

	return s.editPoll(ctx, "", &p.ID, helpers.NilIfZero(args.myPollsPage),
		update.Message.MessageID, update.Message.Chat.ID, update.Message.From.ID, false)
}

func getDBModel(
	ctx context.Context, fieldID editfield.Field, input string, telegramUserID int64,
) (swagger.CreatePoll, swagger.UpdatePoll, string, error) {
	create := swagger.CreatePoll{TelegramUserID: telegramUserID}
	update := swagger.UpdatePoll{TelegramUserID: &telegramUserID}

	invalidDateErrMsg := fmt.Sprintf("Can't parse date format.\nIt should be %s", time.RFC3339)

	switch fieldID {
	case editfield.Field_TITLE:
		create.Title, update.Title = input, &input
	case editfield.Field_DESCRIPTION:
		create.Description, update.Description = input, &input
	case editfield.Field_START_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's start poll date")
			return swagger.CreatePoll{}, swagger.UpdatePoll{}, invalidDateErrMsg, nil
		}

		create.Start, update.Start = date, &date
	case editfield.Field_FINISH_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's finish poll date")
			return swagger.CreatePoll{}, swagger.UpdatePoll{}, invalidDateErrMsg, nil
		}

		create.Finish, update.Finish = date, &date
	case editfield.Field_UNDEFINED:
		return swagger.CreatePoll{}, swagger.UpdatePoll{}, "", fmt.Errorf("field is undefined")
	default:
		return swagger.CreatePoll{}, swagger.UpdatePoll{}, "",
			fmt.Errorf("unknown field %d", fieldID)
	}

	return create, update, "", nil
}

func parseDate(input string) (strfmt.DateTime, error) {
	t, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return strfmt.DateTime{}, fmt.Errorf("unable to parse date %s: %s", input, err.Error())
	}

	return strfmt.DateTime(t), nil
}

type commandArgs struct {
	pollID      int32
	field       editfield.Field
	myPollsPage int32
}

const minCommandParts = 3

func parseCommandArgs(text string) (commandArgs, error) {
	newLineIDx := strings.Index(text, "\n")
	if newLineIDx == -1 {
		return commandArgs{}, fmt.Errorf("no new line found")
	}

	strArr := strings.Split(text[:newLineIDx], " ")
	if len(strArr) < minCommandParts {
		return commandArgs{}, fmt.Errorf("expected at least %d command parts, got %d", minCommandParts, len(strArr))
	}

	command, pollIDStr, field := strArr[0], strArr[1], strArr[2]
	if command != models.EditPollCommand {
		return commandArgs{}, fmt.Errorf("expected %s command, got %s", models.EditPollCommand, command)
	}

	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		return commandArgs{}, fmt.Errorf("unable to parse pollID: %s", err.Error())
	}

	fieldID, ok := editfield.Field_value[field]
	if !ok {
		return commandArgs{}, fmt.Errorf("unknown field %s", field)
	}

	var myPollsPage int64
	if len(strArr) > minCommandParts {
		myPollsPage, err = strconv.ParseInt(strArr[3], 10, 32)
		if err != nil {
			return commandArgs{}, fmt.Errorf("unable to parse myPollsPage: %s", err.Error())
		}
	}

	return commandArgs{
		pollID:      int32(pollID),
		field:       editfield.Field(fieldID),
		myPollsPage: int32(myPollsPage),
	}, nil
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editpoll.EditPoll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetPollId() == 0 {
		return s.createPoll("", req.ReferrerMyPollsPage,
			message.MessageID, chat.ID, true)
	}

	return s.editPoll(ctx, "", req.PollId, req.ReferrerMyPollsPage,
		message.MessageID, chat.ID, upd.CallbackQuery.From.ID, true)
}

func (s *Service) editPoll(
	ctx context.Context, validationErr string, pollID, myPollsPage *int32,
	messageID int, chatID, userID int64, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, *pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	keyboard, err := pollKeyboardMarkup(pollID, myPollsPage, p.Options)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editPoll page: %s", err.Error())
	}

	txt := editPollTxt(validationErr, p)

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func editPollTxt(validationErr string, p swagger.PollWithOptions) string {
	var sb render.StringBuilder

	if validationErr != "" {
		sb.Printf("<b>🚨🚨🚨\n%s\n🚨🚨🚨</b>\n\n", validationErr)
	}

	sb.Printf("Title:\n<b>%s</b>\n", p.Title)
	sb.Printf("\nDescription\n:<b>%s</b>\n", p.Description)
	sb.Printf("\nStart date: <b>%s</b>\n", render.FormatTime(p.Start))
	sb.Printf("Finish date: <b>%s</b>\n", render.FormatTime(p.Finish))

	for i, op := range p.Options {
		sb.Printf("	%d. %s\n", i+1, op.Title)
	}

	return sb.String()
}

func (s *Service) createPoll(
	validationErrMsg string, myPollsPage *int32, messageID int, chatID int64, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := pollKeyboardMarkup(nil, myPollsPage, nil)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for createPoll page: %s", err.Error())
	}

	txt := createPollTxt(validationErrMsg)

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard), "", nil
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func createPollTxt(validationErrMsg string) string {
	txt := `🚀 Create Your Own Poll! 🚀

Step into the creator’s chair!
Define your poll’s title, set a description, choose a start date, and more.
Your audience is waiting for your questions - let’s get started!`

	var sb render.StringBuilder
	if validationErrMsg != "" {
		sb.Printf("<b>🚨🚨🚨\n%s\n🚨🚨🚨</b>\n\n", validationErrMsg)
	}

	sb.WriteString(txt)

	return sb.String()
}

const (
	fieldsMaxCountInRow  = 3
	addOptionButtonWidth = 3
)

func pollKeyboardMarkup(pollID, myPollsPage *int32, options []*swagger.Option) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []models.EditButton[editfield.Field]{
		{"Title", editfield.Field_TITLE},
		{"Description", editfield.Field_DESCRIPTION},
		{"Start date", editfield.Field_START_DATE},
		{"Finish date", editfield.Field_FINISH_DATE},
	}

	buttonsCount := len(editButtons) + 1 // +1 for Go back button

	fieldsKeyboardBuilder := render.NewKeyboardBuilder(fieldsMaxCountInRow, buttonsCount)

	for i := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditFieldRoute, &editfield.EditField{
			PollId:              pollID,
			Field:               &editButtons[i].Field,
			ReferrerMyPollsPage: myPollsPage,
		})

		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal callback data for %s button: %s", editButtons[i].Text, err.Error())
		}

		fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         editButtons[i].Text,
			CallbackData: callbackData,
		})
	}

	var currentPage *int32
	if myPollsPage != nil {
		currentPage = myPollsPage
	} else {
		currentPage = helpers.Ptr[int32](1)
	}

	goBackData, err := proto.MarshalCallbackData(models.MyPollsRoute, &mypolls.MyPolls{
		CurrentPage: currentPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for go back button: %s", err.Error())
	}

	fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Go back",
		CallbackData: goBackData,
	})

	optionsKeyboardBuilder := render.NewKeyboardBuilder(models.MaxCountInRow, len(options)+addOptionButtonWidth)

	for i, op := range options {
		var editOptionData *string
		editOptionData, err = proto.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
			PollId:              pollID,
			OptionId:            helpers.Ptr(int32(op.ID)),
			ReferrerMyPollsPage: myPollsPage,
		})
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal callback data for edit option button: %s", err.Error())
		}

		optionsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: editOptionData,
		})
	}

	addOptionData, err := proto.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
		PollId:              pollID,
		ReferrerMyPollsPage: myPollsPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for add option button: %s", err.Error())
	}

	optionsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Add option",
		CallbackData: addOptionData,
	})

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: append(fieldsKeyboardBuilder.Rows(), optionsKeyboardBuilder.Rows()...),
	}, nil
}
