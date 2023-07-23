// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	bot "github.com/babadro/forecaster/internal/core/forecaster"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/babadro/forecaster/internal/infra/restapi/handlers"
	"github.com/babadro/forecaster/internal/infra/restapi/operations"
	"github.com/caarlos0/env"
	oerrors "github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate swagger generate server --target ../../../../forecaster --name PollAPI --spec ../../../swagger.yaml --model-package internal/models/swagger --server-package internal/infra/restapi --principal interface{}

type envVars struct {
	TelegramToken    string `env:"TELEGRAM_TOKEN,required"`
	DBConn           string `env:"DB_CONN,required"`
	NgrokAgentAddr   string `env:"NGROK_AGENT_ADDR,required"`
	StartTelegramBot bool   `env:"START_TELEGRAM_BOT" envDefault:"true"`
}

func configureFlags(_ *operations.PollAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PollAPIAPI) http.Handler {
	var envs envVars
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	var telegramAPI *handlers.Tg
	if envs.StartTelegramBot {
		publicUrl, err := getNgrokURL(envs.NgrokAgentAddr)
		if err != nil {
			log.Fatalf("Unable to get ngrok url: %v\n", err)
		}

		tgBot, err := initBot(publicUrl+"/telegram-updates", envs.TelegramToken)
		if err != nil {
			log.Fatalf("Unable to init bot: %v\n", err)
		}

		telegramAPI = handlers.NewTelegram(tgBot)
	}

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}

	forecastDB := postgres.NewForecasterDB(dbPool)

	svc := bot.NewService(forecastDB)
	pollsAPI := handlers.NewPolls(svc)

	// configure the api here
	api.ServeError = oerrors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetSeriesByIDHandler = operations.GetSeriesByIDHandlerFunc(pollsAPI.GetSeriesByID)
	api.GetPollByIDHandler = operations.GetPollByIDHandlerFunc(pollsAPI.GetPollByID)

	api.CreateSeriesHandler = operations.CreateSeriesHandlerFunc(pollsAPI.CreateSeries)
	api.CreatePollHandler = operations.CreatePollHandlerFunc(pollsAPI.CreatePoll)
	api.CreateOptionHandler = operations.CreateOptionHandlerFunc(pollsAPI.CreateOption)

	api.UpdateSeriesHandler = operations.UpdateSeriesHandlerFunc(pollsAPI.UpdateSeries)
	api.UpdatePollHandler = operations.UpdatePollHandlerFunc(pollsAPI.UpdatePoll)
	api.UpdateOptionHandler = operations.UpdateOptionHandlerFunc(pollsAPI.UpdateOption)

	api.DeleteSeriesHandler = operations.DeleteSeriesHandlerFunc(pollsAPI.DeleteSeries)
	api.DeletePollHandler = operations.DeletePollHandlerFunc(pollsAPI.DeletePoll)
	api.DeleteOptionHandler = operations.DeleteOptionHandlerFunc(pollsAPI.DeleteOption)

	api.ReceiveTelegramUpdatesHandler = operations.ReceiveTelegramUpdatesHandlerFunc(telegramAPI.ReceiveUpdates)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		dbPool.Close()
		if envs.StartTelegramBot {
			telegramAPI.Wait()
		}
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(_ *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything,
// this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

type ngrokResp struct {
	Tunnels []struct {
		PublicURL string `json:"public_url"`
	} `json:"tunnels"`
}

func getNgrokURL(agentAddr string) (string, error) {
	resp, err := http.Get(agentAddr + "/api/tunnels")
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error while closing response body: %v", err)
		}
	}(resp.Body)

	var response ngrokResp
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if len(response.Tunnels) > 0 {
		return response.Tunnels[0].PublicURL, nil
	}

	return "", errors.New("no active ngrok tunnels found")
}

func initBot(link, token string) (*tgbotapi.BotAPI, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	botAPI.Debug = true

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	wh, err := tgbotapi.NewWebhook(link)
	if err != nil {
		return nil, err
	}

	_, err = botAPI.Request(wh)
	if err != nil {
		return nil, err
	}

	info, err := botAPI.GetWebhookInfo()
	if err != nil {
		return nil, err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return botAPI, nil
}
