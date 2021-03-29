// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNonce_NewNonce_ok(t *testing.T) {
	for i := MinNonceSize; i <= MaxNonceSize; i++ {
		tv := make([]byte, i)

		nonce, err := NewNonce(tv)

		assert.Nil(t, err)
		assert.NotNil(t, nonce)
	}
}

func TestNonce_NewNonce_too_small(t *testing.T) {
	tv := make([]byte, MinNonceSize-1)

	nonce, err := NewNonce(tv)

	expectedError := fmt.Sprintf(
		"a nonce must be between %d and %d bytes long; found %d",
		MinNonceSize, MaxNonceSize, len(tv),
	)

	assert.Nil(t, nonce)
	assert.EqualError(t, err, expectedError)
}

func TestNonce_NewNonce_too_big(t *testing.T) {
	tv := make([]byte, MaxNonceSize+1)

	nonce, err := NewNonce(tv)

	expectedError := fmt.Sprintf(
		"a nonce must be between %d and %d bytes long; found %d",
		MinNonceSize, MaxNonceSize, len(tv),
	)

	assert.Nil(t, nonce)
	assert.EqualError(t, err, expectedError)
}

func TestNonce_MarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	nonce, _ := NewNonce([]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef})

	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	expected := append([]byte{0x48}, nonce.Get()...)

	actual, err := nonce.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonce_UnmarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	expected, err := NonceFromHex("abadcafeabadcafe")
	require.Nil(t, err)

	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	data := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	actual := new(Nonce)
	err = actual.UnmarshalCBOR(data)

	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonce_UnmarshalCBOR_bad_cbor(t *testing.T) {
	data := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca}

	actual := Nonce{}
	err := actual.UnmarshalCBOR(data)

	assert.EqualError(t, err, "unexpected EOF")
}

func TestNonce_Validate(t *testing.T) {
	assert := assert.New(t)

	n1 := Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}}
	assert.Nil(n1.Validate())

	// 7 bytes
	n2 := Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe}}
	assert.EqualError(n2.Validate(), "a nonce must be between 8 and 64 bytes long; found 7")

	// 65 bytes
	n3 := Nonce{[]byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde,
	}}
	assert.EqualError(n3.Validate(), "a nonce must be between 8 and 64 bytes long; found 65")
}

func TestNonces_MarshalCBOR_Multiple(t *testing.T) {
	assert := assert.New(t)

	value := Nonces{
		Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}},
		Nonce{[]byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}},
	}

	//82                     # array(2)
	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	expected := []byte{
		0x82, 0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe,
	}

	actual, err := value.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonces_MarshalCBOR_Single(t *testing.T) {
	assert := assert.New(t)

	nonces := Nonces{Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}}}

	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	expected := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	actual, err := nonces.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonces_UnmarshalCBOR_Multiple(t *testing.T) {
	assert := assert.New(t)

	expected := Nonces{
		Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}},
		Nonce{[]byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}},
	}

	//82                     # array(2)
	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	data := []byte{
		0x82, 0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe,
	}

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonces_UnmarshalCBOR_Single(t *testing.T) {
	assert := assert.New(t)

	expected := Nonces{Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}}}

	//   48                  # bytes(8)
	//      deadbeefdeadbeef # "\xDE\xAD\xBE\xEF\xDE\xAD\xBE\xEF"
	data := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonces_UnmarshalCBOR_bad_cbor(t *testing.T) {
	data := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe}

	actual := Nonces{}
	err := actual.UnmarshalCBOR(data)

	assert.EqualError(t, err, "unexpected EOF")
}

func TestNonces_Validate(t *testing.T) {
	assert := assert.New(t)

	ns1 := Nonces{
		Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}},
		Nonce{[]byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}},
	}
	assert.Nil(ns1.Validate())

	ns2 := Nonces{}
	assert.Nil(ns2.Validate())

	ns3 := Nonces{
		Nonce{[]byte{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}},
		Nonce{[]byte{0xab, 0xad, 0xca, 0xfe}},
	}
	assert.EqualError(ns3.Validate(), "a nonce must be between 8 and 64 bytes long; found 4")
}

