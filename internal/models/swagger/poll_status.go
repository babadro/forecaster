// Code generated by go-swagger; DO NOT EDIT.

package swagger

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// PollStatus poll status
//
// swagger:model PollStatus
type PollStatus int32

// for schema
var pollStatusEnum []interface{}

func init() {
	var res []PollStatus
	if err := json.Unmarshal([]byte(`[0,1,2,3]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		pollStatusEnum = append(pollStatusEnum, v)
	}
}

func (m PollStatus) validatePollStatusEnum(path, location string, value PollStatus) error {
	if err := validate.EnumCase(path, location, value, pollStatusEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this poll status
func (m PollStatus) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validatePollStatusEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this poll status based on context it is used
func (m PollStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
