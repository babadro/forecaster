// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/babadro/forecaster/internal/models/swagger"
)

// CreateOptionCreatedCode is the HTTP code returned for type CreateOptionCreated
const CreateOptionCreatedCode int = 201

/*
CreateOptionCreated Option created successfully

swagger:response createOptionCreated
*/
type CreateOptionCreated struct {

	/*
	  In: Body
	*/
	Payload *swagger.Option `json:"body,omitempty"`
}

// NewCreateOptionCreated creates CreateOptionCreated with default headers values
func NewCreateOptionCreated() *CreateOptionCreated {

	return &CreateOptionCreated{}
}

// WithPayload adds the payload to the create option created response
func (o *CreateOptionCreated) WithPayload(payload *swagger.Option) *CreateOptionCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create option created response
func (o *CreateOptionCreated) SetPayload(payload *swagger.Option) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateOptionCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateOptionBadRequestCode is the HTTP code returned for type CreateOptionBadRequest
const CreateOptionBadRequestCode int = 400

/*
CreateOptionBadRequest Bad request

swagger:response createOptionBadRequest
*/
type CreateOptionBadRequest struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewCreateOptionBadRequest creates CreateOptionBadRequest with default headers values
func NewCreateOptionBadRequest() *CreateOptionBadRequest {

	return &CreateOptionBadRequest{}
}

// WithPayload adds the payload to the create option bad request response
func (o *CreateOptionBadRequest) WithPayload(payload *swagger.Error) *CreateOptionBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create option bad request response
func (o *CreateOptionBadRequest) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateOptionBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateOptionInternalServerErrorCode is the HTTP code returned for type CreateOptionInternalServerError
const CreateOptionInternalServerErrorCode int = 500

/*
CreateOptionInternalServerError Internal server error

swagger:response createOptionInternalServerError
*/
type CreateOptionInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewCreateOptionInternalServerError creates CreateOptionInternalServerError with default headers values
func NewCreateOptionInternalServerError() *CreateOptionInternalServerError {

	return &CreateOptionInternalServerError{}
}

// WithPayload adds the payload to the create option internal server error response
func (o *CreateOptionInternalServerError) WithPayload(payload *swagger.Error) *CreateOptionInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create option internal server error response
func (o *CreateOptionInternalServerError) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateOptionInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
CreateOptionDefault error

swagger:response createOptionDefault
*/
type CreateOptionDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewCreateOptionDefault creates CreateOptionDefault with default headers values
func NewCreateOptionDefault(code int) *CreateOptionDefault {
	if code <= 0 {
		code = 500
	}

	return &CreateOptionDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the create option default response
func (o *CreateOptionDefault) WithStatusCode(code int) *CreateOptionDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the create option default response
func (o *CreateOptionDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the create option default response
func (o *CreateOptionDefault) WithPayload(payload *swagger.Error) *CreateOptionDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create option default response
func (o *CreateOptionDefault) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateOptionDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
