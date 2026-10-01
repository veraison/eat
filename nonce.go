// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// nonce-type = bstr .size (8..64)
const (
	MinNonceSize = 8
	MaxNonceSize = 64
)

func isValidNonce(v []byte) error {
	nonceSize := len(v)
	if nonceSize < MinNonceSize || nonceSize > MaxNonceSize {
		return fmt.Errorf(
			"a nonce must be between %d and %d bytes long; found %d",
			MinNonceSize, MaxNonceSize, nonceSize,
		)
	}
	return nil
}

type nonce struct {
	value BinaryData
}

// newNonce returns a nonce initialized with the supplied byte slice or an error
// if the supplied buffer is either too big (more than 64 bytes) or too small
// (less than 8 bytes)
func newNonce(v []byte) (*nonce, error) {
	if err := isValidNonce(v); err != nil {
		return nil, err
	}

	return &nonce{v}, nil
}

// get returns the nonce value
func (o *nonce) get() []byte {
	return o.value
}

// validate checks that the nonce is between 8 and 64 bytes, as is required by
// the EAT spec
func (o *nonce) validate() error {
	return isValidNonce(o.value)
}

// MarshalCBOR encodes the nonce as a CBOR byte string
func (o nonce) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

// UnmarshalCBOR decodes a CBOR byte string and uses it as the nonce value
func (o *nonce) UnmarshalCBOR(data []byte) error {
	var value BinaryData

	if err := dm.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("CBOR decoding failed for nonce: %w", err)
	}

	o.value = value

	return nil
}

// MarshalJSON encodes the receiver (non-array) nonce as a JSON string
func (o *nonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

// UnmarshalJSON decodes the supplied JSON data to a (non-array) nonce
func (o *nonce) UnmarshalJSON(data []byte) error {
	var value BinaryData

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	o.value = value

	return nil
}

// A nonce-claim may be single Nonce or an array of two or more.
//
//	nonce-claim = (
//	    nonce => nonce-type / [ 2* nonce-type ]
//	)
type Nonce struct {
	values []*nonce
}

// MustNewNonceFromBytes creates a new Nonce containing a single entry created
// from the provided bytes. If the provided bytes do not constitute a valid
// nonce, this function panics.
func MustNewNonceFromBytes(values ...[]byte) *Nonce {
	ret, err := NewNonceFromBytes(values...)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewNonceFromBytes creates a new Nonce containing a single entry created from
// the provided bytes. If the provided bytes do not constitute a valid nonce,
// an error is returned.
func NewNonceFromBytes(values ...[]byte) (*Nonce, error) {
	if len(values) == 0 {
		return nil, errors.New("at least one nonce value must be specified")
	}

	var ret Nonce
	for i, value := range values {
		if err := ret.Add(value); err != nil {
			return nil, fmt.Errorf("value[%d]: %w", i, err)
		}
	}

	return &ret, nil
}

// Add the supplied nonce, provided as a byte array, to the Nonce receiver.
func (o *Nonce) Add(v []byte) error {
	n, err := newNonce(v)
	if err != nil {
		return err
	}

	o.values = append(o.values, n)

	return nil
}

// Len returns the number of nonce values carried in the Nonce receiver.
func (o Nonce) Len() int {
	return len(o.values)
}

// Get returns the nonce found at the supplied index (counting from 0) or nil
// if the index is out of bounds.
func (o Nonce) Get(index int) []byte {
	if index < 0 || index >= o.Len() {
		return nil
	}

	return o.values[index].get()
}

// AddHex provides the same functionality as Add except it takes the nonce value
// as a hex-encoded string.
func (o *Nonce) AddHex(text string) error {
	value, err := hex.DecodeString(text)
	if err != nil {
		return fmt.Errorf("decoding nonce failed: %w", err)
	}

	return o.Add(value)
}

// MarshalCBOR provides a suitable CBOR encoding for the receiver Nonce. In
// case there is only one nonce, the encoded produces a single bstr. If there
// are multiple, the encoder produces an array of bstr, one for each nonce.
func (o Nonce) MarshalCBOR() ([]byte, error) {
	if err := validateNonceValues(o.values); err != nil {
		return nil, err
	}

	if len(o.values) == 1 {
		return em.Marshal(o.values[0])
	}

	return em.Marshal(o.values)
}

// UnmarshalCBOR decodes a EAT nonce. This may be a single byte string
// between 8 and 64 bytes long, or an array of two or more such strings.
func (o *Nonce) UnmarshalCBOR(data []byte) error {
	var values []*nonce

	if isCBORArray(data) {
		if err := dm.Unmarshal(data, &values); err != nil {
			return err
		}
	} else {
		var n nonce

		if err := dm.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("CBOR decoding failed for nonce: %w", err)
		}

		values = []*nonce{&n}
	}

	if err := validateNonceValues(values); err != nil {
		return err
	}

	o.values = values

	return nil
}

// MarshalJSON encodes the receiver Nonce as either a JSON string containing
// the base64 encoding of the binary nonce (if the array comprises only one
// element) or as an array of base64-encoded JSON strings.
//
// NOTE: While RFC 9711 (EAT) does not restrict the nonce format to base64-encoded
// JSON strings, VERAISON/eat imposes a narrower interpretation for
// JSON <-> CBOR conversion.
// See discussion: https://github.com/ietf-rats-wg/eat/pull/421
func (o Nonce) MarshalJSON() ([]byte, error) {
	if err := validateNonceValues(o.values); err != nil {
		return nil, err
	}

	if len(o.values) == 1 {
		return json.Marshal(o.values[0])
	}

	return json.Marshal(o.values)
}

// UnmarshalJSON decodes a EAT nonce in JSON format.
func (o *Nonce) UnmarshalJSON(data []byte) error {
	var values []*nonce

	if isJSONArray(data) {
		if err := json.Unmarshal(data, &values); err != nil {
			return err
		}
	} else {
		var n nonce

		if err := json.Unmarshal(data, &n); err != nil {
			return fmt.Errorf("JSON decoding failed for nonce: %w", err)
		}

		values = []*nonce{&n}
	}

	if err := validateNonceValues(values); err != nil {
		return err
	}

	o.values = values

	return nil
}

// validateNonceValues checks that all nonce values (of which there must be at least one)
// stored in the Nonce receiver are valid according to the EAT syntax.
func validateNonceValues(values []*nonce) error {
	if len(values) == 0 {
		return fmt.Errorf("empty nonce")
	}

	for i, n := range values {
		if n == nil {
			return fmt.Errorf("found nil nonce at index %d", i)
		}

		if err := n.validate(); err != nil {
			return fmt.Errorf("found invalid nonce at index %d: %w", i, err)
		}
	}

	return nil
}
