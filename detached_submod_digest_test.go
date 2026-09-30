// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDigestBytes = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	}
)

func TestHashAlgorithm_conversion(t *testing.T) {
	alg := MustHashAlgorithmFromInt(Sha256)
	assert.True(t, alg.IsInt())
	assert.False(t, alg.IsString())
	assert.Equal(t, Sha256, alg.Int())
	assert.Equal(t, "sha-256", alg.String())
}

func TestHashAlgorithmFromString(t *testing.T) {
	testCases := []struct {
		title    string
		text     string
		expected HashAlgorithm
	}{
		{
			title:    "known string",
			text:     "sha-256",
			expected: HashAlgorithm{Sha256},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			alg, err := HashAlgorithmFromString(tc.text)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.expected, alg)
		})
	}
}

func TestHashAlgorithmFromAny(t *testing.T) {
	testCases := []struct {
		title    string
		value    any
		expected HashAlgorithm
		err      string
	}{
		{
			title:    "int",
			value:    Sha256,
			expected: HashAlgorithm{Sha256},
		},
		{
			title:    "int64",
			value:    int64(-16),
			expected: HashAlgorithm{-16},
		},
		{
			title:    "float64",
			value:    -16.0,
			expected: HashAlgorithm{-16},
		},
		{
			title:    "string",
			value:    "sha-256",
			expected: HashAlgorithm{-16},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			alg, err := HashAlgorithmFromAny(tc.value)
			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, alg)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestHashAlgorithm_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		value        HashAlgorithm
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title:        "known int",
			value:        MustHashAlgorithmFromInt(Sha256),
			expectedCBOR: []byte{0x2f},
			expectedJSON: `-16`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var alg HashAlgorithm
			err = dm.Unmarshal(bytes, &alg)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, alg)

			bytes, err = json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &alg)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, alg)
		})
	}

}

func TestHashAlgorithm_UnmarshalCBOR_bad(t *testing.T) {
	var alg HashAlgorithm
	err := alg.UnmarshalCBOR([]byte{})
	assert.ErrorContains(t, err, "buffer too short")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "19"))
	assert.ErrorContains(t, err, "unexpected EOF")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "64ffffffff"))
	assert.ErrorContains(t, err, "invalid UTF-8 string")

	err = alg.UnmarshalCBOR(MustHexDecode(t, "f4"))
	assert.ErrorContains(t, err, "unexpected CBOR major type")
}

func TestHashAlgorithm_UnmarshalJSON_bad(t *testing.T) {
	var alg HashAlgorithm
	err := alg.UnmarshalJSON([]byte{})
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	err = alg.UnmarshalJSON([]byte("true"))
	assert.ErrorContains(t, err, "unexpected algorithm value: true(bool)")
}

func TestDetachedSubmodDigestFromString(t *testing.T) {
	digest, err := DetachedSubmodDigestFromString("sha-256;AAECAwQFBgcAAQIDBAUGBwABAgMEBQYHAAECAwQFBgc")
	assert.NoError(t, err)
	assert.EqualValues(t, MustNewDetachedSubmodDigestIntAlg(Sha256, testDigestBytes), digest)

	_, err = DetachedSubmodDigestFromString("sha-256;AQID")
	assert.ErrorContains(t, err, "length mismatch for hash algorithm sha-256")

	_, err = DetachedSubmodDigestFromString("sha-256")
	assert.ErrorContains(t, err, `expected exactly two ;-separated parts, got "sha-256"`)

	_, err = DetachedSubmodDigestFromString("sha-256;@@@")
	assert.ErrorContains(t, err, "val: illegal base64 data")
}

func TestDetachedSubmodDigest_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		value        *DetachedSubmodDigest
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "int",
			value: MustNewDetachedSubmodDigestIntAlg(Sha256, testDigestBytes),
			expectedCBOR: []byte{
				0x82,       // array(2)
				0x2f,       // . [0]-16 [sha-256]
				0x58, 0x20, // . [1]bstr(32)
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
			},
			expectedJSON: `[-16, "AAECAwQFBgcAAQIDBAUGBwABAgMEBQYHAAECAwQFBgc"]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			bytes, err := em.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, bytes)

			var digest DetachedSubmodDigest
			err = dm.Unmarshal(bytes, &digest)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, &digest)

			bytes, err = json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(bytes))

			err = json.Unmarshal(bytes, &digest)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.value, &digest)
		})
	}

}

func TestDetachedSubmodDigest_UnmarshalJSON_bad(t *testing.T) {
	var digest DetachedSubmodDigest

	err := digest.UnmarshalJSON([]byte(""))
	assert.ErrorContains(t, err, "unexpected end of JSON input")

	err = digest.UnmarshalJSON([]byte(`{"alg": -16, "value": "foo"}`))
	assert.ErrorContains(t, err, "cannot unmarshal object into Go value of type []interface {}")

	err = digest.UnmarshalJSON([]byte("[1, 2, 3]"))
	assert.ErrorContains(t, err, "expected array with two elements")

	err = digest.UnmarshalJSON([]byte(`[true, "bar"]`))
	assert.ErrorContains(t, err, "invalid hash algorithm: true(bool)")

	err = digest.UnmarshalJSON([]byte(`[-16, "@@@"]`))
	assert.ErrorContains(t, err, "val: illegal base64 data")
}
