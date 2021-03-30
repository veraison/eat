// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonces_Add_ok(t *testing.T) {
	nonces := Nonces{}

	for i := MinNonceSize; i <= MaxNonceSize; i++ {
		tv := make([]byte, i)

		err := nonces.Add(tv)

		assert.Nil(t, err)
	}
}

func TestNonces_Add_too_small(t *testing.T) {
	nonces := Nonces{}

	tv := make([]byte, MinNonceSize-1)

	err := nonces.Add(tv)

	expectedError := fmt.Sprintf(
		"a nonce must be between %d and %d bytes long; found %d",
		MinNonceSize, MaxNonceSize, len(tv),
	)

	assert.EqualError(t, err, expectedError)
}

func TestNonces_Add_too_big(t *testing.T) {
	nonces := Nonces{}

	tv := make([]byte, MaxNonceSize+1)

	err := nonces.Add(tv)

	expectedError := fmt.Sprintf(
		"a nonce must be between %d and %d bytes long; found %d",
		MinNonceSize, MaxNonceSize, len(tv),
	)

	assert.EqualError(t, err, expectedError)
}

func TestNonces_AddHex_ok(t *testing.T) {
	nonces := Nonces{}

	err := nonces.AddHex("deadbeefbeefdead")

	assert.Nil(t, err)
}

func TestNonces_AddHex_bad_hex(t *testing.T) {
	nonces := Nonces{}

	err := nonces.AddHex("dea")

	expected := "decoding nonce failed: encoding/hex: odd length hex string"

	assert.EqualError(t, err, expected)
}

