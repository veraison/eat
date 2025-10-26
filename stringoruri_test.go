// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrURI_Basic(t *testing.T) {
	text := "% Acme Inc %"

	var s StringOrURI

	err := s.FromString(text)
	assert.Nil(t, err)
	assert.Equal(t, text, s.String())
	assert.False(t, s.IsURI())

	_, err = s.ToURL()
	assert.NotNil(t, err)

	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(t, err)

	s.FromURL(u)
	assert.Equal(t, urlText, s.String())
	assert.True(t, s.IsURI())

	err = s.FromString(urlText)
	assert.Nil(t, err)
	assert.True(t, s.IsURI())

	newU, err := s.ToURL()
	assert.Nil(t, err)
	assert.Equal(t, u, newU)
}

func TestStringOrURI_MarshalCBOR(t *testing.T) {
	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(t, err)

	var s StringOrURI
	s.FromURL(u)

	// d8 20                                   # tag(32)
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	expected := []byte{
		0xd8, 0x20, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	actual, err := s.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = s.FromString("% Acme Inc. %")
	assert.Nil(t, err)

	// 6d                            # text(13)
	// 252041636d6520496e632e2025 # "% Acme Inc. %"
	expected = []byte{
		0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x20, 0x25,
	}

	actual, err = s.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestStringOrURI_UnmarshalCBOR(t *testing.T) {
	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(t, err)

	var expected StringOrURI
	expected.FromURL(u)

	// d8 20                                   # tag(32)
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data := []byte{
		0xd8, 0x20, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	var actual StringOrURI
	err = actual.UnmarshalCBOR(data)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// 6b                        # text(11)
	//    3a3a41636d6520496e632e # "::Acme Inc."
	data = []byte{
		0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x20, 0x25,
	}
	err = expected.FromString("% Acme Inc. %")
	assert.Nil(t, err)

	err = actual.UnmarshalCBOR(data)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// Bad tag value (corrupted initial byte)
	// d7 20                                   # tag [corrupted]
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data = []byte{
		0xd7, 0x20, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "bytes of extraneous data starting at index")

	// Bad tag value (corrupted initial byte)
	// d8 21                                   # tag [corrupted]
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data = []byte{
		0xd8, 0x21, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must be URI")

	// not a text string or URI tag
	// 19 0539 # unsigned(1337)
	data = []byte{
		0x19, 0x05, 0x39,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "must be a text string or a URI tag")

	// Non-string inside URI tag
	// d8 20      # tag(32)
	//    19 0539 # unsigned(1337)
	data = []byte{
		0xd8, 0x20, 0x19, 0x05, 0x39,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "URI tag value must be a string")
}

func TestStringOrUri_UnmarshalJSON(t *testing.T) {
	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(t, err)

	var expected StringOrURI
	expected.FromURL(u)

	data := []byte(`"` + urlText + `"`)

	var actual StringOrURI
	err = actual.UnmarshalJSON(data)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	err = expected.FromString("% Acme Inc. %")
	assert.Nil(t, err)

	data = []byte(`"% Acme Inc. %"`)

	err = actual.UnmarshalJSON(data)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestStringOrUri_UnmarshalJSON_NG(t *testing.T) {
	var s StringOrURI

	// not a string
	data := []byte(`134`)

	err := s.UnmarshalJSON(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "cannot unmarshal number into Go value of type string")

	// invalid URI
	data = []byte(`"http://[::1"`)

	err = s.UnmarshalJSON(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "parse")
}
