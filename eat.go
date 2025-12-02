// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Eat is the internal representation of a EAT token
type Eat struct {
	Nonce             *Nonce                     `cbor:"10,keyasint,omitempty" json:"eat_nonce,omitempty"`
	UEID              *UEID                      `cbor:"256,keyasint,omitempty" json:"ueid,omitempty"`
	SUEIDs            *map[string]UEID           `cbor:"257,keyasint,omitempty" json:"sueids,omitempty"`
	OemID             *OEMID                     `cbor:"258,keyasint,omitempty" json:"oemid,omitempty"`
	HardwareModel     *B64Url                    `cbor:"259,keyasint,omitempty" json:"hwmodel,omitempty"`
	HardwareVersion   *Version                   `cbor:"260,keyasint,omitempty" json:"hwversion,omitempty"`
	Uptime            *uint                      `cbor:"261,keyasint,omitempty" json:"uptime,omitempty"`
	OemBoot           *bool                      `cbor:"262,keyasint,omitempty" json:"oemboot,omitempty"`
	DebugStatus       *Debug                     `cbor:"263,keyasint,omitempty" json:"dbgstat,omitempty"`
	Location          *Location                  `cbor:"264,keyasint,omitempty" json:"location,omitempty"`
	Profile           *Profile                   `cbor:"265,keyasint,omitempty" json:"eat-profile,omitempty"`
	Submods           *Submods                   `cbor:"266,keyasint,omitempty" json:"submods,omitempty"`
	BootCount         *uint                      `cbor:"267,keyasint,omitempty" json:"bootcount,omitempty"`
	BootSeed          *B64Url                    `cbor:"268,keyasint,omitempty" json:"bootseed,omitempty"`
	DLOAs             *Dloa                      `cbor:"269,keyasint,omitempty" json:"dloas,omitempty"`
	SoftwareName      *StringOrURI               `cbor:"270,keyasint,omitempty" json:"swname,omitempty"`
	SoftwareVersion   *Version                   `cbor:"271,keyasint,omitempty" json:"swversion,omitempty"`
	Manifests         *[]Manifest                `cbor:"272,keyasint,omitempty" json:"manifests,omitempty"`
	Measurements      *[]Measurement             `cbor:"273,keyasint,omitempty" json:"measurements,omitempty"`
	MeasrementResults *[]MeasurementResultsGroup `cbor:"274,keyasint,omitempty" json:"measres,omitempty"`
	IntendedUse       *IntendedUse               `cbor:"275,keyasint,omitempty" json:"intuse,omitempty"`
	CWTClaims
}

// FromCBOR deserializes the supplied CBOR encoded EAT into the receiver Eat
func (e *Eat) FromCBOR(data []byte) error {
	return dm.Unmarshal(data, e)
}

// ToCBOR serializes the receiver Eat into CBOR encoded EAT
//
//nolint:gocritic
func (e Eat) ToCBOR() ([]byte, error) {
	return em.Marshal(e)
}

// FromJSON deserializes the supplied JSON encoded EAT into the receiver Eat
func (e *Eat) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ToJSON serializes the receiver Eat into JSON encoded EAT
//
//nolint:gocritic
func (e Eat) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// B64Url is base64url (§5 of RFC4648) without padding.
// bstr MUST be base64url encoded as per EAT §7.2.2 "JSON Interoperability".
type B64Url []byte

func (o B64Url) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		base64.RawURLEncoding.EncodeToString(o),
	)
}

func (b *B64Url) UnmarshalJSON(data []byte) error {
	// get string body
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}

	// while base64.RawURLEncoding.DecodeString("") returns
	// no err, we need to return err because the CDDL definition is here,
	// base64-url-text = tstr .regexp "[A-Za-z0-9_-]+"
	if encoded == "" {
		return fmt.Errorf("base64url must be a non-empty string")
	}

	// decode base64url-encoded string
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("base64url decode error: %w", err)
	}

	*b = decoded
	return nil
}
