// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// CreatePollCreatedCode is the HTTP code returned for type CreatePollCreated
const CreatePollCreatedCode int = 201

/*
CreatePollCreated Poll created

swagger:response createPollCreated
*/
type CreatePollCreated struct {
}

// NewCreatePollCreated creates CreatePollCreated with default headers values
func NewCreatePollCreated() *CreatePollCreated {

	return &CreatePollCreated{}
}

// WriteResponse to the client
func (o *CreatePollCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(201)
}

// CreatePollBadRequestCode is the HTTP code returned for type CreatePollBadRequest
const CreatePollBadRequestCode int = 400

/*
CreatePollBadRequest Invalid input

swagger:response createPollBadRequest
*/
type CreatePollBadRequest struct {
}

// NewCreatePollBadRequest creates CreatePollBadRequest with default headers values
func NewCreatePollBadRequest() *CreatePollBadRequest {

	return &CreatePollBadRequest{}
}

// WriteResponse to the client
func (o *CreatePollBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}
