// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOEMID_MarshalJSON_OK(t *testing.T) {
	e1 := OEMID{Value: 76543}
	t1, err := json.Marshal(e1)
	assert.Nil(t, err)
	assert.Equal(t, []byte("76543"), t1)

	var o1 OEMID
	err = json.Unmarshal([]byte("76543"), &o1)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, o1.Value)

	e2 := OEMID{Value: B64Url{0x89, 0x48, 0x23}}
	t2, err := json.Marshal(e2)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"iUgj"`), t2)

	var o2 OEMID
	err = json.Unmarshal([]byte(`"iUgj"`), &o2)
	assert.Nil(t, err)
	assert.Equal(t, e2.Value, o2.Value)

	e3 := OEMID{Value: B64Url{
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
	}}
	t3, err := json.Marshal(e3)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"m--Hh-uhPiyPbny0sfRhmg"`), t3)

	var o3 OEMID
	err = json.Unmarshal([]byte(`"m--Hh-uhPiyPbny0sfRhmg"`), &o3)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, o1.Value)
}

func TestOEMID_MarshalJSON_NG(t *testing.T) {
	e1 := OEMID{Value: 1.5}
	_, err := json.Marshal(e1)
	assert.NotNil(t, err)

	var o1 OEMID
	err = json.Unmarshal([]byte("1.5"), &o1)
	assert.NotNil(t, err)

	e2 := OEMID{Value: B64Url{0x89, 0x48}}
	_, err = json.Marshal(e2)
	assert.NotNil(t, err)

	var o2 OEMID
	err = json.Unmarshal([]byte(`"iUg"`), &o2)
	assert.NotNil(t, err)

	e3 := OEMID{Value: B64Url{
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
		0x00, // extra bytes
	}}
	_, err = json.Marshal(e3)
	assert.NotNil(t, err)

	var o3 OEMID
	err = json.Unmarshal([]byte(`"m--Hh-uhPiyPbny0sfRhmgA"`), &o3)
	assert.NotNil(t, err)

	// see section 2 of RFC 9711
	// base64url encoding: base64 encoding using the URL- and filename-safe
	//    character set defined in Section 5 of [RFC4648], with all trailing
	//    '=' characters omitted and without the inclusion of any line
	//    breaks, whitespace, or other additional characters [RFC7515].
	var o4 OEMID
	err = json.Unmarshal([]byte(`"m--Hh-uhPiyPbny0sfRhmg=="`), &o4)
	assert.NotNil(t, err)
}

func TestOEMID_MarshalCBOR_OK(t *testing.T) {
	oemidPen := []byte{0x1A, 0x00, 0x01, 0x2A, 0xFF}
	e1 := OEMID{Value: 76543}
	assert.True(t, e1.Valid())

	t1, err := em.Marshal(e1)
	assert.Nil(t, err)
	assert.Equal(t, oemidPen, t1)

	var o1 OEMID
	err = dm.Unmarshal(t1, &o1)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, o1.Value)

	oemidIEEE := []byte{0x43, 0x89, 0x48, 0x23}
	e2 := OEMID{Value: B64Url{0x89, 0x48, 0x23}}
	assert.True(t, e2.Valid())

	t2, err := em.Marshal(e2)
	assert.Nil(t, err)
	assert.Equal(t, oemidIEEE, t2)

	var o2 OEMID
	err = dm.Unmarshal(t2, &o2)
	assert.Nil(t, err)
	assert.Equal(t, e2.Value, o2.Value)

	oemidRandom := []byte{
		0x50,
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
	}
	e3 := OEMID{Value: B64Url{
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
	}}
	t3, err := em.Marshal(e3)
	assert.True(t, e3.Valid())

	assert.Nil(t, err)
	assert.Equal(t, oemidRandom, t3)

	var o3 OEMID
	err = dm.Unmarshal(t3, &o3)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, o1.Value)
}

func TestOEMID_MarshalCBOR_NG(t *testing.T) {
	e1 := OEMID{Value: 1.5}
	assert.False(t, e1.Valid())
	_, err := em.Marshal(e1)
	assert.NotNil(t, err)

	var o1 OEMID
	err = dm.Unmarshal([]byte{0xF9, 0x3E, 0x00}, &o1)
	assert.NotNil(t, err)

	e2 := OEMID{Value: B64Url{0x89, 0x48}}
	assert.False(t, e2.Valid())
	_, err = json.Marshal(e2)
	assert.NotNil(t, err)

	var o2 OEMID
	err = json.Unmarshal([]byte{0x42, 0x89, 0x48}, &o2)
	assert.NotNil(t, err)

	e3 := OEMID{Value: B64Url{
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
		0x00, // extra bytes
	}}
	assert.False(t, e3.Valid())
	_, err = json.Marshal(e3)
	assert.NotNil(t, err)

	var o3 OEMID
	err = json.Unmarshal([]byte{
		0x51,
		0x9b, 0xef, 0x87, 0x87, 0xeb, 0xa1, 0x3e, 0x2c,
		0x8f, 0x6e, 0x7c, 0xb4, 0xb1, 0xf4, 0x61, 0x9a,
		0x00,
	}, &o3)
	assert.NotNil(t, err)
}
