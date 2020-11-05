// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import cbor "github.com/fxamacker/cbor/v2"

// Eat is the internal representation of a EAT token
type Eat struct{}

// FromCBOR deserializes the supplied CBOR encoded EAT into the receiver Eat
func (e *Eat) FromCBOR(data []byte) error {
	return cbor.Unmarshal(data, e)
}

// ToCBOR serializes the receiver Eat into CBOR encoded EAT
func (e Eat) ToCBOR() ([]byte, error) {
	return cbor.Marshal(e)
}
