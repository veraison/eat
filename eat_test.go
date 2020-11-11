// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ueID = UEID{
		0x01, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}
	oemID       = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	nonce       = Nonce{[]byte{0xde, 0xad, 0xbe, 0xef}}
	AcmeInc     = "Acme Inc."
	origination = StringOrURI{text: &AcmeInc}
	secLevel    = SecurityLevel(SecLevelHardware)
	secBoot     = true
	debug       = Debug(DebugDisabled)
	location    = Location{Latitude: 12.34, Longitude: 56.78}
	uptime      = uint(60)

	fatEat = Eat{
		Nonce:         &Nonces{nonce},
		UEID:          &ueID,
		Origination:   &origination,
		OemID:         &oemID,
		SecurityLevel: &secLevel,
		SecureBoot:    &secBoot,
		Debug:         &debug,
		Location:      &location,
		Uptime:        &uptime,
	}
)

func cborRoundTripper(t *testing.T, tv Eat, expected []byte) {
	data, err := tv.ToCBOR()

	assert.Nil(t, err)
	assert.Equal(t, expected, data)

	actual := Eat{}
	err = actual.FromCBOR(data)

	assert.Nil(t, err)
	assert.Equal(t, tv, actual)
}

func jsonRoundTripper(t *testing.T, tv Eat, expected string) {
	data, err := tv.ToJSON()

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
	   a9                                       # map(9)
	      0a                                    # unsigned(10)
	      44                                    # bytes(4)
	         deadbeef                           # "\xDE\xAD\xBE\xEF"
	      0b                                    # unsigned(11)
	      51		                    # bytes(17)
	         01deadbeefdeadbeefdeadbeefdeadbeef # "\x01\xDE\xAD\[...]xBE\xEF"
	      0c                                    # unsigned(12)
	      69                                    # text(9)
	         41636d6520496e632e                 # "Acme Inc."
	      0d                                    # unsigned(13)
	      46                                    # bytes(6)
	         ffffffffffff                       # "\xFF\xFF\xFF\xFF\xFF\xFF"
	      0e                                    # unsigned(14)
	      03                                    # unsigned(3)
	      0f                                    # unsigned(15)
	      f5                                    # primitive(21)
	      10                                    # unsigned(16)
	      01                                    # unsigned(1)
	      11                                    # unsigned(17)
	      a2                                    # map(2)
	         01                                 # unsigned(1)
	         fb 4028ae147ae147ae                # primitive(4623136420479977390)
	         02                                 # unsigned(2)
	         fb 404c63d70a3d70a4                # primitive(4633187891898314916)
	      13                                    # unsigned(19)
	      18 3c                                 # unsigned(60)
	*/
	expected := []byte{
		0xa9, 0x0a, 0x44, 0xde, 0xad, 0xbe, 0xef, 0x0b, 0x51, 0x01,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0xde, 0xad,
		0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, 0x0c, 0x69, 0x41, 0x63,
		0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x0d, 0x46, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x0e, 0x03, 0x0f, 0xf5, 0x10,
		0x01, 0x11, 0xa2, 0x01, 0xfb, 0x40, 0x28, 0xae, 0x14, 0x7a,
		0xe1, 0x47, 0xae, 0x02, 0xfb, 0x40, 0x4c, 0x63, 0xd7, 0x0a,
		0x3d, 0x70, 0xa4, 0x13, 0x18, 0x3c,
	}

	cborRoundTripper(t, tv, expected)
}

func TestEat_Full_RoundtripJSON(t *testing.T) {
	tv := fatEat
	expected := `
{
	"nonce": "3q2+7w==",
	"origination": "Acme Inc.",
	"oemid": "////////",
	"security-level": 3,
	"secure-boot": true,
	"debug-disable": 1,
	"location": {
		"lat": 12.34,
		"long": 56.78
	},
	"ueid": "Ad6tvu/erb7v3q2+796tvu8=",
	"uptime": 60
}`
	jsonRoundTripper(t, tv, expected)
}
