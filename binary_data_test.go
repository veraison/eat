// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestBinaryData_marshaling(t *testing.T) {
	data := BinaryData(MustHexDecode(t, "deadbeef"))

	expected := []byte{
		0x44, // bstr(4)
		0xde, 0xad, 0xbe, 0xef,
	}
	encoded, err := cbor.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, encoded)

	expected = []byte(`"3q2-7w"`)
	encoded, err = json.Marshal(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, encoded)
}

func TestBinaryData_unmarshaling(t *testing.T) {
	t.Run("cbor", func(t *testing.T) {
		var data BinaryData
		encoded := MustHexDecode(t, "44deadbeef")

		err := cbor.Unmarshal(encoded, &data)
		assert.NoError(t, err)
		assert.Equal(t, BinaryData(MustHexDecode(t, "deadbeef")), data)
	})

	testCases := []struct {
		title    string
		input    string
		expected BinaryData
		err      string
	}{
		{
			title:    "json ok",
			input:    `"3q2-7w"`,
			expected: BinaryData(MustHexDecode(t, "deadbeef")),
		},
		{
			title: "json bad base64",
			input: `"3q2+7w=="`,
			err:   "illegal base64 data at input byte 3",
		},
		{
			title: "json wrong type",
			input: `42`,
			err:   "json: cannot unmarshal number into Go value of type string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var out BinaryData
			err := json.Unmarshal([]byte(tc.input), &out)

			if tc.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, out)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}
