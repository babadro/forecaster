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
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/deletepoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editoption"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpoll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/editpollfield"
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

	var (
		updateModel   swagger.UpdatePoll
		createModel   swagger.CreatePoll
		p             swagger.Poll
		validationErr string
	)

	if args.pollID == 0 {
		createModel, validationErr, err = getCreateModel(ctx, args.field, update.Message.Text, update.Message.From.ID)
		if err != nil {
			return nil, "", fmt.Errorf("unable to get create model: %s", err.Error())
		}

		if validationErr == "" {
			if p, err = s.db.CreatePoll(ctx, createModel, time.Now()); err != nil {
				return nil, "", fmt.Errorf("unable to create poll: %s", err.Error())
			}
		}

		return s.editPollDialog(ctx, validationErr, p.ID, args.myPollsPage,
			updateModel, false,
			update.Message.MessageID, update.Message.Chat.ID, update.Message.From.ID, false)
	}

	updateModel, validationErr, err = getUpdateModel(ctx, args.field, update.Message.Text, update.Message.From.ID)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get update model: %s", err.Error())
	}

	return s.editPollDialog(ctx, validationErr, args.pollID, args.myPollsPage,
		updateModel, validationErr == "",
		update.Message.MessageID, update.Message.Chat.ID, update.Message.From.ID, false)
}

const invalidDateErrMsg = "Can't parse date format.\nIt should be " + time.RFC3339

func getUpdateModel(
	ctx context.Context, fieldID editpollfield.Field, input string, telegramUserID int64,
) (swagger.UpdatePoll, string, error) {
	update := swagger.UpdatePoll{TelegramUserID: &telegramUserID}

	switch fieldID {
	case editpollfield.Field_TITLE:
		update.Title = &input
	case editpollfield.Field_DESCRIPTION:
		update.Description = &input
	case editpollfield.Field_START_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's start poll date")
			return swagger.UpdatePoll{}, invalidDateErrMsg, nil
		}

		update.Start = &date
	case editpollfield.Field_FINISH_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's finish poll date")
			return swagger.UpdatePoll{}, invalidDateErrMsg, nil
		}

		update.Finish = &date
	case editpollfield.Field_UNDEFINED:
		return swagger.UpdatePoll{}, "", fmt.Errorf("field is undefined")
	default:
		return swagger.UpdatePoll{}, "",
			fmt.Errorf("unknown field %d", fieldID)
	}

	return update, "", nil
}

func getCreateModel(
	ctx context.Context, fieldID editpollfield.Field, input string, telegramUserID int64,
) (swagger.CreatePoll, string, error) {
	create := swagger.CreatePoll{TelegramUserID: telegramUserID}

	switch fieldID {
	case editpollfield.Field_TITLE:
		create.Title = input
	case editpollfield.Field_DESCRIPTION:
		create.Description = input
	case editpollfield.Field_START_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's start poll date")
			return swagger.CreatePoll{}, invalidDateErrMsg, nil
		}

		create.Start = date
	case editpollfield.Field_FINISH_DATE:
		date, err := parseDate(input)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Msg("unable to parse user's finish poll date")
			return swagger.CreatePoll{}, invalidDateErrMsg, nil
		}

		create.Finish = date
	case editpollfield.Field_UNDEFINED:
		return swagger.CreatePoll{}, "", fmt.Errorf("field is undefined")
	}

	return create, "", nil
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
	field       editpollfield.Field
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

	fieldID, ok := editpollfield.Field_value[field]
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
		field:       editpollfield.Field(fieldID),
		myPollsPage: int32(myPollsPage),
	}, nil
}

func (s *Service) RenderCallback(
	ctx context.Context, req *editpoll.EditPoll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	if req.GetPollId() == 0 {
		return s.createPollDialog("", req.GetReferrerMyPollsPage(),
			message.MessageID, chat.ID, true)
	}

	return s.editPollDialog(ctx, "", req.GetPollId(), req.GetReferrerMyPollsPage(),
		swagger.UpdatePoll{}, false,
		message.MessageID, chat.ID, upd.CallbackQuery.From.ID, true)
}

