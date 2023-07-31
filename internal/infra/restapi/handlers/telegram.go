package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	"github.com/babadro/forecaster/internal/models/swagger"
	"github.com/go-openapi/runtime/middleware"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/hlog"
)

func (p *Polls) ReceiveTelegramUpdates(params operations.ReceiveTelegramUpdatesParams) middleware.Responder {
	var update tgbotapi.Update
	logger := hlog.FromRequest(params.HTTPRequest)

	bodyBytes, err := io.ReadAll(params.HTTPRequest.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Unable to read body")

		return operations.NewReceiveTelegramUpdatesBadRequest().WithPayload(&swagger.Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Unable to read body: %v", err),
		})
	}

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&update)
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
				Bytes("update", bodyBytes).
				Msg("Unable to process telegram update")
		}
	}()

	return operations.NewReceiveTelegramUpdatesOK()
}

func (p *Polls) WaitTelegram() {
	p.wg.Wait()
}
