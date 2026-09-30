// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOemID_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		value        OemID
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title:        "ok int",
			value:        OemID{1},
			expectedCBOR: []byte{0x01},
			expectedJSON: "1",
		},
		{
			title: "ok bytes",
			value: OemID{BinaryData{0xff, 0xff, 0xff}},
			expectedCBOR: []byte{
				0x43, // bstr(3)
				0xff, 0xff, 0xff,
			},
			expectedJSON: `"____"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.value.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedValue OemID
			err = decodedValue.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.EqualValues(t, decodedValue, tc.value)

			encodedJSON, err := tc.value.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = decodedValue.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.EqualValues(t, decodedValue, tc.value)
		})
	}
}

func TestOemID_UnmarshalCBOR_negative(t *testing.T) {
	var id OemID
	err := id.UnmarshalCBOR([]byte{})
	assert.ErrorContains(t, err, "truncated input")
	err = id.UnmarshalCBOR([]byte{0xf5})
	assert.ErrorContains(t, err, "unexpected major type 7 for OemID")
	err = id.UnmarshalCBOR([]byte{0x44, 0xff, 0xff, 0xff, 0xff})
	assert.ErrorContains(t, err, "incorrect length 4 for OemID (must be 3 or 16)")
}

func TestOemID_UnmarshalJSON_negative(t *testing.T) {
	var id OemID
	err := id.UnmarshalJSON([]byte{})
	assert.ErrorContains(t, err, "truncated input")
	err = id.UnmarshalJSON([]byte("{}"))
	assert.ErrorContains(t, err, "invalid OemID (must be int or string)")
	err = id.UnmarshalJSON([]byte(`"______"`))
	assert.ErrorContains(t, err, "incorrect length 4 for OemID (must be 3 or 16)")
}

func TestOemID_misc(t *testing.T) {
	id := NewPenOemID(1)
	assert.True(t, id.IsInt())
	assert.False(t, id.IsBytes())
	assert.True(t, id.IsPEN())
	assert.False(t, id.IsIEEE())
	assert.False(t, id.IsRandom())
	assert.Equal(t, 1, id.Int())
	assert.Nil(t, id.Bytes())

	ieeeBytes := [3]byte{0xff, 0xff, 0xff}
	id = NewIeeeOemID(ieeeBytes)
	assert.False(t, id.IsInt())
	assert.True(t, id.IsBytes())
	assert.False(t, id.IsPEN())
	assert.True(t, id.IsIEEE())
	assert.False(t, id.IsRandom())
	assert.Equal(t, 0, id.Int())
	assert.Equal(t, ieeeBytes[:], id.Bytes())

	randomBytes := [16]byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	id = NewRandomOemID(randomBytes)
	assert.False(t, id.IsInt())
	assert.True(t, id.IsBytes())
	assert.False(t, id.IsPEN())
	assert.False(t, id.IsIEEE())
	assert.True(t, id.IsRandom())
	assert.Equal(t, 0, id.Int())
	assert.Equal(t, randomBytes[:], id.Bytes())
}
