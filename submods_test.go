// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestSubmod_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		submod       Submod
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok nested token",
			submod: Submod{&NestedToken{
				Type: NestedTokenDigest,
				Data: []byte(`[-16,"3q2-7w"]`),
			}},
			expectedCBOR: []byte{
				0x6e,                                           // tstr(14)
				0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x22, 0x33, 0x71, // . "[-16,\"3q"
				0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //             . "2-7w\"]"
			},
			expectedJSON: `["DIGEST",[-16,"3q2-7w"]]`,
		},
		{
			title: "ok claims set",
			submod: Submod{&ClaimsSet{
				BootCount:     Ptr(uint(1)),
				privateClaims: make(map[IntOrString]any),
			}},
			expectedCBOR: []byte{
				0xa1,             //             map(1)
				0x19, 0x01, 0x0b, // . key: 267
				0x01, //             . value: 1
			},
			expectedJSON: `{"bootcount":1}`,
		},
		{
			title:  "ok detached submod digest",
			submod: Submod{MustNewDetachedSubmodDigestIntAlg(Sha256, testDigestBytes)},
			expectedCBOR: []byte{
				0x82,       // array(2)
				0x2f,       // . [0]-16 [sha-256]
				0x58, 0x20, // . [1]bstr(32)
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
			},
			expectedJSON: `["DIGEST",[-16, "AAECAwQFBgcAAQIDBAUGBwABAgMEBQYHAAECAwQFBgc"]]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.submod.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedSubmod Submod
			err = decodedSubmod.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.submod, decodedSubmod)

			encodedJSON, err := tc.submod.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = decodedSubmod.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)

			if tc.title == "ok detached submod digest" {
				// As JSON submod encoding does not allow a detanched submod digest,
				// it gets wrapped in a nested token.
				expectedSubmod := Submod{&NestedToken{
					Type: NestedTokenDigest,
					Data: []byte(`[-16,"AAECAwQFBgcAAQIDBAUGBwABAgMEBQYHAAECAwQFBgc"]`),
				}}
				assert.EqualValues(t, expectedSubmod, decodedSubmod)
			} else {
				assert.EqualValues(t, tc.submod, decodedSubmod)
			}
		})
	}
}

func TestSubmod_UnmarshalJSON_empty(t *testing.T) {
	var submod Submod
	err := submod.UnmarshalJSON([]byte{})
	assert.ErrorContains(t, err, "truncated input")
}

func TestSubmods_add_get(t *testing.T) {
	submods := NewSubmods().
		AddClaimsSet("foo", &ClaimsSet{
			BootCount:     Ptr(uint(1)),
			privateClaims: make(map[IntOrString]any),
		}).
		AddDigest("bar", MustNewDetachedSubmodDigestIntAlg(Sha256, testDigestBytes)).
		AddNestedToken("qux", &NestedToken{
			Type: NestedTokenDigest,
			Data: []byte(`[-16,"3q2-7w"]`),
		})

	claimsSet, err := submods.GetClaimsSet("foo")
	assert.NoError(t, err)
	assert.EqualValues(t, Ptr(uint(1)), claimsSet.BootCount)

	_, err = submods.GetClaimsSet("bar")
	assert.ErrorContains(t, err, `submod "bar" is not a claims set`)

	_, err = submods.GetClaimsSet("zot")
	assert.ErrorContains(t, err, `no submod named "zot"`)

	digest, err := submods.GetDigest("bar")
	assert.NoError(t, err)
	assert.EqualValues(t, HashAlgorithmSha256, digest.Algorithm)

	_, err = submods.GetDigest("foo")
	assert.ErrorContains(t, err, `submod "foo" is not a detached submod digest`)

	_, err = submods.GetDigest("zot")
	assert.ErrorContains(t, err, `no submod named "zot"`)

	token, err := submods.GetNestedToken("qux")
	assert.NoError(t, err)
	assert.EqualValues(t, NestedTokenDigest, token.Type)

	_, err = submods.GetNestedToken("foo")
	assert.ErrorContains(t, err, `submod "foo" is not a nested token`)

	_, err = submods.GetNestedToken("zot")
	assert.ErrorContains(t, err, `no submod named "zot"`)

	assert.Panics(t, func() { submods.AddClaimsSet("zot", nil) })
	assert.Panics(t, func() { submods.AddDigest("zot", nil) })
	assert.Panics(t, func() { submods.AddNestedToken("zot", nil) })
}

func TestSubmods_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		submods      Submods
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok claims set",
			submods: *NewSubmods().AddClaimsSet("foo", &ClaimsSet{
				BootCount:     Ptr(uint(1)),
				privateClaims: make(map[IntOrString]any),
			}),
			expectedCBOR: []byte{
				0xa1,             // map(1)
				0x63,             // . key: tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
				0xa1,             // . value: map(1) [ClaimSet]
				0x19, 0x01, 0x0b, // . . key: 267
				0x01, //             . . value: 1
			},
			expectedJSON: `{"foo": {"bootcount":1}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := cbor.Marshal(tc.submods)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedSubmods Submods
			err = cbor.Unmarshal(encodedCBOR, &decodedSubmods)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.submods, decodedSubmods)

			encodedJSON, err := json.Marshal(tc.submods)
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = json.Unmarshal(encodedJSON, &decodedSubmods)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.submods, decodedSubmods)
		})
	}
}
