package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Tg struct {
	tgBot *tgbotapi.BotAPI
	wg    *sync.WaitGroup
}

func NewTelegram(tgBot *tgbotapi.BotAPI) *Tg {
	return &Tg{tgBot: tgBot, wg: &sync.WaitGroup{}}
}

func (t *Tg) ReceiveUpdates(params operations.ReceiveTelegramUpdatesParams) middleware.Responder {
	var update tgbotapi.Update
	err := json.NewDecoder(params.HTTPRequest.Body).Decode(&update)
	if err != nil {
		return operations.NewReceiveTelegramUpdatesBadRequest().WithPayload(&swagger.Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to decode update: %v", err),
		})
	}

	t.wg.Add(1)

	go func() {
		defer t.wg.Done()

		go processUpdate(update, t.tgBot)
	}()

	return operations.NewReceiveTelegramUpdatesOK()
}

func (t *Tg) Wait() {
	t.wg.Wait()
}

func processUpdate(upd tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(upd.Message.Chat.ID, upd.Message.Text)
	if _, sendErr := bot.Send(msg); sendErr != nil {
		log.Printf("Unable to send message: %v\n", sendErr)
	}

	// marshal update to json and output to stdout
	updJSON, err := json.Marshal(upd)
	if err != nil {
		log.Printf("Unable to marshal update: %v\n", err)
		return
	}

	log.Printf("Update: %s\n", updJSON)
}
