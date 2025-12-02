// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

type OEMID struct {
	// The type of oemid has multiple types, []byte or int
	//
	// oemid-type => oemid-pen / oemid-ieee / oemid-random
	// oemid-pen = int
	// oemid-ieee-cbor = bstr .size 3
	// oemid-ieee-json = base64-url-text .size 4
	// oemid-random-cbor = bstr .size 16
	// oemid-random-json = base64-url-text .size 24
	Value interface{}
}

func (o OEMID) Valid() bool {
	switch v := o.Value.(type) {
	case B64Url:
		t := v
		return len(t) == 3 || len(t) == 16
	case int:
		return true
	default:
		return false
	}
}

func (o OEMID) MarshalJSON() ([]byte, error) {
	if !o.Valid() {
		return nil, fmt.Errorf("invalid value %#v", o.Value)
	}
	switch o.Value.(type) {
	case B64Url:
		return json.Marshal(o.Value)
	case int:
		return json.Marshal(o.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", o.Value)
	}
}

func (o *OEMID) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		o.Value = i
		return nil
	}

	var s B64Url
	if err := json.Unmarshal(data, &s); err == nil {
		if len(s) != 3 && len(s) != 16 {
			return fmt.Errorf("%v is neither oemid-ieee nor oemid-random", data)
		}
		o.Value = s
		return nil
	}

	return fmt.Errorf("invalid value: expected oemid-pem, oemid-ieee or oemid-random")
}

func (o OEMID) MarshalCBOR() ([]byte, error) {
	if !o.Valid() {
		return nil, fmt.Errorf("invalid value %#v", o.Value)
	}
	switch o.Value.(type) {
	case B64Url:
		return em.Marshal(o.Value)
	case int:
		return em.Marshal(o.Value)
	default:
		return nil, fmt.Errorf("invalid type %T", o.Value)
	}
}

func (o *OEMID) UnmarshalCBOR(data []byte) error {
	var i int
	if err := dm.Unmarshal(data, &i); err == nil {
		o.Value = i
		return nil
	}

	var b B64Url
	if err := dm.Unmarshal(data, &b); err == nil {
		if len(b) != 3 && len(b) != 16 {
			return fmt.Errorf("%v is neither oemid-ieee nor oemid-random", data)
		}
		o.Value = b
		return nil
	}

	return fmt.Errorf("invalid value: expected oemid-pem, oemid-ieee or oemid-random")
}
