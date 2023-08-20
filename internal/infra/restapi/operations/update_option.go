// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateOptionHandlerFunc turns a function with the right signature into a update option handler
type UpdateOptionHandlerFunc func(UpdateOptionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateOptionHandlerFunc) Handle(params UpdateOptionParams) middleware.Responder {
	return fn(params)
}

// UpdateOptionHandler interface for that can handle valid update option params
type UpdateOptionHandler interface {
	Handle(UpdateOptionParams) middleware.Responder
}

// NewUpdateOption creates a new http.Handler for the update option operation
func NewUpdateOption(ctx *middleware.Context, handler UpdateOptionHandler) *UpdateOption {
	return &UpdateOption{Context: ctx, Handler: handler}
}

/*
UpdateOption swagger:route PUT /options/{pollId}/{optionId} updateOption

Update an existing Option
*/
type UpdateOption struct {
	Context *middleware.Context
	Handler UpdateOptionHandler
}

func (o *UpdateOption) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateOptionParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
