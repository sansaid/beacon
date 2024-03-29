// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ServerBaseResponse server base response
//
// swagger:model server.BaseResponse
type ServerBaseResponse struct {

	// error
	Error string `json:"error,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this server base response
func (m *ServerBaseResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this server base response based on context it is used
func (m *ServerBaseResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ServerBaseResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ServerBaseResponse) UnmarshalBinary(b []byte) error {
	var res ServerBaseResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
