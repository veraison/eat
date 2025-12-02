// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

// StringOrBinary is either tstr or binary-data (B64Url, to be encoded with base64url for JSON)
type StringOrBinary struct {
	Value interface{}
}

func (s StringOrBinary) MarshalJSON() ([]byte, error) {
	switch s.Value.(type) {
	case B64Url:
		return json.Marshal(s.Value)
	case string:
		return json.Marshal(s.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", s.Value)
	}
}

func (s *StringOrBinary) UnmarshalJSON(data []byte) error {
	// try to decode as base64url
	var b B64Url
	if err := json.Unmarshal(data, &b); err == nil {
		s.Value = b
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.Value = str
		return nil
	}

	return fmt.Errorf("invalid value: expected tstr or binary-data")
}

func (s StringOrBinary) MarshalCBOR() ([]byte, error) {
	switch s.Value.(type) {
	case B64Url:
		return em.Marshal(s.Value)
	case string:
		return em.Marshal(s.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", s.Value)
	}
}

func (s *StringOrBinary) UnmarshalCBOR(data []byte) error {
	var b B64Url
	if err := dm.Unmarshal(data, &b); err == nil {
		s.Value = b
		return nil
	}

	var str string
	if err := dm.Unmarshal(data, &str); err == nil {
		s.Value = str
		return nil
	}

	return fmt.Errorf("invalid value: expected tstr or binary-data")
}
