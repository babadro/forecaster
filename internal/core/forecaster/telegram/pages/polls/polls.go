package polls

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto2.Message, *polls.Polls) {
	v := new(polls.Polls)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	currentPageStr := upd.Message.Text[len(models.ShowPollsStartCommandPrefix):]

	_ = currentPageStr

	return nil, "", nil
	//	return s.render()
}

func (s *Service) RenderCallback(
	ctx context.Context, req *polls.Polls, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	return s.render(ctx, req.GetCurrentPage(), chat.ID, message.MessageID, true)
}

const pageSize = 10

func (s *Service) render(
	ctx context.Context, currentPage int32, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	pollsArr, totalCount, err := s.db.GetPolls(ctx, currentPage, pageSize)
	if err != nil {
		errMsgForUser, errToLog := "", fmt.Errorf("unable to get polls: %s", err.Error())
		if errors.Is(err, domain.ErrNotFound) {
			errMsgForUser = "Polls not found"
		}

		return nil, errMsgForUser, errToLog
	}

	keyboardIn := keyboardInput{
		pollsArr:    pollsArr,
		currentPage: currentPage,
		prev:        currentPage > 1,
		next:        currentPage*pageSize < totalCount,
	}

	keyboard, err := keyboardMarkup(keyboardIn)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txtMsg(pollsArr), keyboard),
			"", nil
	}

	return render.NewMessageWithKeyboard(chatID, txtMsg(pollsArr), keyboard), "", nil
}

func txtMsg(pollsArr []swagger.Poll) string {
	var sb render.StringBuilder
	for i, p := range pollsArr {
		sb.Printf("%d. %s\n", i+1, p.Title)
	}

	return sb.String()
}

type keyboardInput struct {
	pollsArr    []swagger.Poll
	currentPage int32
	prev, next  bool
}

func keyboardMarkup(in keyboardInput) (tgbotapi.InlineKeyboardMarkup, error) {
	var firstRow []tgbotapi.InlineKeyboardButton

	var err error
	firstRow, err = appendNaviButton(firstRow, in.prev, in.currentPage-1, "Prev")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	firstRow, err = appendNaviButton(firstRow, in.next, in.currentPage+1, "Next")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	if len(firstRow) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, firstRow)
	}

	rowsCount := len(in.pollsArr) / models.MaxCountInRow
	if len(in.pollsArr)%models.MaxCountInRow > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i := range in.pollsArr {
		p := in.pollsArr[i]
		var pollData *string
		pollData, err = proto.MarshalCallbackData(models.PollRoute, &poll.Poll{
			PollId: helpers.Ptr(p.ID),
		})
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{},
				fmt.Errorf("unable to marshal poll callback data: %s", err.Error())
		}

		rowIdx := i / models.MaxCountInRow

		rows[rowIdx] = append(rows[rowIdx], tgbotapi.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: pollData,
		})
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows...)

	return keyboard, nil
}

func appendNaviButton(row []tgbotapi.InlineKeyboardButton, exists bool, page int32, name string) ([]tgbotapi.InlineKeyboardButton, error) {
	if !exists {
		return row, nil
	}

	data, err := proto.MarshalCallbackData(models.PollsRoute, &polls.Polls{
		CurrentPage: helpers.Ptr(page),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal %s callback data: %s", name, err.Error())
	}

	row = append(row, tgbotapi.InlineKeyboardButton{
		Text:         name,
		CallbackData: data,
	})

	return row, nil
}
