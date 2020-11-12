// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAudience_MarshalCBOR_Multiple(t *testing.T) {
	assert := assert.New(t)

	s := "% Acme Inc. %"
	u, err := url.Parse("http://example.com")
	assert.Nil(err)

	value := Audience{
		StringOrURI{text: &s},
		StringOrURI{uri: u},
	}

	//82                                            # array(2)
	//   6d                                         # text(13)
	//      252041636d6520496e632e2025              # "% Acme Inc. %"
	//   d8 20                                      # tag(32)
	//      72                                      # text(18)
	//         687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	expected := []byte{
		0x82, 0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49,
		0x6e, 0x63, 0x2e, 0x20, 0x25, 0xd8, 0x20, 0x72, 0x68, 0x74,
		0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	actual, err := value.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestAudience_MarshalCBOR_Single(t *testing.T) {
	assert := assert.New(t)

	s := "% Acme Inc. %"

	value := Audience{StringOrURI{text: &s}}

	// 6d                            # text(13)
	//    252041636d6520496e632e2025 # "% Acme Inc. %"
	expected := []byte{
		0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x20, 0x25,
	}

	actual, err := value.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestAudience_UnmarshalCBOR_Multiple(t *testing.T) {
	assert := assert.New(t)

	s := "% Acme Inc. %"
	u, err := url.Parse("http://example.com")
	assert.Nil(err)

	expected := Audience{
		StringOrURI{text: &s},
		StringOrURI{uri: u},
	}

	//82                                            # array(2)
	//   6d                                         # text(13)
	//      252041636d6520496e632e2025              # "% Acme Inc. %"
	//   d8 20                                      # tag(32)
	//      72                                      # text(18)
	//         687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data := []byte{
		0x82, 0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49,
		0x6e, 0x63, 0x2e, 0x20, 0x25, 0xd8, 0x20, 0x72, 0x68, 0x74,
		0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70,
		0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	actual := Audience{}
	err = actual.UnmarshalCBOR(data)

	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestAudience_UnmarshalCBOR_Single(t *testing.T) {
	assert := assert.New(t)

	s := "% Acme Inc. %"
	u, err := url.Parse("http://example.com")
	assert.Nil(err)

	expected := Audience{StringOrURI{text: &s}}

	actual := Audience{}

	// 6d                            # text(13)
	//    252041636d6520496e632e2025 # "% Acme Inc. %"
	data := []byte{
		0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e, 0x20, 0x25,
	}

	err = actual.UnmarshalCBOR(data)
	assert.Nil(err)
	assert.Equal(expected, actual)

	//81                                            # array(1)
	//   6d                                         # text(13)
	//      252041636d6520496e632e2025              # "% Acme Inc. %"
	data2 := []byte{
		0x81, 0x6d, 0x25, 0x20, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49,
		0x6e, 0x63, 0x2e, 0x20, 0x25,
	}

	err = actual.UnmarshalCBOR(data2)
	assert.Nil(err)
	assert.Equal(expected, actual)

	// d8 20                                   # tag(32)
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data3 := []byte{
		0xd8, 0x20, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	expected3 := Audience{StringOrURI{uri: u}}

	err = actual.UnmarshalCBOR(data3)
	assert.Nil(err)
	assert.Equal(expected3, actual)

}
