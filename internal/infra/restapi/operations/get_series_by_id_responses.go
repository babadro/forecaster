// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/babadro/forecaster/internal/models/swagger"
)

// GetSeriesByIDOKCode is the HTTP code returned for type GetSeriesByIDOK
const GetSeriesByIDOKCode int = 200

/*
GetSeriesByIDOK Series found successfully

swagger:response getSeriesByIdOK
*/
type GetSeriesByIDOK struct {

	/*
	  In: Body
	*/
	Payload *swagger.Series `json:"body,omitempty"`
}

// NewGetSeriesByIDOK creates GetSeriesByIDOK with default headers values
func NewGetSeriesByIDOK() *GetSeriesByIDOK {

	return &GetSeriesByIDOK{}
}

// WithPayload adds the payload to the get series by Id o k response
func (o *GetSeriesByIDOK) WithPayload(payload *swagger.Series) *GetSeriesByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get series by Id o k response
func (o *GetSeriesByIDOK) SetPayload(payload *swagger.Series) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSeriesByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetSeriesByIDNotFoundCode is the HTTP code returned for type GetSeriesByIDNotFound
const GetSeriesByIDNotFoundCode int = 404

/*
GetSeriesByIDNotFound Series not found

swagger:response getSeriesByIdNotFound
*/
type GetSeriesByIDNotFound struct {
}

// NewGetSeriesByIDNotFound creates GetSeriesByIDNotFound with default headers values
func NewGetSeriesByIDNotFound() *GetSeriesByIDNotFound {

	return &GetSeriesByIDNotFound{}
}

// WriteResponse to the client
func (o *GetSeriesByIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

// GetSeriesByIDInternalServerErrorCode is the HTTP code returned for type GetSeriesByIDInternalServerError
const GetSeriesByIDInternalServerErrorCode int = 500

/*
GetSeriesByIDInternalServerError Internal server error

swagger:response getSeriesByIdInternalServerError
*/
type GetSeriesByIDInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetSeriesByIDInternalServerError creates GetSeriesByIDInternalServerError with default headers values
func NewGetSeriesByIDInternalServerError() *GetSeriesByIDInternalServerError {

	return &GetSeriesByIDInternalServerError{}
}

// WithPayload adds the payload to the get series by Id internal server error response
func (o *GetSeriesByIDInternalServerError) WithPayload(payload *swagger.Error) *GetSeriesByIDInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get series by Id internal server error response
func (o *GetSeriesByIDInternalServerError) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSeriesByIDInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*
GetSeriesByIDDefault error

swagger:response getSeriesByIdDefault
*/
type GetSeriesByIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *swagger.Error `json:"body,omitempty"`
}

// NewGetSeriesByIDDefault creates GetSeriesByIDDefault with default headers values
func NewGetSeriesByIDDefault(code int) *GetSeriesByIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetSeriesByIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get series by ID default response
func (o *GetSeriesByIDDefault) WithStatusCode(code int) *GetSeriesByIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get series by ID default response
func (o *GetSeriesByIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get series by ID default response
func (o *GetSeriesByIDDefault) WithPayload(payload *swagger.Error) *GetSeriesByIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get series by ID default response
func (o *GetSeriesByIDDefault) SetPayload(payload *swagger.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetSeriesByIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
