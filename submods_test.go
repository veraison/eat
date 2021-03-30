// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmods_Add_OK(t *testing.T) {
	var s Submods

	emptyEatToken := []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}

	err := s.Add("eat-claims", Eat{})
	assert.Nil(t, err)

	err = s.Add("eat-token", emptyEatToken)
	assert.Nil(t, err)

	assert.Equal(t, Eat{}, s.Get("eat-claims"))
	assert.Equal(t, emptyEatToken, s.Get("eat-token"))
}

func TestSubmods_Add_FAIL(t *testing.T) {
	var s Submods

	justTagsNoEatToken := []byte{0xd8, 0x3d, 0xd2}

	err := s.Add("eat-token", justTagsNoEatToken)
	assert.EqualError(t, err, "not enough bytes")

	noTagsJustRandomStuff := []byte{0x00, 0x01, 0x02, 0x03, 0x04}

	err = s.Add("eat-token", noTagsJustRandomStuff)
	assert.EqualError(t, err, "CWT and COSE Sign1 tags not found")

	badSubmodType := 12.34

	err = s.Add("eat-token", badSubmodType)
	assert.EqualError(t, err, "submod must be Eat or []byte")
}

func TestSubmods_JSONMarshal_Simple(t *testing.T) {
	var s Submods

	require.Nil(t, s.Add("0", Eat{Nonce: &Nonces{nonce{nonceBytes}}}))
	require.Nil(t, s.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	expected := `{
		"0": {
			"nonce": "3q2+7w=="
		},
		"xyz": "2D3SQaA="
	}`

	actual, err := json.Marshal(s)
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestSubmods_JSONMarshal_Nested(t *testing.T) {
	var inner Submods
	require.Nil(t, inner.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	eat := Eat{Submods: &inner}

	var outer Submods
	require.Nil(t, outer.Add("0", eat))

	expected := `{
		"0": {
			"submods": {
				"xyz": "2D3SQaA="
			}
		}
	}`

	actual, err := json.Marshal(outer)
	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestSubmods_JSONUnmarshal_Simple(t *testing.T) {
	tv := []byte(`{
		"0": {
			"nonce": "3q2+7w=="
		},
		"xyz": "2D3SQaA="
	}`)

	var s Submods

	err := json.Unmarshal(tv, &s)
	assert.Nil(t, err)

	assert.Equal(t, Eat{Nonce: &Nonces{nonce{nonceBytes}}}, s.Get("0"))
	assert.Equal(t, []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}, s.Get("xyz"))
}

func TestSubmods_JSONUnmarshal_Nested(t *testing.T) {
	tv := []byte(`{
		"0": {
			"submods": {
				"xyz": "2D3SQaA="
			}
		}
	}`)

	var outer Submods

	err := json.Unmarshal(tv, &outer)
	assert.Nil(t, err)

	var inner Submods
	require.Nil(t, inner.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	assert.Equal(t, Eat{Submods: &inner}, outer.Get("0"))
}

func TestSubmods_CBORMarshal_Simple(t *testing.T) {
	var s Submods

	require.Nil(t, s.Add("0", Eat{Nonce: &Nonces{nonce{nonceBytes}}}))
	require.Nil(t, s.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	// echo "{\"0\": {10: h'deadbeef'}, \"xyz\": h'd83dd241a0'}" | diag2cbor.rb | xxd -i
	expected := []byte{
		0xa2, 0x61, 0x30, 0xa1, 0x0a, 0x44, 0xde, 0xad, 0xbe, 0xef, 0x63,
		0x78, 0x79, 0x7a, 0x45, 0xd8, 0x3d, 0xd2, 0x41, 0xa0,
	}

	actual, err := em.Marshal(s)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubmods_CBORMarshal_Nested(t *testing.T) {
	var inner Submods
	require.Nil(t, inner.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	eat := Eat{Submods: &inner}

	var outer Submods
	require.Nil(t, outer.Add("0", eat))

	// echo "{\"0\": {20: {\"xyz\": h'd83dd241a0'}}}" | diag2cbor.rb | xxd -i
	expected := []byte{
		0xa1, 0x61, 0x30, 0xa1, 0x14, 0xa1, 0x63, 0x78, 0x79, 0x7a, 0x45,
		0xd8, 0x3d, 0xd2, 0x41, 0xa0,
	}

	actual, err := em.Marshal(outer)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubmods_CBORUnmarshal_Simple(t *testing.T) {
	// echo "{\"0\": {10: h'deadbeef'}, \"xyz\": h'd83dd241a0'}" | diag2cbor.rb | xxd -i
	tv := []byte{
		0xa2, 0x61, 0x30, 0xa1, 0x0a, 0x44, 0xde, 0xad, 0xbe, 0xef, 0x63,
		0x78, 0x79, 0x7a, 0x45, 0xd8, 0x3d, 0xd2, 0x41, 0xa0,
	}

	var s Submods

	err := dm.Unmarshal(tv, &s)
	assert.Nil(t, err)

	assert.Equal(t, Eat{Nonce: &Nonces{nonce{nonceBytes}}}, s.Get("0"))
	assert.Equal(t, []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}, s.Get("xyz"))
}

func TestSubmods_CBORUnmarshal_SimpleWithNegativeKey(t *testing.T) {
	// echo "{\"-1\": {10: h'deadbeef'}, \"xyz\": h'd83dd241a0'}" | diag2cbor.rb | xxd -i
	tv := []byte{
		0xa2, 0x62, 0x2d, 0x31, 0xa1, 0x0a, 0x44, 0xde, 0xad, 0xbe, 0xef,
		0x63, 0x78, 0x79, 0x7a, 0x45, 0xd8, 0x3d, 0xd2, 0x41, 0xa0,
	}

	var s Submods

	err := dm.Unmarshal(tv, &s)
	assert.Nil(t, err)

	assert.Equal(t, Eat{Nonce: &Nonces{nonce{nonceBytes}}}, s.Get("-1"))
	assert.Equal(t, []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}, s.Get("xyz"))
}

func TestSubmods_CBORUnmarshal_Nested(t *testing.T) {
	// echo "{ \"0\": { 20: { \"xyz\": h'd83dd241a0' } } }" | diag2cbor.rb | xxd -i
	tv := []byte{
		0xa1, 0x61, 0x30, 0xa1, 0x14, 0xa1, 0x63, 0x78, 0x79, 0x7a, 0x45,
		0xd8, 0x3d, 0xd2, 0x41, 0xa0,
	}

	var outer Submods

	err := dm.Unmarshal(tv, &outer)
	assert.Nil(t, err)

	var inner Submods
	require.Nil(t, inner.Add("xyz", []byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}))

	assert.Equal(t, Eat{Submods: &inner}, outer.Get("0"))
}
