package poll

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/userpollresult"
	proto2 "google.golang.org/protobuf/proto"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/votepreview"
	"github.com/babadro/forecaster/internal/domain"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto2.Message, *poll.Poll) {
	v := new(poll.Poll)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	pollIDStr := upd.Message.Text[len(models.ShowPollStartCommandPrefix):]

	pollID, err := strconv.ParseInt(pollIDStr, 10, 32)
	if err != nil {
		return nil,
			fmt.Sprintf("Oops, can't parse poll id %s", pollIDStr),
			fmt.Errorf("unable to parse poll id: %s", err.Error())
	}

	chat := upd.Message.Chat

	user := upd.Message.From
	if user == nil {
		return nil, "", fmt.Errorf("user is nil")
	}

	return s.render(
		ctx, int32(pollID), user.ID, chat.ID,
		upd.Message.MessageID, false)
}

func (s *Service) RenderCallback(
	ctx context.Context, req *poll.Poll, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	user := upd.CallbackQuery.From
	if user == nil {
		return nil, "", errors.New("user is nil")
	}

	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	return s.render(ctx, *req.PollId, user.ID, chat.ID, message.MessageID, true)
}

const (
	cantShowPollMsg = "Sorry, something went wrong, I can't show this poll right now"
)

func (s *Service) render(
	ctx context.Context, pollID int32, userID int64, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, err := s.db.GetPollByID(ctx, pollID)

	if err != nil {
		return nil,
			fmt.Sprintf("oops, can't find poll with id %d", pollID),
			fmt.Errorf("unable to get poll by id: %s", err.Error())
	}

	userAlreadyVoted := false

	lastVote, err := s.db.GetUserVote(ctx, userID, p.ID)
	if err == nil {
		userAlreadyVoted = true
	} else if !errors.Is(err, domain.ErrNotFound) {
		return nil,
			cantShowPollMsg,
			fmt.Errorf("unable to get last vote: %s", err.Error())
	}

	txt, err := txtMsg(p, userAlreadyVoted, lastVote)
	if err != nil {
		return nil,
			cantShowPollMsg,
			fmt.Errorf("unable to create text message: %s", err.Error())
	}

	keyboard, err := keyboardMarkup(p, userID)
	if err != nil {
		return nil,
			cantShowPollMsg,
			fmt.Errorf("unable to create keyboard markup: %s", err.Error())
	}

	var res tgbotapi.Chattable
	if editMessage {
		res = render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard)
	} else {
		res = render.NewMessageWithKeyboard(chatID, txt, keyboard)
	}

	return res, "", nil
}

func keyboardMarkup(poll swagger.PollWithOptions, userID int64) (tgbotapi.InlineKeyboardMarkup, error) {
	length := len(poll.Options)
	if swagger.HasOutcome(poll.Options) {
		length++ // ++ for "show results" button
	}

	rowsCount := length / models.MaxCountInRow

	if length%models.MaxCountInRow > 0 {
		rowsCount++
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, rowsCount)

	for i := 0; i < length; i++ {
		if i == len(poll.Options) {
			showMyResultsData, err := proto.MarshalCallbackData(models.UserPollResultRoute, &userpollresult.UserPollResult{
				UserId: helpers.Ptr[int64](userID),
				PollId: helpers.Ptr[int32](poll.ID),
			})
			if err != nil {
				return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal user poll result callback data: %w", err)
			}

			rowIdx := i / models.MaxCountInRow

			rows[rowIdx] = append(rows[rowIdx], tgbotapi.InlineKeyboardButton{
				Text:         "Show Results",
				CallbackData: showMyResultsData,
			})

			break
		}

		op := poll.Options[i]
		votePreview := votepreview.VotePreview{
			PollId:   helpers.Ptr(poll.ID),
			OptionId: helpers.Ptr(int32(op.ID)),
		}

		callbackData, err := proto.MarshalCallbackData(models.VotePreviewRoute, &votePreview)
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

func txtMsg(p swagger.PollWithOptions, userAlreadyVoted bool, lastVote swagger.Vote) (string, error) {
	var sb render.StringBuilder

	start, finish := render.FormatTime(p.Start), render.FormatTime(p.Finish)

	sb.Printf("<b>%s</b>\n", p.Title)
	sb.Printf("<i>Start Date: %s</i>\n", start)
	sb.Printf("<i>End Date: %s</i>\n", finish)
	sb.Printf("\n")

	timeToGo := time.Until(time.Time(p.Finish)).Seconds()
	if timeToGo > 0 {
		days := int(timeToGo/models.Seconds3600) / models.Hours24
		hours := int(timeToGo/models.Seconds3600) % models.Hours24

		sb.Printf(
			"<b>%d days %d hours minutes to go</b>\n", days, hours,
		)
	} else {
		sb.Printf("<b>Poll Status: Ended %s</b>\n", finish)
	}

	sb.WriteString("\n")

	sb.WriteString("<b>Options:</b>\n")

	for i, op := range p.Options {
		sb.Printf("	%d. %s\n", i+1, op.Title)
	}

	sb.WriteString("\n")

	if timeToGo <= 0 {
		sb.WriteString("<b>This poll has expired!</b>\n")
	}

	if userAlreadyVoted {
		votedOption, idx := swagger.FindOptionByID(p.Options, lastVote.OptionID)
		if idx == -1 {
			return "", fmt.Errorf("unable to find voted option %d for poll %d", lastVote.OptionID, p.ID)
		}

		sb.Printf("<b>Last time you voted for: %d. </b> %s\n", idx, votedOption.Title)
	}

	return sb.String(), nil
}
