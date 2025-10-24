// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"

	"github.com/veraison/swid"
)

type MeasuredComponent struct {
	Id             ComponentID     `cbor:"1,keyasint" json:"id"`
	Measurement    *swid.HashEntry `cbor:"2,keyasint,omitempty" json:"measurement,omitempty"`
	Signers        *[][]byte       `cbor:"3,keyasint,omitempty" json:"signers,omitempty"`
	Flags          *[]byte         `cbor:"4,keyasint,omitempty" json:"flags,omitempty"`
	RawMeasurement *[]byte         `cbor:"5,keyasint,omitempty" json:"raw-measurement,omitempty"`
}

type ComponentID struct {
	_       struct{} `cbor:",toarray"`
	Name    string   `cbor:"0,keyasint"`
	Version *Version `cbor:"1,keyasint,omitempty"`
}

// FromCBOR deserializes the supplied CBOR encoded MeasuredComponent into the receiver MeasuredComponent
func (mc *MeasuredComponent) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, mc)
}

// ToCBOR serializes the receiver MeasuredComponent into CBOR encoded MeasuredComponent
func (mc MeasuredComponent) ToCBOR() ([]byte, error) {
	return em.Marshal(mc)
}

// FromJSON deserializes the supplied JSON encoded MeasuredComponent into the receiver MeasuredComponent
func (mc *MeasuredComponent) FromJSON(data []byte) error {
	return json.Unmarshal(data, mc)
}

// ToJSON serializes the receiver MeasuredComponent into JSON encoded MeasuredComponent
func (mc MeasuredComponent) ToJSON() ([]byte, error) {
	return json.Marshal(mc)
}
