// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ueID = MustUEIDFromBytes([]byte{
		0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	})
	oemID      = NewIeeeOemID([3]byte{0xff, 0xff, 0xff})
	cwtID      = BinaryData{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nonceBytes = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	AcmeInc     = "Acme Inc."
	origination = StringOrURI{text: &AcmeInc}
	oemBoot     = true
	debug       = DebugDisabled
	location    = Location{Latitude: 12.34, Longitude: 56.78}
	uptime      = uint(60)
	issuer      = AcmeInc
	subject     = "rr-trap"
	audience    = Audience{origination}
	epoch       = NumericDate(time.Unix(0, 0))

	fatEat = Eat{
		ClaimsSet{
			Nonce:       MustNewNonceFromBytes(nonceBytes),
			UEID:        ueID,
			OemID:       oemID,
			OemBoot:     &oemBoot,
			DebugStatus: &debug,
			Location:    &location,
			Uptime:      &uptime,

			Issuer:     &issuer,
			Subject:    &subject,
			Audience:   &audience,
			Expiration: &epoch,
			NotBefore:  &epoch,
			IssuedAt:   &epoch,
			CwtID:      &cwtID,

			privateClaims: make(map[IntOrString]any),
		},
	}
)

func cborRoundTripper(t *testing.T, tv Eat, expected []byte) {
	data, err := tv.MarshalCBOR()

	t.Logf("CBOR: %x", data)

	assert.Nil(t, err)
	assert.Equal(t, expected, data)

	actual := Eat{}
	err = actual.UnmarshalCBOR(data)

	assert.Nil(t, err)
	assert.EqualValues(t, tv, actual)
}

func jsonRoundTripper(t *testing.T, tv Eat, expected string) {
	data, err := tv.MarshalJSON()

	t.Logf("JSON: '%s'", string(data))

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(data))

	actual := Eat{}
	err = actual.UnmarshalJSON(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func TestEat_Full_RoundtripCBOR(t *testing.T) {
	tv := fatEat
	expected := []byte{
		0xae,                                           // map(14)
		0x01,                                           // . key: 1
		0x69,                                           // . value: tstr(9)
		0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, // . . "Acme Inc"
		0x2e,                                     //       . . "."
		0x02,                                     //       . key: 2
		0x67,                                     //       . value: tstr(7)
		0x72, 0x72, 0x2d, 0x74, 0x72, 0x61, 0x70, //       . . "rr-trap"
		0x03,                                           // . key: 3
		0x69,                                           // . value: tstr(9)
		0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, // . . "Acme Inc"
		0x2e,                               //             . . "."
		0x04,                               //             . key: 4
		0xc1,                               //             . value: tag(1)
		0x00,                               //             . . 0
		0x05,                               //             . key: 5
		0xc1,                               //             . value: tag(1)
		0x00,                               //             . . 0
		0x06,                               //             . key: 6
		0xc1,                               //             . value: tag(1)
		0x00,                               //             . . 0
		0x07,                               //             . key: 7
		0x46,                               //             . value: bstr(6)
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, //
		0x0a,                                           // . key: 10
		0x48,                                           // . value: bstr(8)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //
		0x19, 0x01, 0x00, //                               . key: 256
		0x51,                                           // . value: bstr(17)
		0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, //
		0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, //
		0xef,             //
		0x19, 0x01, 0x02, //                               . key: 258
		0x43,             //                               . value: bstr(3)
		0xff, 0xff, 0xff, //
		0x19, 0x01, 0x05, //                               . key: 261
		0x18, 0x3c, //                                     . value: 60
		0x19, 0x01, 0x06, //                               . key: 262
		0xf5,             //                               . value: true
		0x19, 0x01, 0x07, //                               . key: 263
		0x01,             //                               . value: 1
		0x19, 0x01, 0x08, //                               . key: 264
		0xa2,                                           // . value: map(2)
		0x01,                                           // . . key: 1
		0xfb, 0x40, 0x28, 0xae, 0x14, 0x7a, 0xe1, 0x47, // . . value: 12.34
		0xae,                                           //
		0x02,                                           // . . key: 2
		0xfb, 0x40, 0x4c, 0x63, 0xd7, 0x0a, 0x3d, 0x70, // . . value: 56.78
		0xa4, //
	}

	cborRoundTripper(t, tv, expected)
}

func TestEat_Full_RoundtripJSON(t *testing.T) {
	tv := fatEat
	expected := `
{
	"eat_nonce": "AAAAAAAAAAA",
	"oemid": "____",
	"oemboot": true,
	"dbgstat": "disabled",
	"location": {
		"lat": 12.34,
		"long": 56.78
	},
	"ueid": "Ad6tvu_erb7v3q2-796tvu8",
	"uptime": 60,
	"iss": "Acme Inc.",
	"sub": "rr-trap",
	"aud": "Acme Inc.",
	"exp": 0,
	"nbf": 0,
	"iat": 0,
	"cti": "________"
}`
	// NOTE: cti is not in JSON EAT though
	jsonRoundTripper(t, tv, expected)
}