func TestNonces_UnmarshalCBOR_Repeatedely(t *testing.T) {
	assert := assert.New(t)
	data := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}
	data2 := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	expected := Nonces{Nonce{[]byte{0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}}}

	actual := Nonces{}

	err := actual.UnmarshalCBOR(data)
	assert.Nil(err)
	err = actual.UnmarshalCBOR(data2)
	assert.Nil(err)

	assert.Equal(expected, actual)
}

func TestNonce_MarshalJSON_ok(t *testing.T) {
	tv, err := NewNonce([]byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	})
	require.Nil(t, err)

	expected := []byte(`"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="`)

	actual, err := tv.MarshalJSON()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestNonce_UnmarshalJSON_ok(t *testing.T) {
	expected := []byte{
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}

	tv := []byte(`"3q2+796tvu/erb7v3q2+796tvu/erb7v3q2+796tvu8="`)

	actual := Nonce{}
	err := actual.UnmarshalJSON(tv)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual.Get())
}

func TestNonce_UnmarshalJSON_invalid_json(t *testing.T) {
	tv := []byte(`"unterminated string`)

	actual := Nonce{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestNonce_UnmarshalJSON_invalid_base64(t *testing.T) {
	tv := []byte(`"0"`)

	actual := Nonce{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, "illegal base64 data at input byte 0")
}

func TestNonce_UnmarshalJSON_not_a_string(t *testing.T) {
	tv := []byte(`{ "a": 1 }`)

	actual := Nonce{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, "invalid nonce input map[string]interface {}")
}

func TestNonces_MarshalJSON_one_entry(t *testing.T) {
	tv := Nonces{}

	n1, err := NewNonce(
		[]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	)
	require.Nil(t, err)
	tv.Append(*n1)

	expected := `"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="`

	actual, err := tv.MarshalJSON()

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestNonces_MarshalJSON_two_entries(t *testing.T) {
	tv := Nonces{}

	n1, err := NewNonce(
		[]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	)
	require.Nil(t, err)
	tv.Append(*n1)

	n2, err := NewNonce(
		[]byte{
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
			0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		},
	)
	require.Nil(t, err)
	tv.Append(*n2)

	expected := `[
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		"AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE="
	]`

	actual, err := tv.MarshalJSON()

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestNonces_GetI_ok(t *testing.T) {
	tv := Nonces{}

	expected0 := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	n1, err := NewNonce(expected0)
	require.Nil(t, err)
	tv.Append(*n1)

	expected1 := []byte{
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
	}

	n2, err := NewNonce(expected1)
	require.Nil(t, err)
	tv.Append(*n2)

	actual0 := tv.GetI(0)
	assert.Equal(t, expected0, actual0)

	actual1 := tv.GetI(1)
	assert.Equal(t, expected1, actual1)
}

func TestNonces_GetI_out_of_bounds(t *testing.T) {
	tv := Nonces{}

	for i := -10; i < 10; i++ {
		actual := tv.GetI(i)
		assert.Nil(t, actual)
	}
}

func TestNonceFromHex_bad_hex(t *testing.T) {
	_, err := NonceFromHex("0")

	assert.EqualError(t, err, "encoding/hex: odd length hex string")
}

func TestNonces_UnmarshalJSON_one_entry_ok(t *testing.T) {
	tv := []byte(`"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="`)

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	expected := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, actual.GetI(0))
}

func TestNonces_UnmarshalJSON_two_entries_ok(t *testing.T) {
	tv := []byte(`[
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		"AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE="
	]`)

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	expected0 := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected1 := []byte{
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
	}

	assert.Nil(t, err)
	assert.Equal(t, expected0, actual.GetI(0))
	assert.Equal(t, expected1, actual.GetI(1))
}

func TestNonces_UnmarshalJSON_bad_json(t *testing.T) {
	tv := []byte(`"unterminated string`)

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestNonces_UnmarshalJSON_invalid_base64(t *testing.T) {
	tv := []byte(`"0"`)

	actual := Nonces{}
	err := actual.UnmarshalJSON(tv)

	assert.EqualError(t, err, "illegal base64 data at input byte 0")
}
