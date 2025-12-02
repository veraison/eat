// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrBinary_MarshalJSON_OK(t *testing.T) {
	e1 := StringOrBinary{Value: "not base64url encoded string"}
	t1, err := json.Marshal(e1)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"not base64url encoded string"`), t1)

	var s1 StringOrBinary
	err = json.Unmarshal([]byte(`"not base64url encoded string"`), &s1)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, s1.Value)

	// a base64 encoded text is treated as a tstr
	e2 := StringOrBinary{Value: "AA/="}
	t2, err := json.Marshal(e2)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"AA/="`), t2)

	var s2 StringOrBinary
	err = json.Unmarshal([]byte(`"AA/="`), &s2)
	assert.Nil(t, err)
	assert.Equal(t, e2.Value, s2.Value)

	// an empty string must be treated as tstr
	e3 := StringOrBinary{Value: ""}
	t3, err := json.Marshal(e3)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`""`), t3)

	var s3 StringOrBinary
	err = json.Unmarshal([]byte(`""`), &s3)
	assert.Nil(t, err)
	assert.Equal(t, e3.Value, s3.Value)

	// only base64url encoded text is treated as B64Url
	e4 := StringOrBinary{Value: B64Url{0x00, 0x00}}
	t4, err := json.Marshal(e4)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"AAA"`), t4)

	var s4 StringOrBinary
	err = json.Unmarshal([]byte(`"AAA"`), &s4)
	assert.Nil(t, err)
	assert.Equal(t, e4.Value, s4.Value)
}

func TestStringOrBinary_MarshalJSON_NG(t *testing.T) {
	e1 := StringOrBinary{Value: 1.5}
	_, err := json.Marshal(e1)
	assert.NotNil(t, err)

	var s1 StringOrBinary
	err = json.Unmarshal([]byte("1.5"), &s1)
	assert.NotNil(t, err)
}

func TestStringOrBinary_MarshalCBOR_OK(t *testing.T) {
	encodedTextString := []byte{0x64, 0x74, 0x65, 0x73, 0x74}
	e1 := StringOrBinary{Value: "test"}
	t1, err := em.Marshal(e1)
	assert.Nil(t, err)
	assert.Equal(t, encodedTextString, t1)

	var s1 StringOrBinary
	err = dm.Unmarshal(encodedTextString, &s1)
	assert.Nil(t, err)
	assert.Equal(t, e1.Value, s1.Value)

	encodedByteString := []byte{0x42, 0x00, 0x00}
	e2 := StringOrBinary{Value: B64Url{0x00, 0x00}}
	t2, err := em.Marshal(e2)
	assert.Nil(t, err)
	assert.Equal(t, encodedByteString, t2)

	var s2 StringOrBinary
	err = dm.Unmarshal(encodedByteString, &s2)
	assert.Nil(t, err)
	assert.Equal(t, e2.Value, s2.Value)
}

func TestStringOrBinary_MarshalCBOR_NG(t *testing.T) {
	e1 := StringOrBinary{Value: 1.5}
	_, err := em.Marshal(e1)
	assert.NotNil(t, err)

	var s1 StringOrBinary
	err = dm.Unmarshal([]byte{0xF9, 0x3E, 0x00}, &s1)
	assert.NotNil(t, err)
}
