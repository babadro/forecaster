// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	"github.com/babadro/forecaster/internal/models/swagger"
)

// NewCreateSeriesParams creates a new CreateSeriesParams object
//
// There are no default values defined in the spec.
func NewCreateSeriesParams() CreateSeriesParams {

	return CreateSeriesParams{}
}

// CreateSeriesParams contains all the bound params for the create series operation
// typically these are obtained from a http.Request
//
// swagger:parameters CreateSeries
type CreateSeriesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Series object to be created
	  Required: true
	  In: body
	*/
	Series *swagger.CreateSeries
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateSeriesParams() beforehand.
func (o *CreateSeriesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body swagger.CreateSeries
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("series", "body", ""))
			} else {
				res = append(res, errors.NewParseError("series", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(r.Context())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Series = &body
			}
		}
	} else {
		res = append(res, errors.Required("series", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
