package mainpage

import (
	"context"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RenderCallback(
	_ context.Context, _ *mainpage.MainPage, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	origMessage := upd.CallbackQuery.Message

	return render.NewEditMessageTextWithKeyboard(
		origMessage.Chat.ID, origMessage.MessageID, "main page will be here",
		tgbotapi.InlineKeyboardMarkup{}), "", nil
}