func (s *Service) editPollDialog(
	ctx context.Context, validationErr string, pollID, myPollsPage int32,
	updateModel swagger.UpdatePoll, doUpdate bool,
	messageID int, chatID, userID int64, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	if p.TelegramUserID != userID {
		return nil, "forbidden", fmt.Errorf("user %d is not owner of poll %d", userID, pollID)
	}

	var updatedPoll swagger.Poll
	if doUpdate {
		updatedPoll, err = s.db.UpdatePoll(ctx, pollID, updateModel, time.Now())
		if err != nil {
			return nil, "", fmt.Errorf("unable to update poll: %s", err.Error())
		}

		p = swagger.MergePolls(p, updatedPoll)
	}

	keyboard, err := pollKeyboardMarkup(pollID, myPollsPage, p.Options)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for editPollDialog page: %s", err.Error())
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

	sb.WriteString("\n<b>Options:</b>\n")

	for i, op := range p.Options {
		sb.Printf("	%d. %s\n", i+1, op.Title)
	}

	return sb.String()
}

func (s *Service) createPollDialog(
	validationErrMsg string, myPollsPage int32, messageID int, chatID int64,
	editMessage bool,
) (tgbotapi.Chattable, string, error) {
	keyboard, err := pollKeyboardMarkup(0, myPollsPage, nil)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard for createPollDialog page: %s", err.Error())
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

	goBackButton = 1
	deleteButton = 1
)

func pollKeyboardMarkup(pollID, myPollsPage int32, options []*swagger.Option) (tgbotapi.InlineKeyboardMarkup, error) {
	editButtons := []models.EditButton[editpollfield.Field]{
		{Text: "Title", Field: editpollfield.Field_TITLE},
		{Text: "Description", Field: editpollfield.Field_DESCRIPTION},
		{Text: "Start date", Field: editpollfield.Field_START_DATE},
		{Text: "Finish date", Field: editpollfield.Field_FINISH_DATE},
	}

	buttonsCount := len(editButtons) + goBackButton + deleteButton

	fieldsKeyboardBuilder := render.NewKeyboardBuilder(fieldsMaxCountInRow, buttonsCount)

	pollIDPtr, myPollsPagePtr := helpers.NilIfZero(pollID), helpers.NilIfZero(myPollsPage)

	var err error

	fieldsKeyboardBuilder, err = addEditButtons(fieldsKeyboardBuilder, editButtons, pollIDPtr, myPollsPagePtr)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to add edit buttons: %s", err.Error())
	}

	goBackData, err := proto.MarshalCallbackData(models.MyPollsRoute, &mypolls.MyPolls{
		CurrentPage: helpers.OneIfZero(myPollsPage),
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for go back button: %s", err.Error())
	}

	fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Go back",
		CallbackData: goBackData,
	})

	addOptionData, err := proto.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
		PollId:              pollIDPtr,
		ReferrerMyPollsPage: myPollsPagePtr,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{},
			fmt.Errorf("unable to marshal callback data for add option button: %s", err.Error())
	}

	fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
		Text:         "Add option",
		CallbackData: addOptionData,
	})

	if pollID != 0 {
		var deleteData *string

		deleteData, err = proto.MarshalCallbackData(models.DeletePollRoute, &deletepoll.DeletePoll{
			PollId:              pollIDPtr,
			ReferrerMyPollsPage: myPollsPagePtr,
			NeedConfirmation:    helpers.Ptr(true),
		})
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal delete callback data: %s", err.Error())
		}

		fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         "Delete poll",
			CallbackData: deleteData,
		})
	}

	optionsKeyboardBuilder := render.NewKeyboardBuilder(models.MaxCountInRow, len(options)+addOptionButtonWidth)

	for i, op := range options {
		var editOptionData *string
		editOptionData, err = proto.MarshalCallbackData(models.EditOptionRoute, &editoption.EditOption{
			PollId:              pollIDPtr,
			OptionId:            helpers.Ptr(int32(op.ID)),
			ReferrerMyPollsPage: myPollsPagePtr,
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

	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: append(fieldsKeyboardBuilder.Rows(), optionsKeyboardBuilder.Rows()...),
	}, nil
}

func addEditButtons(fieldsKeyboardBuilder render.KeyboardBuilder,
	editButtons []models.EditButton[editpollfield.Field], pollIDPtr, myPollsPagePtr *int32,
) (render.KeyboardBuilder, error) {
	for i := range editButtons {
		callbackData, err := proto.MarshalCallbackData(models.EditPollFieldRoute, &editpollfield.EditPollField{
			PollId:              pollIDPtr,
			Field:               &editButtons[i].Field,
			ReferrerMyPollsPage: myPollsPagePtr,
		})

		if err != nil {
			return render.KeyboardBuilder{},
				fmt.Errorf("unable to marshal callback data for %s button: %s", editButtons[i].Text, err.Error())
		}

		fieldsKeyboardBuilder.AddButton(tgbotapi.InlineKeyboardButton{
			Text:         editButtons[i].Text,
			CallbackData: callbackData,
		})
	}

	return fieldsKeyboardBuilder, nil
}
