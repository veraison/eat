// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// IntOrString has an underlying value that is either an int or string.
type IntOrString struct {
	val any
}

// IntOrStringFromInt creates a new IntOrString with the specified int as the
// underlying value.
func IntOrStringFromInt(val int) IntOrString {
	return IntOrString{val}
}

// IntOrStringFromInt creates a new IntOrString with the specified string as
// the underlying value.
func IntOrStringFromString(val string) IntOrString {
	var ret IntOrString

	i, err := strconv.Atoi(val)
	if err == nil {
		ret.val = i
	} else {
		ret.val = val
	}

	return ret
}

// IntOrStringFromAny creates a new IntOrString with the provided underlying
// value. If the provided value is not a string or an integer type, an error is
// returned.
func IntOrStringFromAny(val any) (IntOrString, error) {
	var ret IntOrString
	switch t := val.(type) {
	case int8:
		ret.val = int(t)
	case uint8:
		ret.val = int(t)
	case int16:
		ret.val = int(t)
	case uint16:
		ret.val = int(t)
	case int32:
		ret.val = int(t)
	case uint32:
		ret.val = int(t)
	case int64:
		ret.val = int(t)
	case uint64:
		if t > math.MaxInt64 {
			return IntOrString{}, fmt.Errorf("out of range: %d", t)
		}
		ret.val = int(t)
	case float64:
		ret.val = int(t)
	case string:
		i, err := strconv.Atoi(t)
		if err == nil {
			ret.val = i
		} else {
			ret.val = t
		}
	default:
		return IntOrString{}, fmt.Errorf("unexpected algorithm value: %v(%T)", t, t)
	}

	return ret, nil
}

// IsString returns true if the underlying value is a string.
func (o *IntOrString) IsString() bool {
	_, ok := o.val.(string)
	return ok
}

// IsString returns true if the underlying value is an int.
func (o *IntOrString) IsInt() bool {
	_, ok := o.val.(int)
	return ok
}

// String returns the string representation of the underlying value.
func (o *IntOrString) String() string {
	switch t := o.val.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	default:
		return ""
	}
}

// Int returns an integer representation of the underlying value. If the
// underlying value is a string that cannot be parsed as an int, 0 is returned.
func (o *IntOrString) Int() int {
	switch t := o.val.(type) {
	case int:
		return t
	case string:
		i, err := strconv.Atoi(t)
		if err != nil {
			return 0
		}

		return i
	default:
		return 0
	}
}

func (o *IntOrString) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *IntOrString) UnmarshalCBOR(data []byte) error { // nolint:dupl
	if len(data) == 0 {
		return errors.New("truncated input")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 0, 1:
		var val int
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	case 3:
		var val string
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	default:
		return fmt.Errorf("unexpected CBOR major type for DigestAlgID: %d", majorType)
	}
}

func (o *IntOrString) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.val)
}

func (o *IntOrString) UnmarshalJSON(data []byte) error {
	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	ios, err := IntOrStringFromAny(val)
	if err != nil {
		return err
	}

	o.val = ios.val

	return nil
}
