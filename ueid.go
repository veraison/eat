// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

const (
	UEIDTypeInvalid = iota

	// This is a 128, 192 or 256 bit random number generated once and
	// stored in the device. This may be constructed by concatenating
	// enough identifiers to make up an equivalent number of random bits
	// and then feeding the concatenation through a cryptographic hash
	// function. It may also be a cryptographic quality random number
	// generated once at the beginning of the life of the device and
	// stored. It may not be smaller than 128 bits.
	UEIDTypeRAND

	// This makes use of the IEEE company identification registry. An EUI
	// is either an EUI-48, EUI-60 or EUI-64 and made up of an OUI, OUI-36
	// or a CID, different registered company identifiers, and some unique
	// per-device identifier. EUIs are often the same as or similar to MAC
	// addresses. This type includes MAC-48, an obsolete name for EUI-48.
	// (Note that while devices with multiple network interfaces may have
	// multiple MAC addresses, there is only one UEID for a device)
	UEIDTypeEUI

	// This is a 14-digit identifier consisting of an 8-digit Type
	// Allocation Code and a 6-digit serial number allocated by the
	// manufacturer, which SHALL be encoded as byte string of length 14
	// with each byte as the digit's value (not the ASCII encoding of the
	// digit; the digit 3 encodes as 0x03, not 0x33). The IMEI value
	// encoded SHALL NOT include Luhn checksum or SVN information.
	UEIDTypeIMEI
)

func MustUEIDFromBytes(value []byte) *UEID {
	ret, err := UEIDFromBytes(value)
	if err != nil {
		panic(err)
	}

	return ret
}

func UEIDFromBytes(value []byte) (*UEID, error) {
	if err := validateUEIDBytes(value); err != nil {
		return nil, err
	}

	return &UEID{value}, nil
}

type UEID struct {
	value BinaryData
}

func (o *UEID) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

func (o *UEID) UnmarshalCBOR(data []byte) error {
	var value BinaryData
	if err := dm.Unmarshal(data, &value); err != nil {
		return err
	}

	if err := validateUEIDBytes(value); err != nil {
		return err
	}

	o.value = value

	return nil
}

func (o *UEID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *UEID) UnmarshalJSON(data []byte) error {
	var value BinaryData
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if err := validateUEIDBytes(value); err != nil {
		return err
	}

	o.value = value

	return nil
}

func validateUEIDBytes(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty UEID")
	}

	typ := data[0]
	value := data[1:]

	switch typ {
	case UEIDTypeRAND:
		return validateRAND(value)
	case UEIDTypeEUI:
		return validateEUI(value)
	case UEIDTypeIMEI:
		return validateIMEI(value)
	default:
		return fmt.Errorf("invalid UEID type %v", typ)
	}
}

func validateRAND(value []byte) error {
	vlen := len(value)
	if vlen != 16 && vlen != 24 && vlen != 32 {
		return fmt.Errorf("RAND length must be exactly 16, 24, or 32 bytes; found %v bytes", vlen)
	}
	return nil
}

func validateEUI(value []byte) error {
	vlen := len(value)
	if vlen != 6 && vlen != 8 {
		return fmt.Errorf("EUI length must be exactly 6 (EUI-48) or 8 (EUI-60 or EUI-64) bytes; found %v bytes", vlen)
	}
	return nil
}

func validateIMEI(value []byte) error {
	vlen := len(value)
	if vlen != 14 {
		return fmt.Errorf("IMEI length must be exactly 14 bytes; found %v bytes", vlen)
	}
	return nil
}
