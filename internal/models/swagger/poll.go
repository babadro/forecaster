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

// Poll poll
//
// swagger:model Poll
type Poll struct {

	// created at
	// Format: date-time
	CreatedAt strfmt.DateTime `json:"CreatedAt,omitempty"`

	// description
	Description string `json:"Description,omitempty"`

	// finish
	// Format: date-time
	Finish strfmt.DateTime `json:"Finish,omitempty"`

	// ID
	ID int32 `json:"ID,omitempty"`

	// popularity
	Popularity int32 `json:"Popularity,omitempty"`

	// series ID
	SeriesID int32 `json:"SeriesID,omitempty"`

	// start
	// Format: date-time
	Start strfmt.DateTime `json:"Start,omitempty"`

	// status
	Status PollStatus `json:"Status,omitempty"`

	// telegram user ID
	TelegramUserID int64 `json:"TelegramUserID,omitempty"`

	// title
	Title string `json:"Title,omitempty"`

	// updated at
	// Format: date-time
	UpdatedAt strfmt.DateTime `json:"UpdatedAt,omitempty"`
}

// Validate validates this poll
func (m *Poll) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFinish(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpdatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Poll) validateCreatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("CreatedAt", "body", "date-time", m.CreatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Poll) validateFinish(formats strfmt.Registry) error {
	if swag.IsZero(m.Finish) { // not required
		return nil
	}

	if err := validate.FormatOf("Finish", "body", "date-time", m.Finish.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Poll) validateStart(formats strfmt.Registry) error {
	if swag.IsZero(m.Start) { // not required
		return nil
	}

	if err := validate.FormatOf("Start", "body", "date-time", m.Start.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *Poll) validateStatus(formats strfmt.Registry) error {
	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("Status")
		}
		return err
	}

	return nil
}

func (m *Poll) validateUpdatedAt(formats strfmt.Registry) error {
	if swag.IsZero(m.UpdatedAt) { // not required
		return nil
	}

	if err := validate.FormatOf("UpdatedAt", "body", "date-time", m.UpdatedAt.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this poll based on the context it is used
func (m *Poll) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateStatus(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Poll) contextValidateStatus(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	if err := m.Status.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("Status")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("Status")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Poll) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Poll) UnmarshalBinary(b []byte) error {
	var res Poll
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
