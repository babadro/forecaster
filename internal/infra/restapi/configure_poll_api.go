// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/babadro/forecaster/internal/infra/restapi/operations"
)

//go:generate swagger generate server --target ../../../../forecaster --name PollAPI --spec ../../../swagger.yaml --model-package internal/models/swagger --server-package internal/infra/restapi --principal interface{}

func configureFlags(_ *operations.PollAPIAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.PollAPIAPI) http.Handler {
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

	if api.CreateOptionHandler == nil {
		api.CreateOptionHandler = operations.CreateOptionHandlerFunc(
			func(params operations.CreateOptionParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.CreateOption has not yet been implemented")
			})
	}

	if api.CreatePollHandler == nil {
		api.CreatePollHandler = operations.CreatePollHandlerFunc(
			func(params operations.CreatePollParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.CreatePoll has not yet been implemented")
			})
	}

	if api.DeleteOptionHandler == nil {
		api.DeleteOptionHandler = operations.DeleteOptionHandlerFunc(
			func(params operations.DeleteOptionParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.DeleteOption has not yet been implemented")
			})
	}

	if api.DeletePollHandler == nil {
		api.DeletePollHandler = operations.DeletePollHandlerFunc(
			func(params operations.DeletePollParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.DeletePoll has not yet been implemented")
			})
	}

	if api.GetPollByIDHandler == nil {
		api.GetPollByIDHandler = operations.GetPollByIDHandlerFunc(
			func(params operations.GetPollByIDParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.GetPollByID has not yet been implemented")
			})
	}

	if api.ReceiveTelegramUpdatesHandler == nil {
		api.ReceiveTelegramUpdatesHandler = operations.ReceiveTelegramUpdatesHandlerFunc(
			func(params operations.ReceiveTelegramUpdatesParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.ReceiveTelegramUpdates has not yet been implemented")
			})
	}

	if api.UpdateOptionHandler == nil {
		api.UpdateOptionHandler = operations.UpdateOptionHandlerFunc(
			func(params operations.UpdateOptionParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.UpdateOption has not yet been implemented")
			})
	}

	if api.UpdatePollHandler == nil {
		api.UpdatePollHandler = operations.UpdatePollHandlerFunc(
			func(params operations.UpdatePollParams) middleware.Responder {
				return middleware.NotImplemented("operation operations.UpdatePoll has not yet been implemented")
			})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

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
