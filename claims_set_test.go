// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimsSet_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		claims       ClaimsSet
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok minimal",
			claims: ClaimsSet{
				BootCount:     Ptr(uint(1)),
				privateClaims: make(map[IntOrString]any),
			},
			expectedCBOR: []byte{
				0xa1,             // map(1)
				0x19, 0x01, 0x0b, // . key: 267
				0x01, //             . value: 1
			},
			expectedJSON: `{"bootcount":1}`,
		},
		{
			title: "ok private claims only",
			claims: ClaimsSet{
				privateClaims: map[IntOrString]any{
					IntOrStringFromInt(-1): "foo",
				},
			},
			expectedCBOR: []byte{
				0xa1,             // map(1)
				0x20,             // . key: -1
				0x63,             // . value: tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
			},
			expectedJSON: `{"-1":"foo"}`,
		},
		{
			title: "ok private and standard claims",
			claims: ClaimsSet{
				BootCount: Ptr(uint(1)),
				privateClaims: map[IntOrString]any{
					IntOrStringFromInt(-1): "foo",
				},
			},
			expectedCBOR: []byte{
				0xa2,             // map(2)
				0x20,             // . key: -1
				0x63,             // . value: tstr(3)
				0x66, 0x6f, 0x6f, // . . "foo"
				0x19, 0x01, 0x0b, // . key: 267
				0x01, //             . value: 1
			},
			expectedJSON: `{"-1":"foo","bootcount":1}`,
		},
		{
			title: "ok submods",
			claims: ClaimsSet{
				Submods: NewSubmods().AddClaimsSet("foo", &ClaimsSet{
					BootCount:     Ptr(uint(1)),
					privateClaims: make(map[IntOrString]any),
				}),
				privateClaims: make(map[IntOrString]any),
			},
			expectedCBOR: []byte{
				0xa1,             // map(1)
				0x19, 0x01, 0x0a, // . key: 266
				0xa1,             // . value: map(1) [Submods]
				0x63,             // . . key: tstr(3)
				0x66, 0x6f, 0x6f, // . . . "foo"
				0xa1,             // . . value: map(1) [ClaimSet]
				0x19, 0x01, 0x0b, // . . . key: 267
				0x01, //             . . . value: 1
			},
			expectedJSON: `{"submods": {"foo": {"bootcount":1}}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.claims.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedClaims ClaimsSet
			err = decodedClaims.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.claims, decodedClaims)

			encodedJSON, err := tc.claims.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = decodedClaims.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.claims, decodedClaims)
		})
	}
}

func TestClaimsSet_private(t *testing.T) {
	claimsSet := NewClaimsSet()

	err := claimsSet.SetPrivate(IntOrStringFromInt(258), 1)
	assert.ErrorContains(t, err, `key "258" is a standard claim`)
	val := claimsSet.GetPrivate(IntOrStringFromString("foo"))
	assert.Nil(t, val)

	err = claimsSet.SetPrivate(IntOrStringFromString("oemid"), 1)
	assert.ErrorContains(t, err, `key "oemid" is a standard claim`)

	err = claimsSet.SetPrivate(IntOrStringFromString("foo"), 1)
	assert.NoError(t, err)
	val = claimsSet.GetPrivate(IntOrStringFromString("foo"))
	assert.EqualValues(t, 1, val)
}

func TestClaimsSet_RFC_examples_smoke_test(t *testing.T) {
	testCases := []struct {
		title string
		input []byte
	}{
		{
			title: "simple CBOR",
			input: rfcExampleClaimsSetSimpleTeeCBOR,
		},
		{
			title: "submods CBOR",
			input: rfcExampleClaimsSetBoardAndDeviceSubmodsCBOR,
		},
		{
			title: "hw block CBOR",
			input: rfcExampleClaimsSetAttestationHwBlockCBOR,
		},
		{
			title: "key store CBOR",
			input: rfcExampleClaimsSetKeyStoreAttestCBOR,
		},
		{
			title: "sw meas CBOR",
			input: rfcExampleClaimsSetSwMeasuresIotDeviceCBOR,
		},
		{
			title: "attest result JSON",
			input: rfcExampleClaimsSetAttestResultJSON,
		},
		{
			title: "submods JSON",
			input: rfcExampleClaimsSetTokenWithSubmodsJSON,
		},
		{
			title: "detached digest CBOR",
			input: rfcExampleClaimsSetDetachedDigestSubmodCBOR,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var err error
			var claimsSet ClaimsSet

			if strings.HasSuffix(tc.title, "CBOR") {
				err = claimsSet.UnmarshalCBOR(tc.input)
			} else {
				err = claimsSet.UnmarshalJSON(tc.input)
			}

			assert.NoError(t, err)
		})
	}
}
