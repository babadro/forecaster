package errorpage

import (
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

func ErrorPage(logger *zerolog.Logger, txtMsg string, upd tgbotapi.Update) tgbotapi.Chattable {
	if chat := upd.FromChat(); chat != nil {
		msg := render.NewMessageWithKeyboard(
			chat.ID,
			txtMsg,
			tgbotapi.InlineKeyboardMarkup{},
		)

		callbackData, err := proto.MarshalCallbackData(models.MainPageRoute, &mainpage.MainPage{})
		if err != nil {
			logger.Error().Err(err).Msg("errorPage: cant marshal callback data to main route")
		} else {
			msg.ReplyMarkup = render.Keyboard(tgbotapi.InlineKeyboardButton{
				Text:         "Back to main",
				CallbackData: callbackData,
			})
		}

		return msg
	}

	return nil
}
