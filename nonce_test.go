// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonce_MarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	expected := []byte{0x48, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef}

	nonce, err := NonceFromHex("deadbeefdeadbeef")
	assert.Nil(err)

	actual, err := nonce.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestNonce_UnmarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	expected, err := NonceFromHex("abadcafeabadcafe")
	assert.Nil(err)

	//   48                  # bytes(8)
	//      abadcafeabadcafe # "\xAB\xAD\xCA\xFE\xAB\xAD\xCA\xFE"
	data := []byte{0x48, 0xab, 0xad, 0xca, 0xfe, 0xab, 0xad, 0xca, 0xfe}

	actual := new(Nonce)
	err = actual.UnmarshalCBOR(data)

	assert.Nil(err)
	assert.Equal(expected, actual)
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
