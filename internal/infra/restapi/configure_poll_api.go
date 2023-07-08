// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	bot "github.com/babadro/forecaster/internal/core/forecaster"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/babadro/forecaster/internal/infra/restapi/handlers"
	"github.com/caarlos0/env"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"

	"github.com/babadro/forecaster/internal/infra/restapi/operations"
)

//go:generate swagger generate server --target ../../../../forecaster --name PollAPI --spec ../../../swagger.yaml --model-package internal/models/swagger --server-package internal/infra/restapi --principal interface{}

type envVars struct {
	HTTPAddr       string `env:"HTTP_ADDR,required"`
	TelegramToken  string `env:"TELEGRAM_TOKEN,required"`
	DBConn         string `env:"DB_CONN,required"`
	NgrokAuthToken string `env:"NGROK_AUTH_TOKEN,required"`
}

func configureFlags(_ *operations.PollAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PollAPIAPI) http.Handler {
	var envs envVars
	if err := env.Parse(&envs); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	dbPool, err := pgxpool.Connect(context.Background(), envs.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}

	forecastDB := postgres.NewForecasterDB(dbPool)

	svc := bot.NewService(forecastDB)
	pollsAPI := handlers.NewPolls(svc)

	tunnel, err := ngrokRun(context.Background(), envs.NgrokAuthToken)
	if err != nil {
		dbPool.Close()
		log.Fatalf("Unable to run ngrok: %v\n", err)
	}

	tgBot, err := initBot(tunnel.URL(), envs.TelegramToken)
	if err != nil {
		dbPool.Close()
		log.Fatalf("Unable to init bot: %v\n", err)
	}

	telegramAPI := handlers.NewTelegram(tgBot)

	// configure the api here
	api.ServeError = errors.ServeError

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

	api.CreatePollHandler = operations.CreatePollHandlerFunc(pollsAPI.CreatePoll)
	api.GetPollByIDHandler = operations.GetPollByIDHandlerFunc(pollsAPI.GetPollByID)
	api.DeletePollHandler = operations.DeletePollHandlerFunc(pollsAPI.DeletePoll)
	api.UpdatePollHandler = operations.UpdatePollHandlerFunc(pollsAPI.UpdatePoll)

	api.CreateOptionHandler = operations.CreateOptionHandlerFunc(pollsAPI.CreateOption)
	api.UpdateOptionHandler = operations.UpdateOptionHandlerFunc(pollsAPI.UpdateOption)
	api.DeleteOptionHandler = operations.DeleteOptionHandlerFunc(pollsAPI.DeleteOption)

	api.ReceiveTelegramUpdatesHandler = operations.ReceiveTelegramUpdatesHandlerFunc(telegramAPI.ReceiveUpdates)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		dbPool.Close()
		telegramAPI.Wait()
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

func ngrokRun(ctx context.Context, token string) (ngrok.Tunnel, error) {
	tun, err := ngrok.Listen(ctx,
		config.HTTPEndpoint(),
		ngrok.WithAuthtoken(token),
	)

	if err != nil {
		return nil, err
	}

	log.Println("tunnel created:", tun.URL())

	return tun, nil
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
