// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedToken_round_trip(t *testing.T) {
	testCases := []struct {
		title        string
		token        NestedToken
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			title: "ok Detached-Submodule-Digest",
			token: NestedToken{
				Type: NestedTokenDigest,
				Data: []byte(`[-16, "3q2-7w"]`),
			},
			expectedCBOR: []byte{
				0x6f,                                           // tstr(15)
				0x5b, 0x2d, 0x31, 0x36, 0x2c, 0x20, 0x22, 0x33, // . "[-16, \"3"
				0x71, 0x32, 0x2d, 0x37, 0x77, 0x22, 0x5d, //       . "q2-7w\"]"
			},
			expectedJSON: `["DIGEST",[-16,"3q2-7w"]]`,
		},
		{
			title: "ok Detached-EAT-Bundle",
			token: NestedToken{
				Type: NestedTokenBundle,
				Data: []byte(`[["DIGEST",[-16,"3q2-7w"]],{"foo":"eyJpc3MiOiJiYXIifQ"}]`),
			},
			// nolint:gocritic
			expectedCBOR: []byte{
				0x78, 0x38, //                                     tstr(56)
				0x5b, 0x5b, 0x22, 0x44, 0x49, 0x47, 0x45, 0x53, // . "[[\"DIGES"
				0x54, 0x22, 0x2c, 0x5b, 0x2d, 0x31, 0x36, 0x2c, // . "T\",[-16,"
				0x22, 0x33, 0x71, 0x32, 0x2d, 0x37, 0x77, 0x22, // . "\"3q2-7w\""
				0x5d, 0x5d, 0x2c, 0x7b, 0x22, 0x66, 0x6f, 0x6f, // . "]],{\"foo"
				0x22, 0x3a, 0x22, 0x65, 0x79, 0x4a, 0x70, 0x63, // . "\":\"eyJpc"
				0x33, 0x4d, 0x69, 0x4f, 0x69, 0x4a, 0x69, 0x59, // . "3MiOiJiY"
				0x58, 0x49, 0x69, 0x66, 0x51, 0x22, 0x7d, 0x5d, // . "XIifQ\"}]"
			},
			expectedJSON: `["BUNDLE",[["DIGEST",[-16,"3q2-7w"]],{"foo":"eyJpc3MiOiJiYXIifQ"}]]`,
		},
		{
			title: "ok JWT",
			token: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte(`"3q2-7w.3q2-7w.3q2-7w"`),
			},
			expectedCBOR: []byte{
				0x76,                                           // tstr(22)
				0x22, 0x33, 0x71, 0x32, 0x2d, 0x37, 0x77, 0x2e, // . "\"3q2-7w."
				0x33, 0x71, 0x32, 0x2d, 0x37, 0x77, 0x2e, 0x33, // . "3q2-7w.3"
				0x71, 0x32, 0x2d, 0x37, 0x77, 0x22, //             . "q2-7w\""
			},
			expectedJSON: `["JWT","3q2-7w.3q2-7w.3q2-7w"]`,
		},
		{
			title: "ok CBOR",
			token: NestedToken{
				Type: NestedTokenCBOR,
				Data: []byte{
					0xd8, 0x3d, // tag(61)
					0x01, //       . 1
				},
			},
			expectedCBOR: []byte{
				0x43,       // bstr(3)
				0xd8, 0x3d, // . tag(61)
				0x01, //       . . 1
			},
			expectedJSON: `["CBOR","2D0B"]`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			encodedCBOR, err := tc.token.MarshalCBOR()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCBOR, encodedCBOR)

			var decodedToken NestedToken
			err = decodedToken.UnmarshalCBOR(encodedCBOR)
			assert.NoError(t, err)
			assert.EqualValues(t, decodedToken, tc.token)

			encodedJSON, err := tc.token.MarshalJSON()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedJSON, string(encodedJSON))

			err = decodedToken.UnmarshalJSON(encodedJSON)
			assert.NoError(t, err)
			assert.Equal(t, decodedToken.Type, tc.token.Type)
			if decodedToken.Type == NestedTokenCBOR {
				assert.Equal(t, decodedToken.Data, tc.token.Data)
			} else {
				assert.JSONEq(t, string(decodedToken.Data), string(tc.token.Data))
			}
		})
	}
}

