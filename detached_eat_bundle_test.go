// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetachedEatBundle_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		bundle       DetachedEatBundle
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok",
			bundle: DetachedEatBundle{
				token: &NestedToken{
					Type: NestedTokenDigest,
					Data: []byte(`[-16,"3q2-7w"]`),
				},
				claims: map[string]*ClaimsSet{
					"foo": &ClaimsSet{
						Issuer:        Ptr("bar"),
						privateClaims: make(map[IntOrString]any),
					},
				},
			},
			expectedCBOR: []byte{
				0x82,                                           // array(2)
				0x6e,                                           // . [0]tstr(14)
				0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x22, 0x33, 0x71, // . . "[-16,\"3"
				0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //             . . "q2-7w\"]"
				0xa1,             //                               . [1]map(1)
				0x63,             //                               . . key: tstr(3)
				0x66, 0x6f, 0x6f, //                               . . . "foo"
				0x46,             //                               . . value: bstr(6)
				0xa1,             //                               . . . map(1) [ClaimsSet]
				0x01,             //                               . . . . key: 1
				0x63,             //                               . . . . value: tstr(3)
				0x62, 0x61, 0x72, //                               . . . . . "bar"
			},
			expectedJSON: `[["DIGEST",[-16,"3q2-7w"]],{"foo":"eyJpc3MiOiJiYXIifQ"}]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.bundle.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedBundle DetachedEatBundle
			err = decodedBundle.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.bundle, decodedBundle)

			encodedJSON, err := tc.bundle.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tc.expectedJSON, string(encodedJSON))

			err = decodedBundle.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.EqualValues(t, tc.bundle, decodedBundle)
		})
	}
}

func TestDetachedEatBundle_marshaling_negative(t *testing.T) {
	testCases := []struct {
		title  string
		bundle DetachedEatBundle
		err    string
	}{
		{
			title:  "err missing main token",
			bundle: DetachedEatBundle{},
			err:    "nil main token",
		},
		{
			title: "err missing claims sets",
			bundle: DetachedEatBundle{
				token: &NestedToken{
					Type: NestedTokenDigest,
					Data: []byte(`[-16,"3q2-7w"]`),
				},
			},
			err: "no detached claims sets",
		},
		{
			title: "err bad main token",
			bundle: DetachedEatBundle{
				token: &NestedToken{
					Type: NestedTokenJWT,
					Data: []byte(`[-16,"3q2-7w"]`),
				},
				claims: map[string]*ClaimsSet{
					"foo": &ClaimsSet{
						Issuer:        Ptr("bar"),
						privateClaims: make(map[IntOrString]any),
					},
				},
			},
			err: `Type is "JWT", but Data appears to contain "DIGEST"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := tc.bundle.MarshalCBOR()
			assert.ErrorContains(t, err, tc.err)

			_, err = tc.bundle.MarshalJSON()
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestDetachedEatBundle_UnmarshalCBOR_negative(t *testing.T) {
	testCases := []struct {
		title string
		data  []byte
		err   string
	}{
		{
			title: "err empty",
			data:  []byte{},
			err:   "truncated input",
		},
		{
			title: "err wrong outer type",
			data: []byte{
				0x01, // 1
			},
			err: "invalid major type 0 for DetachedEatBundle",
		},
		{
			title: "err wrong main token type",
			data: []byte{
				0x82, // array(2)
				0x01, // . [0]1
				0x02, // . [1]2
			},
			err: "main token: unexpected major type 0",
		},
		{
			title: "err wrong claims sets type",
			data: []byte{
				0x82,                                           // array(2)
				0x6e,                                           // . [0]tstr(14)
				0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x22, 0x33, 0x71, // . . "[-16,\"3"
				0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //             . . "q2-7w\"]"
				0x01, //                                           . [1]1
			},
			err: "detached claims sets: cbor: cannot unmarshal positive integer into Go value",
		},
		{
			title: "err malformed claims set",
			data: []byte{
				0x82,                                           // array(2)
				0x6e,                                           // . [0]tstr(14)
				0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x22, 0x33, 0x71, // . . "[-16,\"3"
				0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //             . . "q2-7w\"]"
				0xa1,             //                               . [1]map(1)
				0x63,             //                               . . key: tstr(3)
				0x66, 0x6f, 0x6f, //                               . . . "foo"
				0x41, //                                           . . value: bstr(1)
				0x01,
			},
			err: `claims set "foo": cbor: cannot unmarshal positive integer into Go value`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var bundle DetachedEatBundle
			err := bundle.UnmarshalCBOR(tc.data)
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestDetachedEatBundle_UnmarshalJSON_negative(t *testing.T) {
	testCases := []struct {
		title string
		data  string
		err   string
	}{
		{
			title: "err malformed JSON",
			data:  `@@@`,
			err:   "invalid character '@'",
		},
		{
			title: "err wrong outer type",
			data:  `{"foo":1}`,
			err:   "json: cannot unmarshal object into Go value",
		},
		{
			title: "err wrong main token type",
			data:  `[1, 2]`,
			err:   "main token: json: cannot unmarshal number into Go value",
		},
		{
			title: "err wrong claims sets type",
			data:  `[["DIGEST",[-16,"3q2-7w"]], 2]`,
			err:   "detached claims sets: json: cannot unmarshal number into Go value",
		},
		{
			title: "err malformed claims set base64",
			data:  `[["DIGEST",[-16,"3q2-7w"]], {"foo": "@@@"}]`,
			err:   `claims set "foo" base64: illegal base64 data at input byte 0`,
		},
		{
			title: "err malformed claims set",
			data:  `[["DIGEST",[-16,"3q2-7w"]], {"foo": "aaa"}]`,
			err:   `claims set "foo": invalid character 'i'`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var bundle DetachedEatBundle
			err := bundle.UnmarshalJSON([]byte(tc.data))
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestDetachedEatBundle_tagged(t *testing.T) {
	bundle := DetachedEatBundle{
		token: &NestedToken{
			Type: NestedTokenDigest,
			Data: []byte(`[-16,"3q2-7w"]`),
		},
		claims: map[string]*ClaimsSet{
			"foo": &ClaimsSet{
				Issuer:        Ptr("bar"),
				privateClaims: make(map[IntOrString]any),
			},
		},
	}
	// nolint:gocritic
	expectedCBOR := []byte{
		0xd9, 0x02, 0x5a, //                               tag(602)
		0x82,                                           // . array(2)
		0x6e,                                           // . . [0]tstr(14)
		0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x22, 0x33, 0x71, // . . . "[-16,\"3"
		0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //             . . . "q2-7w\"]"
		0xa1,             //                               . . [1]map(1)
		0x63,             //                               . . . key: tstr(3)
		0x66, 0x6f, 0x6f, //                               . . . . "foo"
		0x46,             //                               . . . value: bstr(6)
		0xa1,             //                               . . . . map(1) [ClaimsSet]
		0x01,             //                               . . . . . key: 1
		0x63,             //                               . . . . . value: tstr(3)
		0x62, 0x61, 0x72, //                               . . . . . . "bar"
	}

	encodedCBOR, err := bundle.ToTaggedMessageCBOR()
	assert.NoError(t, err)
	assert.Equal(t, expectedCBOR, encodedCBOR)

	var decodedBundle DetachedEatBundle
	err = decodedBundle.UnmarshalCBOR(encodedCBOR)
	assert.NoError(t, err)
	assert.EqualValues(t, bundle, decodedBundle)
}

func TestDetachedEatBundle_misc(t *testing.T) {
	token := NestedToken{
		Type: NestedTokenDigest,
		Data: []byte(`[-16,"3q2-7w"]`),
	}

	claimsSets := map[string]*ClaimsSet{
		"foo": &ClaimsSet{
			Issuer:        Ptr("bar"),
			privateClaims: make(map[IntOrString]any),
		},
	}

	bundle := NewDetachedEatBundle(&token).
		SetDetachedClaimSets(claimsSets)

	assert.Equal(t, &token, bundle.MainToken())
	assert.Equal(t, claimsSets, bundle.DetachedClaimSets())

	assert.Panics(t, func() {
		NewDetachedEatBundle(nil)
	})

	assert.Panics(t, func() {
		bundle.SetDetachedClaimSets(nil)
	})
}

func TestDetachedEatBundle_RFC_examples_smoke_test(t *testing.T) {
	testCases := []struct {
		title string
		input []byte
	}{
		{
			title: "bundle CBOR",
			input: rfcExampleDetachedEatBundleCBOR,
		},
		{
			title: "bundle JSON",
			input: rfcExampleDetachedEatBundleJSON,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var err error
			var bundle DetachedEatBundle

			if strings.HasSuffix(tc.title, "CBOR") {
				err = bundle.UnmarshalCBOR(tc.input)
			} else {
				err = bundle.UnmarshalJSON(tc.input)
			}

			assert.NoError(t, err)
		})
	}
}
