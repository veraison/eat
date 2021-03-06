// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
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
	value []byte
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
func (n nonce) get() []byte {
	return n.value
}

// validate checks that the nonce is between 8 and 64 bytes, as is required by
// the EAT spec
func (n nonce) validate() error {
	return isValidNonce(n.value)
}

// MarshalCBOR encodes the nonce as a CBOR byte string
func (n nonce) MarshalCBOR() ([]byte, error) {
	return em.Marshal(n.value)
}

// UnmarshalCBOR decodes a CBOR byte string and uses it as the nonce value
func (n *nonce) UnmarshalCBOR(data []byte) error {
	var value []byte

	if err := dm.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("CBOR decoding failed for nonce: %w", err)
	}

	n.value = value

	return nil
}

// MarshalJSON encodes the receiver (non-array) nonce as a JSON string
func (n nonce) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.value)
}

// UnmarshalJSON decodes the supplied JSON data to a (non-array) nonce
func (n *nonce) UnmarshalJSON(data []byte) error {
	var v interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("JSON decoding failed for nonce: %w", err)
	}

	switch t := v.(type) {
	case string:
		value, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return err
		}
		n.value = value
		return nil
	default:
		return fmt.Errorf("invalid nonce input %T", t)
	}
}

// A nonce-claim may be single Nonce or an array of two or more.
//
//    nonce-claim = (
//        nonce => nonce-type / [ 2* nonce-type ]
//    )
type Nonce []nonce

// Validate checks that all nonce values (of which there must be at least one)
// stored in the Nonce receiver are valid according to the EAT syntax.
func (ns Nonce) Validate() error {
	if len(ns) == 0 {
		return fmt.Errorf("empty nonce")
	}

	for i, n := range ns {
		if err := n.validate(); err != nil {
			return fmt.Errorf("found invalid nonce at index %d: %w", i, err)
		}
	}

	return nil
}

// Add the supplied nonce, provided as a byte array, to the Nonce receiver.
func (ns *Nonce) Add(v []byte) error {
	n, err := newNonce(v)
	if err != nil {
		return err
	}

	*ns = append(*ns, *n)

	return nil
}

// Len returns the number of nonce values carried in the Nonce receiver.
func (ns Nonce) Len() int {
	return len(ns)
}

// GetI returns the nonce found at the supplied index (counting from 0) or nil
// if the index is out of bounds.
func (ns Nonce) GetI(index int) []byte {
	if index < 0 || index >= ns.Len() {
		return nil
	}

	return ns[index].get()
}

// AddHex provides the same functionality as Add except it takes the nonce value
// as a hex-encoded string.
func (ns *Nonce) AddHex(text string) error {
	value, err := hex.DecodeString(text)
	if err != nil {
		return fmt.Errorf("decoding nonce failed: %w", err)
	}

	return ns.Add(value)
}

// MarshalCBOR provides a suitable CBOR encoding for the receiver Nonce. In
// case there is only one nonce, the encoded produces a single bstr. If there
// are multiple, the encoder produces an array of bstr, one for each nonce.
func (ns Nonce) MarshalCBOR() ([]byte, error) {
	if err := ns.Validate(); err != nil {
		return nil, fmt.Errorf("CBOR encoding failed: %w", err)
	}

	if len(ns) == 1 {
		return em.Marshal(ns[0])
	}

	return em.Marshal([]nonce(ns))
}

// UnmarshalCBOR decodes a EAT nonce. This may be a single byte string
// between 8 and 64 bytes long, or an array of two or more such strings.
func (ns *Nonce) UnmarshalCBOR(data []byte) error {
	if isCBORArray(data) {
		return dm.Unmarshal(data, (*[]nonce)(ns))
	}

	var n nonce

	if err := dm.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("CBOR decoding failed for nonce: %w", err)
	}

	*ns = Nonce{n}

	return nil
}

// MarshalJSON encodes the receiver Nonce as either a JSON string containing
// the base64 encoding of the binary nonce (if the array comprises only one
// element) or as an array of base64-encoded JSON strings.
func (ns Nonce) MarshalJSON() ([]byte, error) {
	if err := ns.Validate(); err != nil {
		return nil, fmt.Errorf("JSON encoding failed: %w", err)
	}

	if len(ns) == 1 {
		return json.Marshal(ns[0])
	}

	return json.Marshal([]nonce(ns))
}

// UnmarshalJSON decodes a EAT nonce in JSON format.
func (ns *Nonce) UnmarshalJSON(data []byte) error {
	if isJSONArray(data) {
		return json.Unmarshal(data, (*[]nonce)(ns))
	}

	var n nonce

	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("JSON decoding failed for nonce: %w", err)
	}

	*ns = Nonce{n}

	return nil
}
