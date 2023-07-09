// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewPollAPIAPI creates a new PollAPI instance
func NewPollAPIAPI(spec *loads.Document) *PollAPIAPI {
	return &PollAPIAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		CreateOptionHandler: CreateOptionHandlerFunc(func(params CreateOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation CreateOption has not yet been implemented")
		}),
		CreatePollHandler: CreatePollHandlerFunc(func(params CreatePollParams) middleware.Responder {
			return middleware.NotImplemented("operation CreatePoll has not yet been implemented")
		}),
		DeleteOptionHandler: DeleteOptionHandlerFunc(func(params DeleteOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation DeleteOption has not yet been implemented")
		}),
		DeletePollHandler: DeletePollHandlerFunc(func(params DeletePollParams) middleware.Responder {
			return middleware.NotImplemented("operation DeletePoll has not yet been implemented")
		}),
		GetPollByIDHandler: GetPollByIDHandlerFunc(func(params GetPollByIDParams) middleware.Responder {
			return middleware.NotImplemented("operation GetPollByID has not yet been implemented")
		}),
		ReceiveTelegramUpdatesHandler: ReceiveTelegramUpdatesHandlerFunc(func(params ReceiveTelegramUpdatesParams) middleware.Responder {
			return middleware.NotImplemented("operation ReceiveTelegramUpdates has not yet been implemented")
		}),
		UpdateOptionHandler: UpdateOptionHandlerFunc(func(params UpdateOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation UpdateOption has not yet been implemented")
		}),
		UpdatePollHandler: UpdatePollHandlerFunc(func(params UpdatePollParams) middleware.Responder {
			return middleware.NotImplemented("operation UpdatePoll has not yet been implemented")
		}),
	}
}

/*PollAPIAPI API for managing Polls and Options */
type PollAPIAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// CreateOptionHandler sets the operation handler for the create option operation
	CreateOptionHandler CreateOptionHandler
	// CreatePollHandler sets the operation handler for the create poll operation
	CreatePollHandler CreatePollHandler
	// DeleteOptionHandler sets the operation handler for the delete option operation
	DeleteOptionHandler DeleteOptionHandler
	// DeletePollHandler sets the operation handler for the delete poll operation
	DeletePollHandler DeletePollHandler
	// GetPollByIDHandler sets the operation handler for the get poll by ID operation
	GetPollByIDHandler GetPollByIDHandler
	// ReceiveTelegramUpdatesHandler sets the operation handler for the receive telegram updates operation
	ReceiveTelegramUpdatesHandler ReceiveTelegramUpdatesHandler
	// UpdateOptionHandler sets the operation handler for the update option operation
	UpdateOptionHandler UpdateOptionHandler
	// UpdatePollHandler sets the operation handler for the update poll operation
	UpdatePollHandler UpdatePollHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *PollAPIAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *PollAPIAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *PollAPIAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *PollAPIAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *PollAPIAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *PollAPIAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *PollAPIAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *PollAPIAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *PollAPIAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the PollAPIAPI
func (o *PollAPIAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.CreateOptionHandler == nil {
		unregistered = append(unregistered, "CreateOptionHandler")
	}
	if o.CreatePollHandler == nil {
		unregistered = append(unregistered, "CreatePollHandler")
	}
	if o.DeleteOptionHandler == nil {
		unregistered = append(unregistered, "DeleteOptionHandler")
	}
	if o.DeletePollHandler == nil {
		unregistered = append(unregistered, "DeletePollHandler")
	}
	if o.GetPollByIDHandler == nil {
		unregistered = append(unregistered, "GetPollByIDHandler")
	}
	if o.ReceiveTelegramUpdatesHandler == nil {
		unregistered = append(unregistered, "ReceiveTelegramUpdatesHandler")
	}
	if o.UpdateOptionHandler == nil {
		unregistered = append(unregistered, "UpdateOptionHandler")
	}
	if o.UpdatePollHandler == nil {
		unregistered = append(unregistered, "UpdatePollHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *PollAPIAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *PollAPIAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *PollAPIAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *PollAPIAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *PollAPIAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *PollAPIAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the poll API API
func (o *PollAPIAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *PollAPIAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/options"] = NewCreateOption(o.context, o.CreateOptionHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/polls"] = NewCreatePoll(o.context, o.CreatePollHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/options/{optionId}"] = NewDeleteOption(o.context, o.DeleteOptionHandler)
	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/polls/{pollId}"] = NewDeletePoll(o.context, o.DeletePollHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/polls/{pollId}"] = NewGetPollByID(o.context, o.GetPollByIDHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/telegram-updates"] = NewReceiveTelegramUpdates(o.context, o.ReceiveTelegramUpdatesHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/options/{optionId}"] = NewUpdateOption(o.context, o.UpdateOptionHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/polls/{pollId}"] = NewUpdatePoll(o.context, o.UpdatePollHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *PollAPIAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *PollAPIAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *PollAPIAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *PollAPIAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *PollAPIAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[um][path] = builder(h)
	}
}
