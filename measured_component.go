// Copyright 2024 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"

	"github.com/veraison/swid"
)

type MeasuredComponent struct {
	Id             ComponentID     `cbor:"1,keyasint" json:"id"`
	Measurement    *swid.HashEntry `cbor:"2,keyasint,omitempty" json:"measurement,omitempty"`
	Signers        *[]B64Url       `cbor:"3,keyasint,omitempty" json:"signers,omitempty"`
	Flags          *B64Url         `cbor:"4,keyasint,omitempty" json:"flags,omitempty"`
	RawMeasurement *B64Url         `cbor:"5,keyasint,omitempty" json:"raw-measurement,omitempty"`
}

type ComponentID struct {
	_       struct{} `cbor:",toarray"`
	Name    string   `cbor:"0,keyasint"`
	Version *Version `cbor:"1,keyasint,omitempty"`
}

func (c *ComponentID) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("expected JSON array: %w", err)
	}

	if len(tmp) < 1 || 2 < len(tmp) {
		return fmt.Errorf("not component-id value: %#v", tmp)
	}
	if err := json.Unmarshal(tmp[0], &c.Name); err != nil {
		return fmt.Errorf("invalid name: %w", err)
	}

	if len(tmp) == 2 {
		if err := json.Unmarshal(tmp[1], &c.Version); err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}
	}

	return nil
}

func (c ComponentID) MarshalJSON() ([]byte, error) {
	if c.Version == nil {
		return json.Marshal([]string{c.Name})
	}
	return json.Marshal([2]interface{}{c.Name, c.Version})
}

// Based on https://github.com/veraison/swid
// Digest does not support stringify() generating "sha-256;ABC..."
// and encodes/decodes the digest value with base64-url for JSON.
type Digest struct {
	_ struct{} `cbor:",toarray"`

	// The number used as a value for hash-alg-id is an integer-based
	// hash algorithm identifier who's value MUST refer to an ID in the
	// IANA "Named Information Hash Algorithm Registry" [IANA.named-information]
	// with a Status of "current" (at the time the generator software was built
	// or later); other hash algorithms MUST NOT be used. If the hash-alg-id is
	// not known, then the integer value "0" MUST be used. This allows for
	// conversion from ISO SWID tags [SWID], which do not allow an algorithm to
	// be identified for this field.
	Alg int `cbor:"0,keyasint"`

	// The digest value will be base64-url encoded for JSON.
	Val B64Url `cbor:"1,keyasint"`
}

func (d *Digest) UnmarshalJSON(data []byte) error {
	var tmp [2]json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("expected JSON array of length 2: %w", err)
	}

	if err := json.Unmarshal(tmp[0], &d.Alg); err != nil {
		return fmt.Errorf("invalid alg-id: %w", err)
	}

	if err := json.Unmarshal(tmp[1], &d.Val); err != nil {
		return fmt.Errorf("invalid base64url hash value: %w", err)
	}

	return nil
}

func (d Digest) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{d.Alg, d.Val})
}
