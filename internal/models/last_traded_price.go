// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// LastTradedPrice last traded price
//
// swagger:model LastTradedPrice
type LastTradedPrice struct {

	// amount
	Amount float64 `json:"amount,omitempty"`

	// pair
	Pair string `json:"pair,omitempty"`
}

// Validate validates this last traded price
func (m *LastTradedPrice) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this last traded price based on context it is used
func (m *LastTradedPrice) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LastTradedPrice) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LastTradedPrice) UnmarshalBinary(b []byte) error {
	var res LastTradedPrice
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
