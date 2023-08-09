package poll

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	helpers2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Poll(ctx context.Context, pollIDStr string, userID int64, scope models.Scope) (models.ProcessTgResult, string, error) {
	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		return models.ProcessTgResult{},
			fmt.Sprintf("Oops, can't parse poll id %s", pollIDStr),
			fmt.Errorf("unable to parse poll id: %s", err.Error())
	}

	poll, err := scope.DB.GetPollByID(ctx, int32(pollID))

	if err != nil {
		return models.ProcessTgResult{},
			fmt.Sprintf("oops, can't find poll with id %d", pollID),
			fmt.Errorf("unable to get poll by id: %s\n", err.Error())
	}

	keyboard, err := keyboardMarkup(poll)
	if err != nil {
		return models.ProcessTgResult{},
			"Sorry, something went wrong, I can't show this option right now",
			fmt.Errorf("unable to create keyboard markup: %s\n", err.Error())
	}

	return models.ProcessTgResult{
		MsgText:        txtMsg(poll),
		InlineKeyboard: keyboard,
	}, "", nil
}

func keyboardMarkup(poll swagger.PollWithOptions) (tgbotapi.InlineKeyboardMarkup, error) {
	length := len(poll.Options)
	rowsCount := length / models.MaxCountInRow

	if length%models.MaxCountInRow > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i, op := range poll.Options {
		votePreview := votepreview.VotePreview{
			PollId:   helpers.Ptr(poll.ID),
			OptionId: helpers.Ptr[int32](int32(op.ID)),
		}

		callbackData, err := helpers2.CallbackData(models.VotePreviewRoute, &votePreview)
		if err != nil {
			return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to create callback data: %w", err)
		}

		rowIdx := i / models.MaxCountInRow
		rows[rowIdx] = append(rows[rowIdx], tgbotapi.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: callbackData,
		})
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows...)

	return keyboard, nil
}

func txtMsg(p swagger.PollWithOptions) string {
	var sb strings.Builder

	start, finish := formatTime(p.Start), formatTime(p.Finish)

	fPrintf(&sb, "<b>%s</b>\n", p.Title)
	fPrintf(&sb, "<i>Start Date: %s</i>\n", start)
	fPrintf(&sb, "<i>End Date: %s</i>\n", finish)
	fPrintf(&sb, "\n")

	timeToGo := time.Until(time.Time(p.Finish))
	if timeToGo > 0 {
		fPrintf(
			&sb, "<b>%d days %d hours to go</b>\n",
			int(timeToGo/models.Seconds3600)/models.Hours24, int(timeToGo/models.Seconds3600)%models.Hours24,
		)
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

func formatTime[T time.Time | strfmt.DateTime](t T) string {
	return time.Time(t).Format(time.RFC822)
}

func fPrintf(sb *strings.Builder, format string, a ...any) {
	_, _ = fmt.Fprintf(sb, format, a...)
}

func fPrint(sb *strings.Builder, a ...any) {
	_, _ = fmt.Fprint(sb, a...)
}
