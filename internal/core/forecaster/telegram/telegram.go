package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/babadro/forecaster/internal/helpers"
	models "github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/strfmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	db  db
	bot tgBot
}

func NewService(db db, b tgBot) *Service {
	return &Service{db: db, bot: b}
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

type tgBot interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func (s *Service) ProcessTelegramUpdate(logger *zerolog.Logger, upd tgbotapi.Update) error {
	if s.bot == nil {
		return fmt.Errorf("telegram bot is not initialized")
	}

	ctx := logger.WithContext(context.Background())

	result := s.processTGUpdate(ctx, upd)

	if result.msgText != "" {
		logger.Info().Msg(result.msgText)

		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, result.msgText)
		msg.ParseMode = "HTML"
		if _, sendErr := s.bot.Send(msg); sendErr != nil {
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

		prefix := "/start showpoll_"
		if strings.HasPrefix(text, prefix) {
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

func txtMsg(p models.PollWithOptions) string {
	var sb strings.Builder

	start, finish := formatTime(p.Start), formatTime(p.Finish)

	fPrintf(&sb, "<b>%s</b>\n", p.Title)
	fPrintf(&sb, "<i>Start Date: %s</i>\n", start)
	fPrintf(&sb, "<i>End Date: %s</i>\n", finish)
	fPrintf(&sb, "\n")

	timeToGo := time.Time(p.Finish).Sub(time.Now())
	if timeToGo > 0 {
		fPrintf(&sb, "<b>%d days %d hours to go</b>\n", int(timeToGo/3600)/24, int(timeToGo/3600)%24)
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
