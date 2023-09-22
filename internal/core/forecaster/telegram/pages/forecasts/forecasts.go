package forecasts

import (
	"context"
	"fmt"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecasts"
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

	keyboardIn := keyboardInput{
		pollsArr:    pollsArr,
		currentPage: currentPage,
		prev:        currentPage > 1,
		next:        currentPage*pageSize < totalCount,
	}

	keyboard, err := keyboardMarkup(keyboardIn)
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	if editMessage {
		return render.NewEditMessageTextWithKeyboard(chatID, messageID, txtMsg(pollsArr), keyboard),
			"", nil
	}

	return render.NewMessageWithKeyboard(chatID, txtMsg(pollsArr), keyboard), "", nil

}
