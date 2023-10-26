package render

import (
	"fmt"
	"strconv"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	proto2 "google.golang.org/protobuf/proto"
)

type ManyItemsKeyboardInput struct {
	IDs                       []int32
	CurrentPage               int32
	Prev, Next                bool
	AllItemsRoute             byte
	SingleItemRoute           byte
	AllItemsProtoMessage      func(page int32) proto2.Message
	SingleItemProtoMessage    func(itemID, referrerAllItemsPage int32) proto2.Message
	FirstRowAdditionalButtons []tgbotapi.InlineKeyboardButton
}

func ManyItemsKeyboardMarkup(in ManyItemsKeyboardInput) (tgbotapi.InlineKeyboardMarkup, error) {
	var firstRow []tgbotapi.InlineKeyboardButton

	var err error

	if firstRow, err = appendMainMenuButton(firstRow); err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	firstRow, err = appendNaviButton(firstRow, in.AllItemsRoute, in.AllItemsProtoMessage,
		in.Prev, in.CurrentPage-1, "Prev")
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	firstRow, err = appendNaviButton(firstRow, in.AllItemsRoute, in.AllItemsProtoMessage,
		in.Next, in.CurrentPage+1, "Next")
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

		pollData, err = proto.MarshalCallbackData(in.SingleItemRoute, in.SingleItemProtoMessage(id, in.CurrentPage))
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
	row []tgbotapi.InlineKeyboardButton,
	route byte,
	protoMessage func(page int32) proto2.Message,
	exists bool, page int32, name string,
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

func appendMainMenuButton(row []tgbotapi.InlineKeyboardButton) ([]tgbotapi.InlineKeyboardButton, error) {
	data, err := proto.MarshalCallbackData(models.MainPageRoute, &mainpage.MainPage{})
	if err != nil {
		return nil, fmt.Errorf("unable to marshal main callback data: %s", err.Error())
	}

	row = append(row, tgbotapi.InlineKeyboardButton{
		Text:         "Main Menu",
		CallbackData: data,
	})

	return row, nil
}

func NewKeyboardBuilder(rowLen, buttonsCount int) KeyboardBuilder {
	rowsCapacity := buttonsCount / rowLen
	if buttonsCount%rowLen > 0 {
		rowsCapacity++
	}

	return KeyboardBuilder{
		rows:   make([][]tgbotapi.InlineKeyboardButton, 0, rowsCapacity),
		rowLen: rowLen,
	}
}

type KeyboardBuilder struct {
	rows         [][]tgbotapi.InlineKeyboardButton
	rowLen       int
	buttonsCount int
}

func (k *KeyboardBuilder) AddButton(button tgbotapi.InlineKeyboardButton) {
	rowIdx := k.buttonsCount / k.rowLen
	if rowIdx == len(k.rows) {
		k.rows = append(k.rows, make([]tgbotapi.InlineKeyboardButton, 0, k.rowLen))
	}

	k.rows[rowIdx] = append(k.rows[rowIdx], button)
	k.buttonsCount++
}

func (k *KeyboardBuilder) MarkUp() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: k.rows}
}

func (k *KeyboardBuilder) Rows() [][]tgbotapi.InlineKeyboardButton {
	return k.rows
}
