// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ueID = UEID{
		0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}
	oemID      = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nonceBytes = []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	AcmeInc     = "Acme Inc."
	origination = StringOrURI{text: &AcmeInc}
	oemBoot     = true
	debug       = Debug(DebugDisabled)
	location    = Location{Latitude: 12.34, Longitude: 56.78}
	uptime      = uint(60)
	submods     = Submods{
		"eat-claims": Submod{Eat{}},
		"eat-token":  Submod{[]byte{0xd8, 0x3d, 0xd2, 0x41, 0xa0}},
	}

	issuer   = AcmeInc
	subject  = "rr-trap"
	audience = Audience{origination}
	epoch    = NumericDate(time.Unix(0, 0))

	fatEat = Eat{
		Nonce:       &Nonce{nonce{nonceBytes}},
		UEID:        &ueID,
		OemID:       &oemID,
		OemBoot:     &oemBoot,
		DebugStatus: &debug,
		Location:    &location,
		Uptime:      &uptime,

		CWTClaims: CWTClaims{
			Issuer:     &issuer,
			Subject:    &subject,
			Audience:   &audience,
			Expiration: &epoch,
			NotBefore:  &epoch,
			IssuedAt:   &epoch,
			CwtID:      &oemID,
		},
	}

	justEatSubmods = Eat{
		Submods: &submods,
	}
)

func cborRoundTripper(t *testing.T, tv Eat, expected []byte) {
	data, err := tv.ToCBOR()

	t.Logf("CBOR: %x", data)

	assert.Nil(t, err)
	assert.Equal(t, expected, data)

	actual := Eat{}
	err = actual.FromCBOR(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func jsonRoundTripper(t *testing.T, tv Eat, expected string) {
	data, err := tv.ToJSON()

	t.Logf("JSON: '%s'", string(data))

	assert.Nil(t, err)
	assert.JSONEq(t, expected, string(data))

	actual := Eat{}
	err = actual.FromJSON(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func TestEat_Full_RoundtripCBOR(t *testing.T) {
	tv := fatEat
	/*
		ae                                      # map(14)
		   01                                   # unsigned(1)
		   69                                   # text(9)
		      41636d6520496e632e                # "Acme Inc."
		   02                                   # unsigned(2)
		   67                                   # text(7)
		      72722d74726170                    # "rr-trap"
		   03                                   # unsigned(3)
		   69                                   # text(9)
		      41636d6520496e632e                # "Acme Inc."
		   04                                   # unsigned(4)
		   c1                                   # tag(1)
		      00                                # unsigned(0)
		   05                                   # unsigned(5)
		   c1                                   # tag(1)
		      00                                # unsigned(0)
		   06                                   # unsigned(6)
		   c1                                   # tag(1)
		      00                                # unsigned(0)
		   07                                   # unsigned(7)
		   46                                   # bytes(6)
		      ffffffffffff                      # "\xFF\xFF\xFF\xFF\xFF\xFF"
		   0a                                   # unsigned(10)
		   48                                   # bytes(8)
		      0000000000000000                  # "\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"
		   19 0100                              # unsigned(256)
		   51                                   # bytes(17)
		      01deadbeefdeadbeefdeadbeefdeadbeef # "\u0001ޭ\xBE\xEFޭ\xBE\xEFޭ\xBE\xEFޭ\xBE\xEF"
		   19 0102                              # unsigned(258)
		   46                                   # bytes(6)
		      ffffffffffff                      # "\xFF\xFF\xFF\xFF\xFF\xFF"
		   19 0105                              # unsigned(261)
		   18 3c                                # unsigned(60)
		   19 0106                              # unsigned(262)
		   f5                                   # primitive(21)
		   19 0107                              # unsigned(263)
		   01                                   # unsigned(1)
		   19 0108                              # unsigned(264)
		   a2                                   # map(2)
		      01                                # unsigned(1)
		      fb 4028ae147ae147ae               # primitive(4623136420479977390)
		      02                                # unsigned(2)
		      fb 404c63d70a3d70a4               # primitive(4633187891898314916)
	*/
	expected := []byte{
		0xae, 0x01, 0x69, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63,
		0x2e, 0x02, 0x67, 0x72, 0x72, 0x2d, 0x74, 0x72, 0x61, 0x70, 0x03,
		0x69, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x04,
		0xc1, 0x00, 0x05, 0xc1, 0x00, 0x06, 0xc1, 0x00, 0x07, 0x46, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x0a, 0x48, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x19, 0x01, 0x00, 0x51, 0x01, 0xde, 0xad,
		0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde,
		0xad, 0xbe, 0xef, 0x19, 0x01, 0x02, 0x46, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x19, 0x01, 0x05, 0x18, 0x3c, 0x19, 0x01, 0x06, 0xf5,
		0x19, 0x01, 0x07, 0x01, 0x19, 0x01, 0x08, 0xa2, 0x01, 0xfb, 0x40,
		0x28, 0xae, 0x14, 0x7a, 0xe1, 0x47, 0xae, 0x02, 0xfb, 0x40, 0x4c,
		0x63, 0xd7, 0x0a, 0x3d, 0x70, 0xa4,
	}

	cborRoundTripper(t, tv, expected)
}

func TestEat_Full_RoundtripJSON(t *testing.T) {
	tv := fatEat
	expected := `
{
	"eat_nonce": "AAAAAAAAAAA=",
	"oemid": "////////",
	"oemboot": true,
	"dbgstat": 1,
	"location": {
		"lat": 12.34,
		"long": 56.78
	},
	"ueid": "Ad6tvu/erb7v3q2+796tvu8=",
	"uptime": 60,
	"iss": "Acme Inc.",
	"sub": "rr-trap",
	"aud": "Acme Inc.",
	"exp": 0,
	"nbf": 0,
	"iat": 0,
	"cti": "////////"
}`
	// NOTE: cti is not in JSON EAT though
	jsonRoundTripper(t, tv, expected)
}

func TestEat_Submods_RoundtripJSON(t *testing.T) {
	tv := justEatSubmods
	expected := `{
		"submods": {
		  "eat-claims": {},
		  "eat-token": "2D3SQaA="
		}
	  }`

	jsonRoundTripper(t, tv, expected)
}
