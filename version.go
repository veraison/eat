// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
)

type Version struct {
	_       struct{} `cbor:",toarray"`
	Version string
	Scheme  *VersionScheme
}

func (v Version) MarshalCBOR() ([]byte, error) {
	r := []interface{}{v.Version}
	if v.Scheme != nil {
		r = append(r, *v.Scheme)
	}
	return cbor.Marshal(r)
}

func (v *Version) UnmarshalCBOR(data []byte) error {
	var raw []interface{}
	if err := cbor.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 1 || len(raw) > 2 {
		return fmt.Errorf("invalid Version CBOR array length: %d", len(raw))
	}

	version, ok := raw[0].(string)
	if !ok {
		return fmt.Errorf("invalid Version type: expected string")
	}
	v.Version = version
	if len(raw) == 2 {
		var scheme VersionScheme
		switch v := raw[1].(type) {
		case int:
			scheme = VersionScheme(v)
		case uint64:
			if v > uint64(^uint(0)>>1) {
				return fmt.Errorf("invalid Version Scheme value: %d", v)
			}
			scheme = VersionScheme(v)
		default:
			return fmt.Errorf("invalid Version Scheme type: expected int %T", raw[1])
		}
		v.Scheme = &scheme
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
	var raw []interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw) < 1 || len(raw) > 2 {
		return fmt.Errorf("invalid Version CBOR array length: %d", len(raw))
	}

	version, ok := raw[0].(string)
	if !ok {
		return fmt.Errorf("invalid Version type: expected string")
	}
	v.Version = version

	if len(raw) == 2 {
		var scheme VersionScheme
		switch v := raw[1].(type) {
		case float64:
			scheme = VersionScheme(int(v))
		case string:
			loc, ok := stringToVersionScheme[v]
			if !ok {
				return fmt.Errorf("invalid VersionScheme string: %s", v)
			}
			scheme = loc
		default:
			return fmt.Errorf("invalid Version Scheme type: expected int %T", raw[1])
		}
		v.Scheme = &scheme
	} else {
		v.Scheme = nil
	}

	return nil
}
