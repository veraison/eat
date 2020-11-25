// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
)

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
type Submods map[string]Submod

// Get retrieves a submod by name (either int64 or string)
func (s Submods) Get(name string) interface{} {
	return s[name].value
}

// Add inserts the named submod in the Submods container. The supplied name must
// be of type string or int64
func (s *Submods) Add(name string, submod interface{}) error {
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

	(*s)[name] = Submod{submod}

	return nil
}
