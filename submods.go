// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// SubmodName is the type of a submod name. A custom type is needed to implement
// a TextMarshaler to encode the naked interface{} key into a string suitable
// for the JSON map.
type SubmodName struct{ value interface{} }

// MarshalText marshals the receiver SubmodName's value into a string
func (sn SubmodName) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%v", sn.value)), nil
}

// UnmarshalText is the dual of MarshalText
func (sn *SubmodName) UnmarshalText(text []byte) error {
	s := string(text)

	i, err := strconv.Atoi(s)
	if err != nil {
		sn.value = s
	} else {
		sn.value = int64(i)
	}

	return nil
}

// MarshalCBOR encodes the submod name wrapped in the SubmodName receiver to
// CBOR
func (sn SubmodName) MarshalCBOR() ([]byte, error) {
	return em.Marshal(sn.value)
}

// UnmarshalCBOR decodes the supplied data into the SubmodName receiver
func (sn *SubmodName) UnmarshalCBOR(data []byte) error {
	var v interface{}

	if err := dm.Unmarshal(data, &v); err != nil {
		return err
	}

	switch t := v.(type) {
	case uint64:
		if t > math.MaxInt64 {
			return errors.New("submod name too big")
		}
		sn.value = int64(t)
	case int64:
		sn.value = t
	case string:
		sn.value = t
	default:
		return fmt.Errorf("submod name must be string or (u)int64, found %T", t)
	}

	return nil
}

// Submod is the type of a submod: either a raw EAT (wrapped in a Sign1 CWT), or
// a map of EAT claims
type Submod struct{ value interface{} }

// MarshalJSON encodes the submod value wrapped in the Submod receiver to JSON
func (s Submod) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.value)
}

// MarshalCBOR encodes the submod value wrapped in the Submod receiver to CBOR
func (s Submod) MarshalCBOR() ([]byte, error) {
	return em.Marshal(s.value)
}

// UnmarshalJSON attempts to decode the supplied JSON data into the Submod
// receiver, peeking into the stream to choose between one of the two target
// formats (i.e., eat-token or eat-claims)
func (s *Submod) UnmarshalJSON(data []byte) error {
	if data[0] == '{' { // eat-claims
		var eatClaims Eat

		if err := eatClaims.FromJSON(data); err != nil {
			return err
		}
		s.value = eatClaims

		return nil
	}

	// eat-token
	b64 := string(data[1 : len(data)-1]) // remove quotes

	eatToken, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}

	if err := s.setEatToken(eatToken); err != nil {
		return err
	}

	return nil
}

func (s *Submod) setEatToken(data []byte) error {
	if err := checkTags(data); err != nil {
		return err
	}

	s.value = data

	return nil
}

// UnmarshalCBOR attempts to decode the supplied CBOR data into the Submod
// receiver, peeking into the stream to choose between one of the two target
// formats (i.e., eat-token or eat-claims)
func (s *Submod) UnmarshalCBOR(data []byte) error {
	if isCBORByteString(data) {
		var eatToken []byte

		if err := dm.Unmarshal(data, &eatToken); err != nil {
			return err
		}

		if err := s.setEatToken(eatToken); err != nil {
			return err
		}

		return nil
	}

	var eatClaims Eat
	if err := eatClaims.FromCBOR(data); err != nil {
		return err
	}

	s.value = eatClaims

	return nil
}

func checkTags(data []byte) error {
	// d8 3d  # tag(61) -- CWT
	// d2  # tag(18) -- Sign1
	prefix := []byte{0xd8, 0x3d, 0xd2}

	if len(data) < len(prefix)+1 {
		return errors.New("not enough bytes")
	}

	if !bytes.HasPrefix(data, prefix) {
		return errors.New("CWT and COSE Sign1 tags not found")
	}

	return nil
}

// Submods models the submods type
type Submods map[SubmodName]Submod

// Get retrieves a submod by name (either int64 or string)
func (s Submods) Get(name interface{}) interface{} {
	switch name.(type) {
	case string, int64:
		return s[SubmodName{name}].value
	default:
		return nil
	}
}

// Add inserts the named submod in the Submods container. The supplied name must
// be of type string or int64
func (s *Submods) Add(name interface{}, submod interface{}) error {
	switch t := name.(type) {
	case string, int64: // OK
	default:
		return fmt.Errorf("submod name must be string or int64, found %T", t)
	}

	switch t := submod.(type) {
	case Eat: // OK as-is
	case []byte: // make sure that the wrapping tags are in the right place
		if err := checkTags(t); err != nil {
			return err
		}
	default:
		return errors.New("submod must be Eat or []byte")
	}

	if *s == nil {
		*s = make(Submods)
	}

	(*s)[SubmodName{name}] = Submod{submod}

	return nil
}
