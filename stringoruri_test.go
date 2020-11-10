package eat

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrURI_Basic(t *testing.T) {
	assert := assert.New(t)

	text := "::Acme Inc."

	var s StringOrURI

	s.FromString(text)
	assert.Equal(text, s.String())
	assert.False(s.IsURI())

	_, err := s.ToURL()
	assert.NotNil(err)

	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(err)

	s.FromURL(u)
	assert.Equal(urlText, s.String())
	assert.True(s.IsURI())

	s.FromString(urlText)
	assert.False(s.IsURI())

	newU, err := s.ToURL()
	assert.Nil(err)
	assert.Equal(u, newU)
}

func TestStringOrURI_MarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(err)

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
	assert.Nil(err)
	assert.Equal(expected, actual)

	s.FromString("::Acme Inc.")

	// 6b                        # text(11)
	//    3a3a41636d6520496e632e # "::Acme Inc."
	expected = []byte{
		0x6b, 0x3a, 0x3a, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e,
	}

	actual, err = s.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestStringOrURI_UnmarshalCBOR(t *testing.T) {
	assert := assert.New(t)

	urlText := "http://example.com"
	u, err := url.Parse(urlText)
	assert.Nil(err)

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
	assert.Nil(err)
	assert.Equal(expected, actual)

	// 6b                        # text(11)
	//    3a3a41636d6520496e632e # "::Acme Inc."
	data = []byte{
		0x6b, 0x3a, 0x3a, 0x41, 0x63, 0x6d, 0x65, 0x20, 0x49, 0x6e, 0x63, 0x2e,
	}
	expected.FromString("::Acme Inc.")

	err = actual.UnmarshalCBOR(data)
	assert.Nil(err)
	assert.Equal(expected, actual)

	// Bad tag value (corrupted initial byte)
	// d7 20                                   # tag [corrupted]
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data = []byte{
		0xd7, 0x20, 0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "must be URI")

	// not a text string or URI tag
	// 19 0539 # unsigned(1337)
	data = []byte{
		0x19, 0x05, 0x39,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "must be a text string or a URI tag")

	// Non-string inside URI tag
	// d8 20      # tag(32)
	//    19 0539 # unsigned(1337)
	data = []byte{
		0xd8, 0x20, 0x19, 0x05, 0x39,
	}

	err = actual.UnmarshalCBOR(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "URI tag value must be a string")
}
