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
		switch raw[1].(type) {
		case int:
			scheme = VersionScheme(raw[1].(int))
		case uint64:
			s, _ := raw[1].(uint64)
			if s > uint64(^uint(0)>>1) {
				return fmt.Errorf("invalid Version Scheme value: %d", s)
			}
			scheme = VersionScheme(s)
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
		switch raw[1].(type) {
		case float64:
			f, _ := raw[1].(float64)
			scheme = VersionScheme(int(f))
		case string:
			loc, ok := stringToVersionScheme[raw[1].(string)]
			if !ok {
				return fmt.Errorf("invalid VersionScheme string: %s", raw[1].(string))
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