func TestNonces_MarshalCBOR_single_ok(t *testing.T) {
	tv := []byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	nonces := Nonces{}
	require.Nil(t, nonces.Add(tv))

	expected := []byte{
		0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	actual, err := nonces.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNonces_MarshalCBOR_multiple_ok(t *testing.T) {
	tv := [][]byte{
		{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef},
		{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe},
	}

	nonces := Nonces{}
	for i := range tv {
		require.Nil(t, nonces.Add(tv[i]))
	}

	//82                     # array(2)
	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	expected := []byte{0x82}
	expected = append(expected, byte(0x48))
	expected = append(expected, tv[0]...)
	expected = append(expected, byte(0x48))
	expected = append(expected, tv[1]...)

	actual, err := nonces.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNonces_UnmarshalCBOR_single_ok(t *testing.T) {
	expected := []byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	data := append([]byte{0x48}, expected...)

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	assert.Nil(t, err)
	assert.Equal(t, 1, actual.Len())
	assert.Equal(t, expected, actual.GetI(0))
}

func TestNonces_UnmarshalCBOR_multiple_ok(t *testing.T) {
	expected := [][]byte{
		{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef},
		{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe},
	}

	//82                     # array(2)
	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	data := []byte{0x82}
	data = append(data, byte(0x48))
	data = append(data, expected[0]...)
	data = append(data, byte(0x48))
	data = append(data, expected[1]...)

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	assert.Nil(t, err)
	assert.Equal(t, 2, actual.Len())
	assert.Equal(t, expected[0], actual.GetI(0))
	assert.Equal(t, expected[1], actual.GetI(1))
}

func TestNonce_UnmarshalCBOR_bad_cbor(t *testing.T) {
	// length (8) doesn't match the number of bytes (7)
	data := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca}

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	expected := "CBOR decoding failed for nonce: unexpected EOF"

	assert.EqualError(t, err, expected)
}

func TestNonces_Validate_ok(t *testing.T) {
	data := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	actual := Nonces{}
	require.Nil(t, actual.UnmarshalCBOR(data))

	err := actual.Validate()
	assert.Nil(t, err)
}

func TestNonces_Validate_empty(t *testing.T) {
	empty := Nonces{}

	err := empty.Validate()
	assert.EqualError(t, err, "empty nonce")
}

func TestNonces_Validate_too_short(t *testing.T) {
	// 7 bytes
	data := []byte{0x47, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe}

	actual := Nonces{}
	require.Nil(t, actual.UnmarshalCBOR(data))

	expected := "found invalid nonce at index 0: "
	expected += "a nonce must be between 8 and 64 bytes long; found 7"

	err := actual.Validate()
	assert.EqualError(t, err, expected)
}

func TestNonces_Validate_too_long(t *testing.T) {
	// 65 bytes
	data := []byte{
		0x58, 0x41, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde,
		0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe,
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad,
		0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde,
		0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde,
	}

	actual := Nonces{}
	require.Nil(t, actual.UnmarshalCBOR(data))

	expected := "found invalid nonce at index 0: "
	expected += "a nonce must be between 8 and 64 bytes long; found 65"

	err := actual.Validate()
	assert.EqualError(t, err, expected)
}

func TestNonces_UnmarshalCBOR_repeatedly(t *testing.T) {
	data1 := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	data2 := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	expected := []byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	actual := Nonces{}

	err := actual.UnmarshalCBOR(data1)
	assert.Nil(t, err)

	// the second decode on the same Nonces clobbers the first
	err = actual.UnmarshalCBOR(data2)
	assert.Nil(t, err)
	assert.Equal(t, 1, actual.Len())
	assert.Equal(t, expected, actual.GetI(0))
}

func TestNonce_MarshalJSON_single_ok(t *testing.T) {
	nonces := Nonces{}

	err := nonces.Add([]byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	})
	require.Nil(t, err)

	expected := []byte(`"3q2+796tvu8="`)

	actual, err := nonces.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNonces_MarshalJSON_multiple_ok(t *testing.T) {
	tv := [][]byte{
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
	}

	nonces := Nonces{}
	for i := range tv {
		require.Nil(t, nonces.Add(tv[i]))
	}

	expected := `[
		"AAAAAAAAAAA=",
		"AQEBAQEBAQE="
	]`

	actual, err := nonces.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestNonces_UnmarshalJSON_single_ok(t *testing.T) {
	tv := []byte(`"3q2+796tvu8="`)

	expected := []byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.Nil(t, err)
	assert.Equal(t, 1, actual.Len())
	assert.Equal(t, expected, actual.GetI(0))
}

func TestNonces_UnmarshalJSON_invalid_json(t *testing.T) {
	tv := []byte(`"unterminated string`)

	expected := "JSON decoding failed for nonce: "
	expected += "unexpected end of JSON input"

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, expected)
}

func TestNonces_UnmarshalJSON_invalid_base64(t *testing.T) {
	tv := []byte(`"0"`)

	expected := "JSON decoding failed for nonce: "
	expected += "illegal base64 data at input byte 0"

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, expected)
}

func TestNonces_UnmarshalJSON_not_a_string(t *testing.T) {
	tv := []byte(`{ "a": 1 }`)

	expected := "JSON decoding failed for nonce: "
	expected += "invalid nonce input map[string]interface {}"

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, expected)
}

func TestNonces_GetI_ok(t *testing.T) {
	tv := [][]byte{
		{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef},
		{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe},
	}

	nonces := Nonces{}
	for i := range tv {
		require.Nil(t, nonces.Add(tv[i]))
	}

	for i := range tv {
		assert.Equal(t, tv[i], nonces.GetI(i))
	}
}

func TestNonces_GetI_out_of_bounds(t *testing.T) {
	tv := Nonces{}

	for i := -10; i < 10; i++ {
		actual := tv.GetI(i)
		assert.Nil(t, actual)
	}
}

func TestNonces_UnmarshalJSON_two_entries_ok(t *testing.T) {
	tv := []byte(`[
		"AAAAAAAAAAA=",
		"AQEBAQEBAQE="
	]`)

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	expected := [][]byte{
		{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
	}

	assert.Nil(t, err)
	for i := range expected {
		assert.Equal(t, expected[i], actual.GetI(i))
	}
}
