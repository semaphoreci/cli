// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// SemaphoreSecretsV1betaSecret semaphore secrets v1beta secret
// swagger:model SemaphoreSecretsV1betaSecret
type SemaphoreSecretsV1betaSecret struct {

	// data
	Data *SemaphoreSecretsV1betaSecretData `json:"data,omitempty"`

	// metadata
	Metadata *SemaphoreSecretsV1betaSecretMetadata `json:"metadata,omitempty"`
}

// Validate validates this semaphore secrets v1beta secret
func (m *SemaphoreSecretsV1betaSecret) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateData(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMetadata(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SemaphoreSecretsV1betaSecret) validateData(formats strfmt.Registry) error {

	if swag.IsZero(m.Data) { // not required
		return nil
	}

	if m.Data != nil {
		if err := m.Data.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("data")
			}
			return err
		}
	}

	return nil
}

func (m *SemaphoreSecretsV1betaSecret) validateMetadata(formats strfmt.Registry) error {

	if swag.IsZero(m.Metadata) { // not required
		return nil
	}

	if m.Metadata != nil {
		if err := m.Metadata.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("metadata")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SemaphoreSecretsV1betaSecret) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SemaphoreSecretsV1betaSecret) UnmarshalBinary(b []byte) error {
	var res SemaphoreSecretsV1betaSecret
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
