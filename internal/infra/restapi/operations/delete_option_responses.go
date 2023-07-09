// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/babadro/forecaster/internal/models/swagger"
)

// DeleteOptionNoContentCode is the HTTP code returned for type DeleteOptionNoContent
const DeleteOptionNoContentCode int = 204

/*DeleteOptionNoContent Option deleted successfully

swagger:response deleteOptionNoContent
*/
type DeleteOptionNoContent struct {
}

// NewDeleteOptionNoContent creates DeleteOptionNoContent with default headers values
func NewDeleteOptionNoContent() *DeleteOptionNoContent {

	return &DeleteOptionNoContent{}
}

// WriteResponse to the client
func (o *DeleteOptionNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteOptionNotFoundCode is the HTTP code returned for type DeleteOptionNotFound
const DeleteOptionNotFoundCode int = 404

/*DeleteOptionNotFound Option not found

swagger:response deleteOptionNotFound
*/
type DeleteOptionNotFound struct {
}

// NewDeleteOptionNotFound creates DeleteOptionNotFound with default headers values
func NewDeleteOptionNotFound() *DeleteOptionNotFound {

	return &DeleteOptionNotFound{}
}

// WriteResponse to the client
func (o *DeleteOptionNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// DeleteOptionInternalServerErrorCode is the HTTP code returned for type DeleteOptionInternalServerError
const DeleteOptionInternalServerErrorCode int = 500

/*DeleteOptionInternalServerError Internal server error

swagger:response deleteOptionInternalServerError
*/
type DeleteOptionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewDeleteOptionInternalServerError creates DeleteOptionInternalServerError with default headers values
func NewDeleteOptionInternalServerError() *DeleteOptionInternalServerError {

	return &DeleteOptionInternalServerError{}
}

// WithPayload adds the payload to the delete option internal server error response
func (o *DeleteOptionInternalServerError) WithPayload(payload *swagger.Error) *DeleteOptionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete option internal server error response
func (o *DeleteOptionInternalServerError) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteOptionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteOptionDefault error

swagger:response deleteOptionDefault
*/
type DeleteOptionDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewDeleteOptionDefault creates DeleteOptionDefault with default headers values
func NewDeleteOptionDefault(code int) *DeleteOptionDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteOptionDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete option default response
func (o *DeleteOptionDefault) WithStatusCode(code int) *DeleteOptionDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete option default response
func (o *DeleteOptionDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete option default response
func (o *DeleteOptionDefault) WithPayload(payload *swagger.Error) *DeleteOptionDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete option default response
func (o *DeleteOptionDefault) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteOptionDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
