package mainpage

import (
	"context"
	"fmt"

	proto2 "github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/proto"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/helpers/render"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/models"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/forecasts"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mainpage"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/mypolls"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram/proto/polls"
	"github.com/babadro/forecaster/internal/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	db models.DB
}

func New(db models.DB) *Service {
	return &Service{db: db}
}

func (s *Service) NewRequest() (proto.Message, *mainpage.MainPage) {
	v := new(mainpage.MainPage)

	return v, v
}

func (s *Service) RenderStartCommand(_ context.Context, upd tgbotapi.Update) (tgbotapi.Chattable, string, error) {
	return s.render(upd.Message.Chat.ID, upd.Message.MessageID, false)
}

func (s *Service) RenderCallback(
	_ context.Context, _ *mainpage.MainPage, upd tgbotapi.Update,
) (tgbotapi.Chattable, string, error) {
	origMessage := upd.CallbackQuery.Message

	return s.render(origMessage.Chat.ID, origMessage.MessageID, true)
}

func (s *Service) render(chatID int64, messageID int, editMsg bool) (tgbotapi.Chattable, string, error) {
	keyboardMarkup, err := keyboard()
	if err != nil {
		return nil, "", fmt.Errorf("unable to create keyboard: %s", err.Error())
	}

	var res tgbotapi.Chattable
	if editMsg {
		res = render.NewEditMessageTextWithKeyboard(chatID, messageID, txtMsg(), keyboardMarkup)
	} else {
		res = render.NewMessageWithKeyboard(chatID, txtMsg(), keyboardMarkup)
	}

	return res, "", nil
}

func txtMsg() string {
	//nolint:lll,stylecheck // lines should not be broken
	return `🚀 Welcome to Forecaster Bot! 🚀

🗳 Your go-to platform for creating and participating in forecasts and polls directly through Telegram! 🗳

With Forecaster Bot, you can:

    Create Forecasts: 🛠 Craft your own polls, ask questions, and gather insights from your audience.
    Vote on Polls: ✅ Engage in diverse polls, and let your opinion be heard.
    Analyze Results: 📊 Dive into real-time analytics of poll outcomes.
    Discuss Outcomes: 💬 Share thoughts and discuss poll results with others.

🔄 Navigate through our menu to explore, create, or vote in ongoing polls, and become a part of our forecasting community!

🕵️‍♂️ Discover polls on various topics, and leverage collective intelligence to glimpse into the future!

Let’s dive in, and may the best forecast win! 🎉`
}

func keyboard() (tgbotapi.InlineKeyboardMarkup, error) {
	currentPage := helpers.Ptr(int32(1))
	pollsData, err := proto2.MarshalCallbackData(models.PollsRoute, &polls.Polls{
		CurrentPage: currentPage,
	})

	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable to marshal polls callback data: %w", err)
	}

	forecastsData, err := proto2.MarshalCallbackData(models.ForecastsRoute, &forecasts.Forecasts{
		CurrentPage: currentPage,
	})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall forecasts callback data: %s", err.Error())
	}

	myPollsData, err := proto2.MarshalCallbackData(models.MyPollsRoute, &mypolls.MyPolls{CurrentPage: helpers.Ptr[int32](1)})
	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, fmt.Errorf("unable marshall myPolls callback data: %s", err.Error())
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{Text: "All polls", CallbackData: pollsData},
			tgbotapi.InlineKeyboardButton{Text: "All forecasts", CallbackData: forecastsData},
			tgbotapi.InlineKeyboardButton{Text: "My polls", CallbackData: myPollsData},
		)), nil
}
