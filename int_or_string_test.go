// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IntOrStringFromAny(t *testing.T) {
	testCases := []struct {
		title string
		input any
		err   string
	}{
		{
			title: "ok int8",
			input: int8(1),
		},
		{
			title: "ok int16",
			input: int16(1),
		},
		{
			title: "ok int32",
			input: int32(1),
		},
		{
			title: "ok int64",
			input: int64(1),
		},
		{
			title: "ok uint8",
			input: uint8(1),
		},
		{
			title: "ok uint16",
			input: uint16(1),
		},
		{
			title: "ok uint32",
			input: uint32(1),
		},
		{
			title: "ok uint64",
			input: uint64(1),
		},
		{
			title: "ok float64",
			input: 1.0,
		},
		{
			title: "ok string",
			input: "foo",
		},
		{
			title: "err bool",
			input: true,
			err:   "unexpected algorithm value: true(bool)",
		},
		{
			title: "err struct",
			input: struct{}{},
			err:   "unexpected algorithm value: {}(struct {})",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := IntOrStringFromAny(tc.input)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestIntOrString_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		value        IntOrString
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title:        "ok zero",
			value:        IntOrStringFromInt(0),
			expectedCBOR: []byte{0x00},
			expectedJSON: "0",
		},
		{
			title:        "ok positive",
			value:        IntOrStringFromInt(42),
			expectedCBOR: []byte{0x18, 0x2a},
			expectedJSON: "42",
		},
		{
			title:        "ok negative",
			value:        IntOrStringFromInt(-1),
			expectedCBOR: []byte{0x20},
			expectedJSON: "-1",
		},
		{
			title: "ok string",
			value: IntOrStringFromString("foo"),
			expectedCBOR: []byte{
				0x63,             // tstr(3)
				0x66, 0x6f, 0x6f, // . "foo"
			},
			expectedJSON: `"foo"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.value.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedValue IntOrString
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

func TestIntOrString_UnmarshalCBOR_negative(t *testing.T) {
	var ios IntOrString
	err := ios.UnmarshalCBOR([]byte{})
	assert.ErrorContains(t, err, "truncated input")

	err = ios.UnmarshalCBOR([]byte{0xf5})
	assert.ErrorContains(t, err, "unexpected CBOR major type for DigestAlgID: 7")
}

func TestIntOrString_UnmarshalJSON_negative(t *testing.T) {
	var ios IntOrString
	err := ios.UnmarshalJSON([]byte{})
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	err = ios.UnmarshalJSON([]byte(`true`))
	assert.ErrorContains(t, err, "unexpected algorithm value: true(bool)")
}

func TestIntOrString_misc(t *testing.T) {
	ios := IntOrStringFromInt(1)
	assert.True(t, ios.IsInt())
	assert.False(t, ios.IsString())
	assert.Equal(t, "1", ios.String())
	assert.Equal(t, 1, ios.Int())

	ios = IntOrStringFromString("foo")
	assert.False(t, ios.IsInt())
	assert.True(t, ios.IsString())
	assert.Equal(t, "foo", ios.String())
	assert.Equal(t, 0, ios.Int())

	ios = IntOrStringFromString("1")
	assert.Equal(t, 1, ios.Int())
}

func TestIntOrString_JSON_key(t *testing.T) {
	toEncode := map[IntOrString]any{
		IntOrStringFromString("foo"): 1,
	}

	expectedJSON := `{"foo": 1}`
	encoded, err := json.Marshal(toEncode)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(encoded))

	var decoded map[IntOrString]any
	err = json.Unmarshal(encoded, &decoded)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, int(decoded[IntOrStringFromString("foo")].(float64)))
}

func TestIntOrString_text_marshaling(t *testing.T) {
	ios := IntOrStringFromInt(1)
	text, err := ios.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, []byte("1"), text)

	var other IntOrString
	err = other.UnmarshalText(text)
	assert.NoError(t, err)
	assert.Equal(t, ios, other)
}
