// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
)

// A nonce-claim may be single Nonce or an array of two or more.
//
//    nonce-claim = (
//        nonce => nonce-type / [ 2* nonce-type ]
//    )
type Nonces []Nonce

// MarshalCBOR encodes Nonces as a byte string, in case there is only
// one, or an array of byte strings if there are multiple.
func (ns Nonces) MarshalCBOR() ([]byte, error) {
	if len(ns) == 1 {
		return em.Marshal(ns[0])
	}

	return em.Marshal([]Nonce(ns))
}

// UnmarshalCBOR decodes nonce claim data. This may be a single byte string
// between 8 and 64 bytes long, or an array of two or more such strings.
func (ns *Nonces) UnmarshalCBOR(data []byte) error {
	if isCBORArray(data) {
		return dm.Unmarshal(data, (*[]Nonce)(ns))
	}

	var n Nonce

	if err := dm.Unmarshal(data, &n); err != nil {
		return err
	}

	*ns = Nonces{n}

	return nil
}

// Validate checks that all Nonce's are valid.
func (ns Nonces) Validate() error {
	for _, n := range ns {
		if err := n.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// NonceFromHex creates a new Nonce instance from a string containing
// hex-encoded byte values and returns a pointer to it.
func NonceFromHex(text string) (*Nonce, error) {
	value, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return NewNonce(value), nil
}

// A Nonce is between 8 and 64 bytes
//    nonce-type = bstr .size (8..64)
type Nonce struct {
	value []byte
}

// NewNonce returns a Nonce initialized with the supplied byte slice
func NewNonce(v []byte) *Nonce {
	return &Nonce{v}
}

// Get returns the nonce value
func (n Nonce) Get() []byte {
	return n.value
}

// MarshalCBOR encodes the Nonce a CBOR byte string.
func (n Nonce) MarshalCBOR() ([]byte, error) {
	return em.Marshal(n.value)
}

// UnmarshalCBOR decodes a CBOR byte string and uses it as the Nonce value.
func (n *Nonce) UnmarshalCBOR(data []byte) error {
	var value []byte

	if err := dm.Unmarshal(data, &value); err != nil {
		return err
	}

	n.value = value

	return nil
}

// Validate checks that the underlying value of the Nonce is between 8 and 64
// bytes, as is required by the EAT spec.
func (n Nonce) Validate() error {
	vlen := len(n.value)
	if vlen > 64 || vlen < 8 {
		return fmt.Errorf("a nonce must be between 8 and 64 bytes long; found %v", vlen)
	}
	return nil
}

// MarshalJSON encodes the receiver Nonces as a JSON string
func (ns Nonces) MarshalJSON() ([]byte, error) {
	if len(ns) == 1 {
		return json.Marshal(ns[0].value)
	}

	return nil, errors.New("TODO handle array of nonces")
}

func (ns *Nonces) UnmarshalJSON(data []byte) error {
	var v interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch t := v.(type) {
	case string:
		value, err := base64.StdEncoding.DecodeString(t)
		if err != nil {
			return err
		}
		*ns = Nonces{Nonce{value}}
		return nil
	default:
		return errors.New("TODO handle array of nonces")
	}
}
