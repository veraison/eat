// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// A nonce-claim may be single Nonce or an array of two or more.
//
//    nonce-claim = (
//        nonce => nonce-type / [ 2* nonce-type ]
//
type Nonces []nonce

// nonce-type = bstr .size (8..64)
const (
	MinNonceSize = 8
	MaxNonceSize = 64
)

// MarshalCBOR provides a suitable CBOR encoding for the receiver Nonces.
// In case there is only one nonce, it is encoded as a bstr. If there are
// multiple, it is encoded as an array of bstr.
func (ns Nonces) MarshalCBOR() ([]byte, error) {
	if err := ns.Validate(); err != nil {
		return nil, fmt.Errorf("CBOR encoding failed: %w", err)
	}

	if len(ns) == 1 {
		return em.Marshal(ns[0])
	}

	return em.Marshal([]nonce(ns))
}

// UnmarshalCBOR decodes nonce claim data. This may be a single byte string
// between 8 and 64 bytes long, or an array of two or more such strings.
func (ns *Nonces) UnmarshalCBOR(data []byte) error {
	if isCBORArray(data) {
		return dm.Unmarshal(data, (*[]nonce)(ns))
	}

	var n nonce

	if err := dm.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("CBOR decoding failed for nonce: %w", err)
	}

	*ns = Nonces{n}

	return nil
}

// Validate checks that all nonces are valid.
func (ns Nonces) Validate() error {
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

// Add the supplied nonce value to the receiver Nonces. There can be multiple
// nonce values carried in a single Nonces instance.
func (ns *Nonces) Add(v []byte) error {
	n, err := newNonce(v)
	if err != nil {
		return err
	}

	*ns = append(*ns, *n)

	return nil
}

// Len returns the number of nonce values carried in the Nonces
func (ns Nonces) Len() int {
	return len(ns)
}

// GetI returns the nonce found at index (starting at 0) or nil if the index is
// out of bounds
func (ns Nonces) GetI(index int) []byte {
	if index < 0 || index >= ns.Len() {
		return nil
	}

	return ns[index].get()
}

// AddHex is the same as Add except it takes the nonce value as a string
// containing hex-encoded byte values
func (ns *Nonces) AddHex(text string) error {
	value, err := hex.DecodeString(text)
	if err != nil {
		return fmt.Errorf("decoding nonce failed: %w", err)
	}

	return ns.Add(value)
}

// A nonce is between 8 and 64 bytes
//    nonce-type = bstr .size (8..64)
type nonce struct {
	value []byte
}

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

// newNonce returns a nonce initialized with the supplied byte slice or an error
// if the supplied buffer is either too big (more than 64 bytes) or too small
// (less than 8 bytes)
func newNonce(v []byte) (*nonce, error) {
	if err := isValidNonce(v); err != nil {
		return nil, err
	}

	return &nonce{v}, nil
}

// Get returns the nonce value
func (n nonce) get() []byte {
	return n.value
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

// validate checks that the underlying value of the Nonce is between 8 and 64
// bytes, as is required by the EAT spec
func (n nonce) validate() error {
	return isValidNonce(n.value)
}

// MarshalJSON encodes the receiver Nonces as either a JSON string containing
// the base64 encoding of the binary nonce (if the array
// comprises only one element) or as an array of JSON strings.
func (ns Nonces) MarshalJSON() ([]byte, error) {
	if err := ns.Validate(); err != nil {
		return nil, fmt.Errorf("JSON encoding failed: %w", err)
	}

	if len(ns) == 1 {
		return json.Marshal(ns[0])
	}

	return json.Marshal([]nonce(ns))
}

// UnmarshalJSON decodes the EAT nonce in JSON format
func (ns *Nonces) UnmarshalJSON(data []byte) error {
	if isJSONArray(data) {
		return json.Unmarshal(data, (*[]nonce)(ns))
	}

	var n nonce

	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("JSON decoding failed for nonce: %w", err)
	}

	*ns = Nonces{n}

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