func TestNestedToken_marshaling_negative(t *testing.T) {
	testCases := []struct {
		title string
		token NestedToken
		err   string
	}{
		{
			title: "bad token type",
			token: NestedToken{
				Type: NestedTokenInvalid,
				Data: []byte("foo"),
			},
			err: "invalid token type: 0",
		},
		{
			title: "malformed inner CBOR",
			token: NestedToken{
				Type: NestedTokenCBOR,
				Data: []byte("foo"),
			},
			err: "inner CBOR: unexpected EOF",
		},
		{
			title: "unexpected inner CBOR tag",
			token: NestedToken{
				Type: NestedTokenCBOR,
				Data: []byte{
					0xd8, 0x2a, // tag(42)
					0x01, //       . 1
				},
			},
			err: "inner CBOR: unexpected tag 42",
		},
		{
			title: "malformed inner JSON",
			token: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte("foo"),
			},
			err: "malformed nested-token",
		},
		{
			title: "mismatched JSON type",
			token: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte(`[-16, "3q2-7w"]`),
			},
			err: `Type is "JWT", but Data appears to contain "DIGEST"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			_, err := tc.token.MarshalCBOR()
			assert.ErrorContains(t, err, tc.err)

			_, err = tc.token.MarshalJSON()
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestNestedToken_UnmarshalCBOR_negative(t *testing.T) {
	testCases := []struct {
		title string
		data  []byte
		err   string
	}{
		{
			title: "malformed CBOR",
			data:  []byte{0x63},
			err:   "unexpected EOF",
		},
		{
			title: "bad CBOR type",
			data: []byte{
				0x01, // 1
			},
			err: "unexpected major type 0",
		},
		{
			title: "unexected tag",
			data: []byte{
				0x43,       // bstr(3)
				0xd8, 0x2a, // . tag(42)
				0x01, //       . . 1
			},
			err: "unexpected tag 42",
		},
		{
			title: "malformed JSON",
			data: []byte{
				0x63,             // tstr(3)
				0x66, 0x6f, 0x6f, // . "foo"
			},
			err: "malformed nested-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var token NestedToken
			err := token.UnmarshalCBOR(tc.data)
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func TestNestedToken_UnmarshalJSON_negative(t *testing.T) {
	testCases := []struct {
		title string
		data  string
		err   string
	}{
		{
			title: "invalid JSON",
			data:  "foo",
			err:   "invalid character 'o' in literal false",
		},
		{
			title: "wrong JSON type",
			data:  `{"1":1}`,
			err:   "cannot unmarshal object into Go value",
		},
		{
			title: "wrong selector type",
			data:  `[1, [-16, "3q2-7w"]]`,
			err:   "selector type: json: cannot unmarshal number into Go value of type string",
		},
		{
			title: "selector mismatch",
			data:  `["DIGEST", "3q2-7w.3q2-7w.3q2-7w"]`,
			err:   `selector type is "DIGEST" but the nested-token appears to be "JWT"`,
		},
		{
			title: "invalid nested-token",
			data:  `["DIGEST", 1]`,
			err:   "malformed nested-token",
		},
		{
			title: "invalid inner CBOR - wrong type",
			data:  `["CBOR", 1]`,
			err:   "json: cannot unmarshal number into Go value of type string",
		},
		{
			title: "invalid inner CBOR - bad base64",
			data:  `["CBOR", "@@@"]`,
			err:   "illegal base64 data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var token NestedToken
			err := token.UnmarshalJSON([]byte(tc.data))
			assert.ErrorContains(t, err, tc.err)
		})
	}
}

func Test_NewNestedToken(t *testing.T) {
	testCases := []struct {
		title    string
		typ      string
		data     []byte
		expected NestedToken
		err      string
	}{
		{
			title: "ok CBOR",
			typ:   "CBOR",
			data: []byte{
				0xd8, 0x3d, // tag(61)
				0x01, //       . 1
			},
			expected: NestedToken{
				Type: NestedTokenCBOR,
				Data: []byte{0xd8, 0x3d, 0x01},
			},
		},
		{
			title: "ok JWT",
			typ:   "JWT",
			data:  []byte(`"3q2-7w.3q2-7w.3q2-7w"`),
			expected: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte(`"3q2-7w.3q2-7w.3q2-7w"`),
			},
		},
		{
			title: "err invalid type",
			typ:   "FOO",
			data:  []byte(`"3q2-7w.3q2-7w.3q2-7w"`),
			err:   `unknown nested token type "FOO"`,
		},
		{
			title: "err malformed nested-token",
			typ:   "DIGEST",
			data:  []byte("@@@"),
			err:   "malformed nested-token",
		},
		{
			title: "err mismatched type and data",
			typ:   "DIGEST",
			data:  []byte(`"3q2-7w.3q2-7w.3q2-7w"`),
			err:   `provided type is "DIGEST", but data appears to contain "JWT"`,
		},
		{
			title: "err unexpected CBOR tag",
			typ:   "CBOR",
			data: []byte{
				0xd8, 0x2a, // . tag(42)
				0x01, //       . . 1
			},
			err: "cbor: unexpected tag 42",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			token, err := NewNestedToken(tc.typ, tc.data)

			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, *token)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestNestedToken_UnwrapDigest(t *testing.T) {
	testCases := []struct {
		title    string
		token    NestedToken
		expected DetachedSubmodDigest
		err      string
	}{
		{
			title: "ok",
			token: NestedToken{
				Type: NestedTokenDigest,
				Data: []byte(`[-16, "3q2-796tvu_erb7v3q2-796tvu_erb7v3q2-796tvu8"]`),
			},
			expected: DetachedSubmodDigest{
				Algorithm: MustHashAlgorithmFromInt(Sha256),
				Value: []byte{
					0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
					0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
					0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
					0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
				},
			},
		},
		{
			title: "err wrong type",
			token: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte(`[-16, "3q2-7w.3q2-7w.3q2-7w"]`),
			},
			err: "token contains JWT, not DIGEST",
		},
		{
			title: "err bad digest",
			token: NestedToken{
				Type: NestedTokenDigest,
				Data: []byte(`[-16, "3q2-7w"]`),
			},
			err: "length mismatch for hash algorithm sha-256",
		},
		{
			title: "err malformed JSON",
			token: NestedToken{
				Type: NestedTokenDigest,
				Data: []byte(`@@@`),
			},
			err: "invalid character '@'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			unwrapped, err := tc.token.UnwrapDigest()

			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, &tc.expected, unwrapped)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestNestedToken_UnwrapBundle(t *testing.T) {
	testCases := []struct {
		title    string
		token    NestedToken
		expected DetachedEatBundle
		err      string
	}{
		{
			title: "ok",
			token: NestedToken{
				Type: NestedTokenBundle,
				Data: []byte(`[["DIGEST",[-16,"3q2-7w"]],{"foo":"eyJpc3MiOiJiYXIifQ"}]`),
			},
			expected: DetachedEatBundle{
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
		},
		{
			title: "err wrong type",
			token: NestedToken{
				Type: NestedTokenJWT,
				Data: []byte(`[-16, "3q2-7w.3q2-7w.3q2-7w"]`),
			},
			err: "token contains JWT, not BUNDLE",
		},
		{
			title: "err malformed JSON",
			token: NestedToken{
				Type: NestedTokenBundle,
				Data: []byte(`@@@`),
			},
			err: "invalid character '@'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			unwrapped, err := tc.token.UnwrapBundle()

			if tc.err == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, &tc.expected, unwrapped)
			} else {
				assert.ErrorContains(t, err, tc.err)
			}
		})
	}
}
