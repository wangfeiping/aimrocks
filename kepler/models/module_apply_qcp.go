// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// ModuleApplyQcp module apply qcp
// swagger:model module.ApplyQcp
type ModuleApplyQcp struct {

	// create time
	CreateTime string `json:"createTime,omitempty"`

	// email
	Email string `json:"email,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// info
	Info string `json:"info,omitempty"`

	// note
	Note string `json:"note,omitempty"`

	// phone
	Phone string `json:"phone,omitempty"`

	// qcp chain Id
	QcpChainID string `json:"qcpChainId,omitempty"`

	// qcp pub
	QcpPub string `json:"qcpPub,omitempty"`

	// qos chain Id
	QosChainID string `json:"qosChainId,omitempty"`

	// status
	Status int64 `json:"status,omitempty"`

	// update time
	UpdateTime string `json:"updateTime,omitempty"`
}

// Validate validates this module apply qcp
func (m *ModuleApplyQcp) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ModuleApplyQcp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ModuleApplyQcp) UnmarshalBinary(b []byte) error {
	var res ModuleApplyQcp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
