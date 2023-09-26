package poll

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/dbwrapper"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecast"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecasts"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/poll"
	proto2 "google.golang.org/protobuf/proto"

	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/helpers"
	"github.com/babadro/forecaster/internal/models/swagger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	db models.DB
	w  dbwrapper.Wrapper
}

func New(db models.DB) *Service {
	return &Service{db: db, w: dbwrapper.New(db)}
}

func (s *Service) NewRequest() (proto2.Message, *forecast.Forecast) {
	v := new(forecast.Forecast)

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
		ctx, int32(pollID), 1, user.ID, chat.ID,
		upd.Message.MessageID, false)
}

func (s *Service) RenderCallback(
	ctx context.Context, req *forecast.Forecast, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	user := upd.CallbackQuery.From
	if user == nil {
		return nil, "", errors.New("user is nil")
	}

	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	forecastsPage := req.GetReferrerForecastsPage()
	if forecastsPage == 0 {
		forecastsPage = 1
	}

	return s.render(ctx, req.GetPollId(), forecastsPage, user.ID, chat.ID, message.MessageID, true)
}

const (
	cantShowForecastMsg = "Sorry, something went wrong, I can't show this forecast right now"
)

func (s *Service) render(
	ctx context.Context, pollID, referrerForecastsPage int32, userID int64, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	p, errMsg, err := s.w.GetPollByID(ctx, pollID)
	if err != nil {
		return nil, errMsg, err
	}

	userVote, userVoteFound, err := s.w.GetUserVote(ctx, userID, p.ID)
	if err != nil {
		return nil, "", err
	}

	if swagger.HasOutcome(p.Options) {
		// edge case when polls outcome became known just before user chose this forecast
		// we need to suggest to see results instead of forecast
		return nil, "This poll already has outcome", fmt.Errorf("poll %d already has outcome", p.ID)
	}

	markup, err := keyboardMarkup(p.ID)
	if err != nil {
		return nil, "", fmt.Errorf("userpoll result: unable to create keyboard markup: %s", err.Error())
	}

	var res tgbotapi.Chattable
	msg := s.txtMsg(txtInputModel)

	if editMessage {
		res = render.NewEditMessageTextWithKeyboard(chatID, messageID, msg, markup)
	} else {
		res = render.NewMessageWithKeyboard(chatID, msg, markup)
	}

	return res, "", nil
}

func keyboardMarkup(pollID, forecastsPage int32) (tgbotapi.InlineKeyboardMarkup, error) {
	forecastsData, err := proto.MarshalCallbackData(models.ForecastsRoute,
		&forecasts.Forecasts{
			CurrentPage: helpers.Ptr(forecastsPage),
		},
	)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall poll callback data: %s", err.Error())
	}

	forecastsBtn := tgbotapi.InlineKeyboardButton{Text: "All Forecasts", CallbackData: forecastsData}

	pollData, err := proto.MarshalCallbackData(models.PollRoute,
		&poll.Poll{PollId: helpers.Ptr(pollID)},
	)
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall poll callback data: %s", err.Error())
	}

	pollBtn := tgbotapi.InlineKeyboardButton{Text: "Show Poll", CallbackData: pollData}

	return render.Keyboard(forecastsBtn, pollBtn), nil
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
