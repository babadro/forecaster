// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/babadro/forecaster/internal/core/forecaster/polls"
	"github.com/babadro/forecaster/internal/core/forecaster/telegram"
	pollshandlers "github.com/babadro/forecaster/internal/infra/restapi/handlers/polls"
	telegramhandlers "github.com/babadro/forecaster/internal/infra/restapi/handlers/telegram"
	"github.com/babadro/forecaster/internal/infra/restapi/middlewares"
	"github.com/go-openapi/runtime/middleware"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/babadro/forecaster/internal/infra/postgres"
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
	l := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	var envs envVars
	if err := env.Parse(&envs); err != nil {
		l.Fatal().Msgf("Unable to parse env vars: %v\n", err)
	}

	ctx := context.Background()

	var tgBot *tgbotapi.BotAPI

	if envs.StartTelegramBot {
		publicURL, err := getNgrokURL(ctx, envs.NgrokAgentAddr)
		if err != nil {
			l.Fatal().Msgf("Unable to get ngrok url: %v\n", err)
		}

		tgBot, err = initBot(publicURL+"/telegram-updates", envs.TelegramToken)
		if err != nil {
			l.Fatal().Msgf("Unable to init bot: %v\n", err)
		}
	}

	dbPool, err := pgxpool.Connect(ctx, envs.DBConn)
	if err != nil {
		l.Fatal().Msgf("Unable to connection to database :%v\n", err)
	}

	forecastDB := postgres.NewForecasterDB(dbPool)

	pollsService := polls.NewService(forecastDB)
	pollsAPI := pollshandlers.NewPolls(pollsService)

	var telegramAPI *telegramhandlers.Telegram
	if envs.StartTelegramBot {
		telegramAPI = telegramhandlers.NewTelegram(telegram.NewService(forecastDB, tgBot))
	}

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

	if envs.StartTelegramBot {
		api.ReceiveTelegramUpdatesHandler = operations.ReceiveTelegramUpdatesHandlerFunc(telegramAPI.ReceiveTelegramUpdates)
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		dbPool.Close()

		if envs.StartTelegramBot {
			telegramAPI.Wait()
		}
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares(l)))
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
	_, _, _ = s, scheme, addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(l zerolog.Logger) middleware.Builder {
	return alice.New(
		middlewares.Logging(l),
	).Then
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

func getNgrokURL(ctx context.Context, agentAddr string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, agentAddr+"/api/tunnels", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

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
