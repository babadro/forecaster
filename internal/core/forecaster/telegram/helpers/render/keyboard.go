package render

import (
	"fmt"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"

	proto2 "google.golang.org/protobuf/proto"
)

type keyboardInput struct {
	ids          []int32
	currentPage  int32
	prev, next   bool
	route        byte
	protoMessage func(page int32) proto2.Message
}

func keyboardMarkup(in keyboardInput) (tgbotapi.InlineKeyboardMarkup, error) {
	var firstRow []tgbotapi.InlineKeyboardButton

	var err error

	firstRow, err = appendNaviButton(in.route, in.protoMessage,
		firstRow, in.prev, in.currentPage-1, "Prev")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	firstRow, err = appendNaviButton(in.route, in.protoMessage,
		firstRow, in.next, in.currentPage+1, "Next")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	if len(firstRow) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, firstRow)
	}

	rowsCount := len(in.ids) / models.MaxCountInRow
	if len(in.ids)%models.MaxCountInRow > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i, id := range in.ids {

		var pollData *string

		pollData, err = proto.MarshalCallbackData(models.PollRoute, &poll.Poll{
			PollId: helpers.Ptr(id),
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

func appendNaviButton(
	route byte,
	protoMessage func(page int32) proto2.Message,
	row []tgbotapi.InlineKeyboardButton, exists bool, page int32, name string,
) ([]tgbotapi.InlineKeyboardButton, error) {
	if !exists {
		return row, nil
	}

	data, err := proto.MarshalCallbackData(route, protoMessage(page))
	if err != nil {
		return nil, fmt.Errorf("unable to marshal %s callback data: %s", name, err.Error())
	}

	row = append(row, tgbotapi.InlineKeyboardButton{
		Text:         name,
		CallbackData: data,
	})

	return row, nil
}
