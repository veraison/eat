// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
)

// Eat is the internal representation of a EAT token
type Eat struct {
	Nonce           *Nonce    `cbor:"10,keyasint,omitempty" json:"eat_nonce,omitempty"`
	UEID            *UEID     `cbor:"256,keyasint,omitempty" json:"ueid,omitempty"`
	OemID           *[]byte   `cbor:"258,keyasint,omitempty" json:"oemid,omitempty"`
	HardwareModel   *[]byte   `cbor:"259,keyasint,omitempty" json:"hwmodel,omitempty"`
	HardwareVersion *Version  `cbor:"260,keyasint,omitempty" json:"hwversion,omitempty"`
	Uptime          *uint     `cbor:"261,keyasint,omitempty" json:"uptime,omitempty"`
	OemBoot         *bool     `cbor:"262,keyasint,omitempty" json:"oemboot,omitempty"`
	DebugStatus     *Debug    `cbor:"263,keyasint,omitempty" json:"dbgstat,omitempty"`
	Location        *Location `cbor:"264,keyasint,omitempty" json:"location,omitempty"`
	Profile         *Profile  `cbor:"265,keyasint,omitempty" json:"eat-profile,omitempty"`
	Submods         *Submods  `cbor:"266,keyasint,omitempty" json:"submods,omitempty"`
	BootCount       *uint     `cbor:"267,keyasint,omitempty" json:"bootcount,omitempty"`
	BootSeed        *[]byte   `cbor:"268,keyasint,omitempty" json:"bootseed,omitempty"`
	// TODO: DLOAs
	SoftwareName    *StringOrURI   `cbor:"270,keyasint,omitempty" json:"swname,omitempty"`
	SoftwareVersion *Version       `cbor:"271,keyasint,omitempty" json:"swversion,omitempty"`
	Manifests       *[]Manifest    `cbor:"272,keyasint,omitempty" json:"manifests,omitempty"`
	Measurements    *[]Measurement `cbor:"273,keyasint,omitempty" json:"measurements,omitempty"`
	// TODO: MeasrementResults
	// TODO: IntendedUse
	CWTClaims
}

// FromCBOR deserializes the supplied CBOR encoded EAT into the receiver Eat
func (e *Eat) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, e)
}

//nolint:gocritic // ToCBOR serializes the receiver Eat into CBOR encoded EAT
func (e Eat) ToCBOR() ([]byte, error) {
	return em.Marshal(e)
}

// FromJSON deserializes the supplied JSON encoded EAT into the receiver Eat
func (e *Eat) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

//nolint:gocritic // ToJSON serializes the receiver Eat into JSON encoded EAT
func (e Eat) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
