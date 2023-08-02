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

type Service struct {
	db    db
	tgBot *tgbotapi.BotAPI
}

func NewService(db db, tgBot *tgbotapi.BotAPI) *Service {
	return &Service{db: db, tgBot: tgBot}
}

type db interface {
	GetSeriesByID(ctx context.Context, id int32) (models.Series, error)
	GetPollByID(ctx context.Context, id int32) (models.PollWithOptions, error)

	CreateSeries(ctx context.Context, s models.CreateSeries) (models.Series, error)
	CreatePoll(ctx context.Context, poll models.CreatePoll) (models.Poll, error)
	CreateOption(ctx context.Context, option models.CreateOption) (models.Option, error)

	UpdateSeries(ctx context.Context, id int32, s models.UpdateSeries) (models.Series, error)
	UpdatePoll(ctx context.Context, id int32, poll models.UpdatePoll) (models.Poll, error)
	UpdateOption(ctx context.Context, id int32, option models.UpdateOption) (models.Option, error)

	DeleteSeries(ctx context.Context, id int32) error
	DeletePoll(ctx context.Context, id int32) error
	DeleteOption(ctx context.Context, id int32) error
}

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.tgBot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result := s.processTGUpdate(ctx, upd)

	if result.msgText != "" {
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, result.msgText)
		if _, sendErr := s.tgBot.Send(msg); sendErr != nil {
			return fmt.Errorf("unable to send message: %s", sendErr.Error())
		}
	}

	return nil
}

type processTGResult struct {
	msgText        string
	inlineKeyboard tgbotapi.InlineKeyboardMarkup
}

func (s *Service) processTGUpdate(ctx context.Context, upd tgbotapi.Update) processTGResult {
	l := zerolog.Ctx(ctx)

	if upd.Message != nil {
		text := upd.Message.Text

		prefix := "/start showpoll"
		if strings.HasPrefix(prefix, text) {
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
				l.Error().Msgf("unable to get poll by id: %v\n", err)

				return processTGResult{
					msgText: fmt.Sprintf("oops, can't find poll with id %d", pollID),
				}
			}

			return processTGResult{
				msgText:        txtMsg(poll),
				inlineKeyboard: keyboardMarkup(poll),
			}
		}
	}

	return processTGResult{}
}

func keyboardMarkup(poll models.PollWithOptions) tgbotapi.InlineKeyboardMarkup {
	length := len(poll.Options)
	rowsCount := length / 8

	if length%8 > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i := range poll.Options {
		rowIdx := i / 8
		rows[rowIdx] = append(rows[rowIdx], tgbotapi.InlineKeyboardButton{
			Text:         strconv.Itoa(i + 1),
			CallbackData: helpers.Ptr(""),
		})
	}

	var keyboard tgbotapi.InlineKeyboardMarkup
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, rows...)
	return keyboard
}

func txtMsg(poll models.PollWithOptions) string {
	var sb strings.Builder
	sb.WriteString(pollToString(poll))
	sb.WriteString("\n")
	if time.Now().Unix() >= int64(time.Time(poll.Finish).Unix()) {
		sb.WriteString("This poll is expired!\n")
		// todo
		//if poll.OutcomeOptionId > 0 {
		//	option := poll.Options[poll.OutcomeOptionId]
		//	sb.WriteString(fmt.Sprintf("Outcome option is: #%d %q", poll.OutcomeOptionId+1, option.Title))
		//}
	}
	// todo
	//if votedOptionID > -1 {
	//	sb.WriteString(fmt.Sprintf("You have already voted option %d\n", votedOptionID+1))
	//}
	return sb.String()
}

func pollToString(p models.PollWithOptions) string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "<b>%s</b>\n", p.Title)
	_, _ = fmt.Fprintf(&sb, "Start: %s\n", time.Unix(time.Time(p.Start).Unix(), 0))
	_, _ = fmt.Fprintf(&sb, "Finish: %s\n", time.Unix(time.Time(p.Finish).Unix(), 0))
	if timeToGo := time.Time(p.Finish).Sub(time.Now()); timeToGo > 0 {
		_, _ = fmt.Fprintf(&sb, "%d days %d hours to go\n", int(timeToGo/3600)/24, int(timeToGo/3600)%24)
	} else {
		_, _ = fmt.Fprintf(&sb, "Poll has ended %s\n", time.Unix(time.Time(p.Finish).Unix(), 0))
	}

	for i, op := range p.Options {
		_, _ = fmt.Fprintf(&sb, "	%d. %s\n", i+1, op.Title)
	}

	return sb.String()
}
