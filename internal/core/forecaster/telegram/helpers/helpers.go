package helpers

import (
	"fmt"

	"github.com/babadro/forecaster/internal/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
)

func CallbackData(route byte, m proto.Message) (*string, error) {
	binaryData, err := proto.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("can't marshal proto message: %w", err)
	}

	res := make([]byte, 0, len(binaryData)+1)
	res = append(res, route)
	res = append(res, binaryData...)

	return helpers.Ptr(string(res)), nil
}

func UnmarshalCallbackData(data string, m proto.Message) error {
	if len(data) < 2 {
		return fmt.Errorf("data is too short")
	}

	// first byte is route
	binaryData := []byte(data[1:])

	return proto.Unmarshal(binaryData, m)
}

func NewMessageWithKeyboard(
	chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup,
) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeHTML

	return msg
}
