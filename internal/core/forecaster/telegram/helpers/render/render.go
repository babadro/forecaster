package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StringBuilder struct {
	strings.Builder
}

func (sb *StringBuilder) WriteString(s string) {
	_, _ = sb.Builder.WriteString(s)
}

func (sb *StringBuilder) WriteStringLn(s string) {
	_, _ = sb.Builder.WriteString(s)
	_, _ = sb.Builder.WriteString("\n")
}

func (sb *StringBuilder) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(sb, format, a...)
}

func FormatTime[T time.Time | strfmt.DateTime](t T) string {
	return time.Time(t).Format(time.RFC822)
}

func NewMessageWithKeyboard(
	chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup,
) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	msg.ParseMode = tgbotapi.ModeHTML

	return msg
}

func NewEditMessageTextWithKeyboard(
	chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	msg.ReplyMarkup = &keyboard
	msg.ParseMode = tgbotapi.ModeHTML

	return msg
}

func Keyboard(buttons ...tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
		buttons,
	}}
}

func FindOptionByID(options []*swagger.Option, id int16) (*swagger.Option, int) {
	for i, op := range options {
		if op.ID == id {
			return op, i
		}
	}

	return nil, -1
}
