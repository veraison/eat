// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUEID_Verify(t *testing.T) {
	u0 := UEID{}
	assert.EqualError(t, u0.Validate(), "empty UEID")

	u1 := UEID{
		0x01, // RAND
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, // 16 bytes
	}
	assert.Nil(t, u1.Validate())

	u2 := UEID{
		0x01, // RAND
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, // 15 bytes
	}
	assert.EqualError(t, u2.Validate(), "RAND length must be exactly 16, 24, or 32 bytes; found 15 bytes")

	u3 := UEID{
		0x02,                               // EUI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 6 bytes
	}
	assert.Nil(t, u3.Validate())

	u4 := UEID{
		0x02, // EUI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, // 15 bytes
	}
	assert.EqualError(t, u4.Validate(), "EUI length must be exactly 6 (EUI-48) or 8 (EUI-60 or EUI-64) bytes; found 15 bytes")

	u5 := UEID{
		0x03, // IMEI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 14 bytes
	}
	assert.Nil(t, u5.Validate())

	u6 := UEID{
		0x03, // IMEI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, // 15 bytes
	}
	assert.EqualError(t, u6.Validate(), "IMEI length must be exactly 14 bytes; found 15 bytes")

	u7 := UEID{
		0xFF, // Invalid
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
	}
	assert.EqualError(t, u7.Validate(), "invalid UEID type 255")

}

func TestUEID_JSONMarshal_OK(t *testing.T) {
	e1 := UEID{
		0x01, // RAND
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, // 16 bytes
	}
	s1 := []byte(`"Ad6tvu_erb7v3q2-796tvu8"`)
	t1, err := json.Marshal(e1)
	assert.Nil(t, err)
	assert.Equal(t, s1, t1)

	var u1 UEID
	err = json.Unmarshal(s1, &u1)
	assert.Nil(t, err)
	assert.Equal(t, e1, u1)

	e2 := UEID{
		0x02,                               // EUI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 6 bytes
	}
	s2 := []byte(`"At6tvu_erQ"`)

	t2, err := json.Marshal(e2)
	assert.Nil(t, err)
	assert.Equal(t, s2, t2)

	var u2 UEID
	err = json.Unmarshal(s2, &u2)
	assert.Nil(t, err)
	assert.Equal(t, e2, u2)

	e3 := UEID{
		0x03, // IMEI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 14 bytes
	}
	s3 := []byte(`"A96tvu_erb7v3q2-796t"`)

	t3, err := json.Marshal(e3)
	assert.Nil(t, err)
	assert.Equal(t, s3, t3)

	var u3 UEID
	err = json.Unmarshal(s3, &u3)
	assert.Nil(t, err)
	assert.Equal(t, e3, u3)
}

func TestUEID_JSONUnmarshal_NG(t *testing.T) {
	// not a string
	s1 := []byte(`0`)
	var u1 UEID
	err := json.Unmarshal(s1, &u1)
	assert.NotNil(t, err)

	// not base64
	s2 := []byte(`"&"`)
	var u2 UEID
	err = json.Unmarshal(s2, &u2)
	assert.NotNil(t, err)

	// base64 but not base64url
	s3 := []byte(`"A96tvu/erb7v3q2+796t"`)
	var u3 UEID
	err = json.Unmarshal(s3, &u3)
	assert.NotNil(t, err)
}
