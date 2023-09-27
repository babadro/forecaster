package render

import (
	"fmt"
	"strconv"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	proto2 "google.golang.org/protobuf/proto"
)

type KeyboardInput struct {
	IDs                    []int32
	CurrentPage            int32
	Prev, Next             bool
	AllItemsRoute          byte
	SingleItemRoute        byte
	AllItemsProtoMessage   func(page int32) proto2.Message
	SingleItemProtoMessage func(itemID int32) proto2.Message
}

func KeyboardMarkup(in KeyboardInput) (tgbotapi.InlineKeyboardMarkup, error) {
	var firstRow []tgbotapi.InlineKeyboardButton

	var err error

	firstRow, err = appendNaviButton(in.AllItemsRoute, in.AllItemsProtoMessage,
		firstRow, in.Prev, in.CurrentPage-1, "Prev")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	firstRow, err = appendNaviButton(in.AllItemsRoute, in.AllItemsProtoMessage,
		firstRow, in.Next, in.CurrentPage+1, "Next")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	if len(firstRow) > 0 {
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, firstRow)
	}

	rowsCount := len(in.IDs) / models.MaxCountInRow
	if len(in.IDs)%models.MaxCountInRow > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i, id := range in.IDs {
		var pollData *string

		pollData, err = proto.MarshalCallbackData(in.SingleItemRoute, in.SingleItemProtoMessage(id))
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
