// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// OemID is the Original Equipment Manufacturer identifier. It comes in three
// flavors:
//
//   - Random number-based. This is a 16-byte random values. This allows OEMs to
//     generate their own IDs using a cryptographic-quality random number
//     generator.
//   - IEEE-based. Obtained from the IEEE global registry of MAC addresses and
//     company IDs. This is a 3-byte value.
//   - IANA Private Enterprise Number-based. Taken from the IANA PEN registry.
//     This is an integer value (taken from the OID in the 1.3.6.1.4.1 namespace
//     managed by IANA).
type OemID struct {
	value any
}

// NewPenOemID returns a pointer to a new OemID containing the specified PEN
// ID.
func NewPenOemID(value int) *OemID {
	return &OemID{value}
}

// NewIeeeOemID returns a pointer to a new OemID containing the specified
// IEEE-based ID.
func NewIeeeOemID(value [3]byte) *OemID {
	return &OemID{BinaryData(value[:])}
}

// NewRandomOemID returns a pointer to a new OemID containing the specified
// random ID.
func NewRandomOemID(value [16]byte) *OemID {
	return &OemID{BinaryData(value[:])}
}

// IsInt return true if the underlying value is an int.
func (o *OemID) IsInt() bool {
	_, ok := o.value.(int)
	return ok
}

// IsBytes returns true if the underlying value is a []byte.
func (o *OemID) IsBytes() bool {
	_, ok := o.value.(BinaryData)
	return ok
}

// Int returns the underlying int value. If the underlying value is a []byte, 0
// is returned instead.
func (o *OemID) Int() int {
	value, ok := o.value.(int)
	if !ok {
		return 0
	}

	return value
}

// Bytes returns the underlying []byte value. If the underlying value is an
// int, nil is returned instead.
func (o *OemID) Bytes() []byte {
	value, ok := o.value.(BinaryData)
	if !ok {
		return nil
	}

	return value
}

// IsPEN returns true if the underlying value is an AINA Private Enterprise
// Number ID.
func (o *OemID) IsPEN() bool {
	return o.IsInt()
}

// IsIEEE returns true if the underlying is taken from the IEEE global
// registry.
func (o *OemID) IsIEEE() bool {
	return len(o.Bytes()) == 3
}

// IsRandom returns true if the underlying value is randomly generated.
func (o *OemID) IsRandom() bool {
	return len(o.Bytes()) == 16
}

func (o *OemID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

func (o *OemID) UnmarshalCBOR(data []byte) error {
	if len(data) < 1 {
		return errors.New("truncated input")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 0, 1:
		var value int
		if err := dm.Unmarshal(data, &value); err != nil {
			return err
		}

		o.value = value
	case 2:
		var value BinaryData
		if err := dm.Unmarshal(data, &value); err != nil {
			return err
		}

		if len(value) != 3 && len(value) != 16 {
			return fmt.Errorf("incorrect length %d for OemID (must be 3 or 16)", len(value))
		}

		o.value = value
	default:
		return fmt.Errorf("unexpected major type %d for OemID", majorType)
	}

	return nil
}

func (o *OemID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *OemID) UnmarshalJSON(data []byte) error {
	text := strings.TrimSpace(string(data))

	if len(text) < 1 {
		return errors.New("truncated input")
	}

	switch text[0] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		var value int
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}

		o.value = value
	case '"':
		var value BinaryData
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}

		if len(value) != 3 && len(value) != 16 {
			return fmt.Errorf("incorrect length %d for OemID (must be 3 or 16)", len(value))
		}

		o.value = value
	default:
		return fmt.Errorf("invalid OemID (must be int or string): %s", text)
	}

	return nil
}
