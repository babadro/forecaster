// Code generated by go-swagger; DO NOT EDIT.

package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CreatePoll create poll
//
// swagger:model CreatePoll
type CreatePoll struct {

	// description
	Description string `json:"Description,omitempty"`

	// finish
	// Format: date-time
	Finish strfmt.DateTime `json:"Finish,omitempty"`

	// start
	// Format: date-time
	Start strfmt.DateTime `json:"Start,omitempty"`

	// title
	Title string `json:"Title,omitempty"`
}

// Validate validates this create poll
func (m *CreatePoll) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFinish(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreatePoll) validateFinish(formats strfmt.Registry) error {
	if swag.IsZero(m.Finish) { // not required
		return nil
	}

	if err := validate.FormatOf("Finish", "body", "date-time", m.Finish.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CreatePoll) validateStart(formats strfmt.Registry) error {
	if swag.IsZero(m.Start) { // not required
		return nil
	}

	if err := validate.FormatOf("Start", "body", "date-time", m.Start.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this create poll based on context it is used
func (m *CreatePoll) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreatePoll) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatePoll) UnmarshalBinary(b []byte) error {
	var res CreatePoll
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
