package poll

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Result struct {
	MsgText        string
	InlineKeyboard tgbotapi.InlineKeyboardMarkup
}

func Process(pollID) (Result, error) {
	pollIDStr := strings.TrimPrefix(text, prefix)

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
