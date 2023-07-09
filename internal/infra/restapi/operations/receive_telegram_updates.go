// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ReceiveTelegramUpdatesHandlerFunc turns a function with the right signature into a receive telegram updates handler
type ReceiveTelegramUpdatesHandlerFunc func(ReceiveTelegramUpdatesParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ReceiveTelegramUpdatesHandlerFunc) Handle(params ReceiveTelegramUpdatesParams) middleware.Responder {
	return fn(params)
}

// ReceiveTelegramUpdatesHandler interface for that can handle valid receive telegram updates params
type ReceiveTelegramUpdatesHandler interface {
	Handle(ReceiveTelegramUpdatesParams) middleware.Responder
}

// NewReceiveTelegramUpdates creates a new http.Handler for the receive telegram updates operation
func NewReceiveTelegramUpdates(ctx *middleware.Context, handler ReceiveTelegramUpdatesHandler) *ReceiveTelegramUpdates {
	return &ReceiveTelegramUpdates{Context: ctx, Handler: handler}
}

/*ReceiveTelegramUpdates swagger:route POST /telegram-updates receiveTelegramUpdates

Receive updates from Telegram

*/
type ReceiveTelegramUpdates struct {
	Context *middleware.Context
	Handler ReceiveTelegramUpdatesHandler
}

func (o *ReceiveTelegramUpdates) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewReceiveTelegramUpdatesParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
