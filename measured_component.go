// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
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
