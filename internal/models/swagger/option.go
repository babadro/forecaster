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

// Option option
//
// swagger:model Option
type Option struct {

	// description
	Description string `json:"Description,omitempty"`

	// ID
	ID int32 `json:"ID,omitempty"`

	// poll ID
	PollID int32 `json:"PollID,omitempty"`

	// title
	Title string `json:"Title,omitempty"`

	// updated at
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"UpdatedAt,omitempty"`
}

// Validate validates this option
func (m *Option) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Option) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("UpdatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this option based on context it is used
func (m *Option) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Option) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Option) UnmarshalBinary(b []byte) error {
	var res Option
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
