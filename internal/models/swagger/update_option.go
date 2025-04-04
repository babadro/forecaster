// Code generated by go-swagger; DO NOT EDIT.

package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UpdateOption update option
//
// swagger:model UpdateOption
type UpdateOption struct {

	// description
	Description *string `json:"Description,omitempty"`

	// is actual outcome
	IsActualOutcome *bool `json:"IsActualOutcome,omitempty"`

	// title
	Title *string `json:"Title,omitempty"`
}

// Validate validates this update option
func (m *UpdateOption) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this update option based on context it is used
func (m *UpdateOption) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *UpdateOption) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *UpdateOption) UnmarshalBinary(b []byte) error {
	var res UpdateOption
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
