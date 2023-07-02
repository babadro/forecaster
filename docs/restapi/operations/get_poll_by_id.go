// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetPollByIDHandlerFunc turns a function with the right signature into a get poll by Id handler
type GetPollByIDHandlerFunc func(GetPollByIDParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetPollByIDHandlerFunc) Handle(params GetPollByIDParams) middleware.Responder {
	return fn(params)
}

// GetPollByIDHandler interface for that can handle valid get poll by Id params
type GetPollByIDHandler interface {
	Handle(GetPollByIDParams) middleware.Responder
}

// NewGetPollByID creates a new http.Handler for the get poll by Id operation
func NewGetPollByID(ctx *middleware.Context, handler GetPollByIDHandler) *GetPollByID {
	return &GetPollByID{Context: ctx, Handler: handler}
}

/*
GetPollByID swagger:route GET /polls/{pollId} getPollById

Get poll by ID
*/
type GetPollByID struct {
	Context *middleware.Context
	Handler GetPollByIDHandler
}

func (o *GetPollByID) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetPollByIDParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
