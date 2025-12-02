// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDloas_MarshalJSON_OK(t *testing.T) {
	expectedUrl, _ := url.Parse("http://example.com/")

	len2data := `["http://example.com/","foo"]`
	var dloa Dloa

	assert.Nil(t, json.Unmarshal([]byte(len2data), &dloa))
	assert.Equal(t, *expectedUrl, dloa.Registrar)
	assert.Equal(t, "foo", dloa.PlatformLabel)
	assert.Nil(t, dloa.ApplicationLabel)

	encoded, err := json.Marshal(dloa)
	assert.Nil(t, err)
	assert.JSONEq(t, len2data, string(encoded))

	len3data := `["http://example.com/", "foo", "bar"]`
	assert.Nil(t, json.Unmarshal([]byte(len3data), &dloa))
	assert.Equal(t, *expectedUrl, dloa.Registrar)
	assert.Equal(t, "foo", dloa.PlatformLabel)
	assert.Equal(t, "bar", *dloa.ApplicationLabel)

	encoded, err = json.Marshal(dloa)
	assert.Nil(t, err)
	assert.JSONEq(t, len3data, string(encoded))
}

func TestDloas_UnmarshalJSON_NG(t *testing.T) {
	var dloa Dloa
	assert.NotNil(t, json.Unmarshal([]byte(`0`), &dloa))
	assert.NotNil(t, json.Unmarshal([]byte(`["not url"]`), &dloa))
	assert.NotNil(t, json.Unmarshal([]byte(`["http://example.com/"]`), &dloa))
	assert.NotNil(t, json.Unmarshal([]byte(`["http://example.com/",0]`), &dloa))
	assert.NotNil(t, json.Unmarshal([]byte(`["http://example.com/","foo","bar","baz"]`), &dloa))
}

func TestDloas_MarshalCBOR_OK(t *testing.T) {
	expectedUrl, _ := url.Parse("http://example.com/")

	len2data := []byte{
		0x82, // array(2)
		0x73, // text(19) "http://example.com/"
		0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
		0x63,             // text(3)
		0x66, 0x6f, 0x6f, // "foo"
	}
	var dloa Dloa

	assert.Nil(t, dm.Unmarshal(len2data, &dloa))
	assert.Equal(t, *expectedUrl, dloa.Registrar)
	assert.Equal(t, "foo", dloa.PlatformLabel)
	assert.Nil(t, dloa.ApplicationLabel)

	encoded, err := em.Marshal(dloa)
	assert.Nil(t, err)
	assert.Equal(t, len2data, encoded)

	len3data := []byte{
		0x83, // array(3)
		0x73, // text(19) "http://example.com/"
		0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
		0x63,             // text(3)
		0x66, 0x6f, 0x6f, // "foo"
		0x63,             // text(3)
		0x62, 0x61, 0x72, // "bar"
	}
	assert.Nil(t, dm.Unmarshal(len3data, &dloa))
	assert.Equal(t, *expectedUrl, dloa.Registrar)
	assert.Equal(t, "foo", dloa.PlatformLabel)
	assert.Equal(t, "bar", *dloa.ApplicationLabel)

	encoded, err = em.Marshal(dloa)
	assert.Nil(t, err)
	assert.Equal(t, len3data, encoded)
}

func TestDloas_UnmarshalCBOR_NG(t *testing.T) {
	var dloa Dloa
	assert.NotNil(t, dm.Unmarshal([]byte{0x00}, &dloa))
	// echo '["not url"]' | diag2cbor.rb | xxd -p
	assert.NotNil(t, dm.Unmarshal([]byte{0x81, 0x67, 0x6e, 0x6f, 0x74, 0x20, 0x75, 0x72, 0x6c}, &dloa))
	// echo '["http://example.com/"]' | diag2cbor.rb | xxd -p
	assert.NotNil(t, dm.Unmarshal([]byte{
		0x81, 0x73, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	}, &dloa))
	// echo '["http://example.com/",0]' | diag2cbor.rb | xxd -p
	assert.NotNil(t, dm.Unmarshal([]byte{
		0x82, 0x73, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x00,
	}, &dloa))
	// echo '["http://example.com/","foo","bar","baz"]' | diag2cbor.rb | xxd -p
	assert.NotNil(t, dm.Unmarshal([]byte{
		0x84, 0x73, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x66, 0x6f,
		0x6f, 0x63, 0x62, 0x61, 0x72, 0x63, 0x62, 0x61, 0x7a,
	}, &dloa))
}
