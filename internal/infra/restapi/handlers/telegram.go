package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/hlog"
)

func (p *Polls) ReceiveUpdates(params operations.ReceiveTelegramUpdatesParams) middleware.Responder {
	var update tgbotapi.Update
	logger := hlog.FromRequest(params.HTTPRequest)

	err := json.NewDecoder(params.HTTPRequest.Body).Decode(&update)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to decode update")

		return operations.NewReceiveTelegramUpdatesBadRequest().WithPayload(&swagger.Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to decode update: %v", err),
		})
	}

	p.wg.Add(1)

	go func() {
		defer p.wg.Done()

		if err = p.svc.ProcessTelegramUpdate(logger, update); err != nil {
			logger.Error().Err(err).
				// todo: reread request body and log it as is with Bytes method of logger
				Msg("Unable to process telegram update")
		}
	}()

	return operations.NewReceiveTelegramUpdatesOK()
}

func (p *Polls) WaitTelegram() {
	p.wg.Wait()
}

func (p *Polls) processUpdate(upd tgbotapi.Update) {
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
