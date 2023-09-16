package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
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

func (sb *StringBuilder) Print(a ...any) {
	_, _ = fmt.Fprint(sb, a...)
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

func GetHighestTimeUnit(d time.Duration) (int, string) {
	switch {
	case d.Hours()/models.Hours24/models.Days365 >= 1:
		return int(d.Hours() / models.Hours24 / models.Days365), "years"
	case d.Hours()/models.Hours24 >= 1:
		return int(d.Hours() / models.Hours24), "days"
	case d.Hours() >= 1:
		return int(d.Hours()), "hours"
	case d.Minutes() >= 1:
		return int(d.Minutes()), "minutes"
	default:
		return int(d.Seconds()), "seconds"
	}
}
