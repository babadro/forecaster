package forecasts

import (
	"context"
	"fmt"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecasts"
	"github.com/babadro/forecaster/internal/helpers"
	models2 "github.com/babadro/forecaster/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	proto2 "google.golang.org/protobuf/proto"
	"strconv"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto2.Message, *forecasts.Forecasts) {
	v := new(forecasts.Forecasts)

	return v, v
}

func (s *Service) RenderStartCommand(ctx context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	currentPageStr := upd.Message.Text[len(models.ShowForecastsStartCommand):]

	currentPage, err := strconv.ParseInt(currentPageStr, 10, 32)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse current page: %s", err.Error())
	}

	return s.render(ctx, int32(currentPage), upd.Message.Chat.ID, upd.Message.MessageID, false)
}

func (s *Service) RenderCallback(
	ctx context.Context, req *forecasts.Forecasts, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	chat := upd.CallbackQuery.Message.Chat
	message := upd.CallbackQuery.Message

	return s.render(ctx, req.GetCurrentPage(), chat.ID, message.MessageID, true)
}

const pageSize = 10

func (s *Service) render(
	ctx context.Context, currentPage int32, chatID int64, messageID int, editMessage bool,
) (tgbotapi.Chattable, string, error) {
	offset, limit := uint64((currentPage-1)*pageSize), uint64(pageSize)

	forecastArr, totalCount, err := s.db.GetForecasts(ctx, offset, limit)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get forecasts: %s", err.Error())
	}

	keyboardIn := render.KeyboardInput{
		IDs:         forecastIDs(forecastArr),
		CurrentPage: currentPage,
		Prev:        currentPage > 1,
		Next:        currentPage*pageSize < totalCount,
		Route:       models.ForecastsRoute,
		ProtoMessage: func(page int32) proto2.Message {
			return &forecasts.Forecasts{CurrentPage: helpers.Ptr(page)}
		},
	}

	keyboard, err := render.KeyboardMarkup(keyboardIn)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	txt, err := txtMsg(forecastArr)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create text message: %s", err.Error())
	}

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txt, keyboard),
			"", nil
	}

	return render.NewMessageWithKeyboard(chatID, txt, keyboard), "", nil
}

func txtMsg(forecastsArr []models2.Forecast) (string, error) {
	if len(forecastsArr) == 0 {
		return "There are no polls yet", nil
	}

	var sb render.StringBuilder
	for i, f := range forecastsArr {
		sb.Printf("%d. %s\n", i+1, f.PollTitle)

		s, err := calculateStatistic(f.Options)
		if err != nil {
			return "", fmt.Errorf(
				"unable to calculate options statistics for pollID %d: %s", f.PollID, err.Error())
		}

		sb.Printf("<b>%s</b>\n", f.PollTitle)
		sb.Printf("Most popular option (%d%% votes):\n", s.topOptionPercentage)
		sb.Printf("<b>%s</b>\n", s.topOption.Title)
	}

	return sb.String(), nil
}

type stat struct {
	topOption           models2.ForecastOption
	totalVotes          int32
	topOptionPercentage int
}

func calculateStatistic(options []models2.ForecastOption) (stat, error) {
	if len(options) == 0 {
		return stat{}, fmt.Errorf("options is empty")
	}

	topOptionIDx, total := 0, 0

	for i, f := range options {
		total++

		if f.TotalVotes > options[topOptionIDx].TotalVotes {
			topOptionIDx = i
		}
	}

	if total == 0 {
		return stat{}, fmt.Errorf("total is zero")
	}

	return stat{
		topOption:  options[topOptionIDx],
		totalVotes: int32(total),
	}, nil
}

func forecastIDs(forecastArr []models2.Forecast) []int32 {
	res := make([]int32, len(forecastArr))

	for i, forecast := range forecastArr {
		res[i] = forecast.PollID
	}

	return res
}
