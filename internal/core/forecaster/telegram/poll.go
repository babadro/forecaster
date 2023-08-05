package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	models "github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

const (
	maxRows     = 8
	hours24     = 24
	seconds3600 = 3600
)

func (s *Service) poll(ctx context.Context, pollIDStr string) processTGResult {
	l := zerolog.Ctx(ctx)

	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		l.Error().Msgf("unable to convert poll id to int: %v\n", err)

		return processTGResult{
			msgText: fmt.Sprintf("invalid poll id: %s", pollIDStr),
		}
	}

	poll, err := s.db.GetPollByID(ctx, int32(pollID))

	if err != nil {
		l.Error().Int64("id", pollID).Msgf("unable to get poll by id: %v\n", err)

		return processTGResult{
			msgText: fmt.Sprintf("oops, can't find poll with id %d", pollID),
		}
	}

	return processTGResult{
		msgText:        txtMsg(poll),
		inlineKeyboard: keyboardMarkup(poll),
	}
}

func keyboardMarkup(poll models.PollWithOptions) tgbotapi.InlineKeyboardMarkup {
	length := len(poll.Options)
	rowsCount := length / maxRows

	if length%maxRows > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i := range poll.Options {
		rowIdx := i / maxRows
		rows[rowIdx] = append(rows[rowIdx], tgbotapi.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: helpers.Ptr(""),
		})
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows...)

	return keyboard
}

func txtMsg(p models.PollWithOptions) string {
	var sb strings.Builder

	start, finish := formatTime(p.Start), formatTime(p.Finish)

	fPrintf(&sb, "<b>%s</b>\n", p.Title)
	fPrintf(&sb, "<i>Start Date: %s</i>\n", start)
	fPrintf(&sb, "<i>End Date: %s</i>\n", finish)
	fPrintf(&sb, "\n")

	timeToGo := time.Until(time.Time(p.Finish))
	if timeToGo > 0 {
		fPrintf(&sb, "<b>%d days %d hours to go</b>\n", int(timeToGo/seconds3600)/hours24, int(timeToGo/seconds3600)%hours24)
	} else {
		fPrintf(&sb, "<b>Poll Status: Ended %s</b>\n", finish)
	}

	fPrintf(&sb, "\n")

	fPrint(&sb, "<b>Options:</b>\n")

	for i, op := range p.Options {
		fPrintf(&sb, "	%d. %s\n", i+1, op.Title)
	}

	fPrint(&sb, "\n")

	if timeToGo <= 0 {
		fPrint(&sb, "<b>This poll has expired!</b>\n")
	}

	return sb.String()
}
