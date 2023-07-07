// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeletePollHandlerFunc turns a function with the right signature into a delete poll handler
type DeletePollHandlerFunc func(DeletePollParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeletePollHandlerFunc) Handle(params DeletePollParams) middleware.Responder {
	return fn(params)
}

// DeletePollHandler interface for that can handle valid delete poll params
type DeletePollHandler interface {
	Handle(DeletePollParams) middleware.Responder
}

// NewDeletePoll creates a new http.Handler for the delete poll operation
func NewDeletePoll(ctx *middleware.Context, handler DeletePollHandler) *DeletePoll {
	return &DeletePoll{Context: ctx, Handler: handler}
}

/*
DeletePoll swagger:route DELETE /polls/{pollId} deletePoll

Delete a Poll by its ID
*/
type DeletePoll struct {
	Context *middleware.Context
	Handler DeletePollHandler
}

func (o *DeletePoll) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeletePollParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
