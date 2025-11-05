// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/veraison/swid"
)

type Version struct {
	_       struct{} `cbor:",toarray"`
	Version string
	Scheme  *swid.VersionScheme
}

func (v Version) MarshalCBOR() ([]byte, error) {
	r := []interface{}{v.Version}
	if v.Scheme != nil {
		r = append(r, *v.Scheme)
	}
	return cbor.Marshal(r)
}

func (v *Version) UnmarshalCBOR(data []byte) error {
	var raw []cbor.RawMessage
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 1 || len(raw) > 2 {
		return fmt.Errorf("invalid Version CBOR array length: %d", len(raw))
	}

	if err := cbor.Unmarshal(raw[0], &v.Version); err != nil {
		return fmt.Errorf("invalid Version type: expected string")
	}
	if len(raw) == 1 {
		// no version scheme
		return nil
	}

	if err := cbor.Unmarshal(raw[1], &v.Scheme); err != nil {
		return err
	}

	return nil
}

func (v Version) MarshalJSON() ([]byte, error) {
	r := []interface{}{v.Version}
	if v.Scheme != nil {
		r = append(r, *v.Scheme)
	}
	return json.Marshal(r)
}

func (v *Version) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 1 || len(raw) > 2 {
		return fmt.Errorf("invalid Version CBOR array length: %d", len(raw))
	}

	if err := json.Unmarshal(raw[0], &v.Version); err != nil {
		return fmt.Errorf("invalid Version type: expected string")
	}
	if len(raw) == 1 {
		// no version scheme
		return nil
	}

	if err := json.Unmarshal(raw[1], &v.Scheme); err != nil {
		return err
	}

	return nil
}
