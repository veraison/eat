// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "encoding/json"

// Eat is the internal representation of a EAT token
type Eat struct {
	Nonce         *Nonces        `cbor:"10,keyasint,omitempty" json:"nonce,omitempty"`
	UEID          *UEID          `cbor:"11,keyasint,omitempty" json:"ueid,omitempty"`
	Origination   *StringOrURI   `cbor:"12,keyasint,omitempty" json:"origination,omitempty"`
	OemID         *[]byte        `cbor:"13,keyasint,omitempty" json:"oemid,omitempty"`
	SecurityLevel *SecurityLevel `cbor:"14,keyasint,omitempty" json:"security-level,omitempty"`
	SecureBoot    *bool          `cbor:"15,keyasint,omitempty" json:"secure-boot,omitempty"`
	Debug         *Debug         `cbor:"16,keyasint,omitempty" json:"debug-disable,omitempty"`
	Location      *Location      `cbor:"17,keyasint,omitempty" json:"location,omitempty"`
	Uptime        *uint          `cbor:"19,keyasint,omitempty" json:"uptime,omitempty"`
	Submods       *Submods       `cbor:"20,keyasint,omitempty" json:"submods,omitempty"`

	CWTClaims
}

// FromCBOR deserializes the supplied CBOR encoded EAT into the receiver Eat
func (e *Eat) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, e)
}

// ToCBOR serializes the receiver Eat into CBOR encoded EAT
func (e Eat) ToCBOR() ([]byte, error) {
	return em.Marshal(e)
}

// FromJSON deserializes the supplied JSON encoded EAT into the receiver Eat
func (e *Eat) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ToJSON serializes the receiver Eat into CBOR encoded EAT
func (e Eat) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
