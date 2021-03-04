package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	InvalidProfile = 100
	InvalidURL     = "abcd"
	EmptyURL       = "  "
	InvalidOID     = "xxxx"
	EmptyOID       = "  "
)

// TestProfileURL_Basic tests the basic setting of Profile as URL string
func TestProfileURL_Basic(t *testing.T) {
	assert := assert.New(t)
	inputURL := "https://samplewebsite.com"

	profile, err := NewProfile(inputURL)
	assert.Nil(err)

	expectedURL, err := profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedURL, inputURL)
	assert.True(profile.IsURI())

	inputURL = "https://samplenewwebsite.co.uk"
	err = profile.SetProfile(inputURL)
	assert.Nil(err)

	expectedURL, err = profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedURL, inputURL)
	assert.True(profile.IsURI())
}

// TestProfileURL_NOK tests the negative cases of invalid URL
func TestProfileURL_NOK(t *testing.T) {
	assert := assert.New(t)
	inputURL := "https://samplewebsite.com"

	profile, err := NewProfile(inputURL)
	assert.Nil(err)

	inputURL = InvalidURL
	err = profile.SetProfile(inputURL)
	assert.NotNil(err)
	assert.Contains(err.Error(), "no valid URI or OID")

	inputURL = EmptyURL
	err = profile.SetProfile(inputURL)
	assert.NotNil(err)
	assert.Contains(err.Error(), "no valid URI or OID")
}

// TestProfileOID_Basic tests the saving and retrieval of OID as profile
func TestProfileOID_Basic(t *testing.T) {
	assert := assert.New(t)
	inputOID := "1.2"

	profile, err := NewProfile(inputOID)
	assert.Nil(err)

	expectedOID, err := profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedOID, inputOID)
	assert.True(profile.IsOID())

	inputOID = "24.43.27.88"
	err = profile.SetProfile(inputOID)
	assert.Nil(err)

	expectedOID, err = profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedOID, inputOID)
	assert.True(profile.IsOID())
}

// TestProfileOID_NOK tests for invalid OID cases
func TestProfileOID_NOK(t *testing.T) {
	assert := assert.New(t)
	inputOID := "1.2"

	profile, err := NewProfile(inputOID)
	assert.Nil(err)

	inputOID = InvalidOID
	err = profile.SetProfile(InvalidOID)
	assert.NotNil(err)
	assert.Contains(err.Error(), "no valid URI or OID")

	inputOID = EmptyOID
	err = profile.SetProfile(inputOID)
	assert.NotNil(err)
	assert.Contains(err.Error(), "no valid URI or OID")
}

// TestProfileURI_MarshalCBOROK tests the CBOR marshalling of profile been set as an URI
func TestProfileURI_MarshalCBOROK(t *testing.T) {
	assert := assert.New(t)
	urlText := "http://example.com"
	profile, err := NewProfile(urlText)
	assert.Nil(err)

	// URI is set as a non tagged string within CBOR payload
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	expected := []byte{
		0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	actual, err := profile.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

// TestProfileURI_UnMarshalCBOR tests the CBOR UnMarshalling of profile claim set as URI
func TestProfileURI_UnMarshalCBOR(t *testing.T) {
	assert := assert.New(t)
	inputURL := "http://example.com"

	profile, err := NewProfile(inputURL)
	assert.Nil(err)

	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data := []byte{
		0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	err = profile.UnmarshalCBOR(data)
	assert.Nil(err)
	recvURL, err := profile.GetProfile()
	assert.Nil(err)
	assert.Equal(inputURL, recvURL)
}

// TestProfile_UnMarshalCBORNOK tests the negative cases of failed CBOR decodiing
// due to malformed CBOR header or a corrupted CBOR PDU
func TestProfile_UnMarshalCBORNOK(t *testing.T) {
	assert := assert.New(t)
	profile := Profile{}
	// (corrupted initial byte to 0x20)
	data := []byte{
		0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}
	err := profile.UnmarshalCBOR(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "malformed CBOR")

	// (malformed profile as CBOR )
	data = []byte{
		0x72, 0x77, 0x81, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}
	err = profile.UnmarshalCBOR(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "decoding of CBOR data failed")
}

// TestProfileOID_MarshalCBOR tests the CBOR marshaling of profile set as OID
func TestProfileOID_MarshalCBOR(t *testing.T) {
	assert := assert.New(t)
	inputOID := "2.5.2.8192"

	profile, err := NewProfile(inputOID)
	expected := []byte{
		0x46, 0x06, 0x04, 0x55, 0x02, 0xC0, 0x00,
	}
	actual, err := profile.MarshalCBOR()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

// TestProfileOID_UnMarshalCBOR tests the CBOR unmarshaling of profile as decoded as OID
func TestProfileOID_UnMarshalCBOR(t *testing.T) {
	assert := assert.New(t)
	inputOID := "2.5.2.8192"

	profile, err := NewProfile(inputOID)
	assert.Nil(err)

	input := []byte{
		0x46, 0x06, 0x04, 0x55, 0x02, 0xC0, 0x00,
	}
	err = profile.UnmarshalCBOR(input)
	assert.Nil(err)
	expectedOID, err := profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedOID, inputOID)
}

// TestProfileURI_MarshalJSON tests the JSON Marshalling for a known URL
func TestProfileURI_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(err)

	expected := []byte{
		0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}

	actual, err := profile.MarshalJSON()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

// TestProfileURI_UnMarshalJSON tests the UnMarshalling of a JSON value as URL string
func TestProfileURI_UnMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(err)

	data := []byte{
		0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}

	err = profile.UnmarshalJSON(data)
	assert.Nil(err)
	recvURL, err := profile.GetProfile()
	assert.Equal(inputURL, recvURL)
}

// TestProfile_UnMarshalJSONNOK tests hthe negative cases of JSON Unmarshalling
// due to corrupted header or malformed JSON PDU
func TestProfile_UnMarshalJSONNOK(t *testing.T) {
	assert := assert.New(t)
	profile := Profile{}
	// Corrupted Header
	data := []byte{
		0x10, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x10,
	}
	err := profile.UnmarshalJSON(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "invalid character")

	// Invalid Input data as string
	data = []byte{
		0x22, 0x88, 0x78, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}
	err = profile.UnmarshalJSON(data)
	assert.NotNil(err)
	assert.Contains(err.Error(), "json decode of profile failed")
}

// TestProfileOID_MarshalJSON validates the JSON Marshalling of OID as profile
func TestProfileOID_MarshalJSON(t *testing.T) {
	assert := assert.New(t)
	inputOID := "2.5.2.1"
	profile, err := NewProfile(inputOID)
	assert.Nil(err)
	expected := []byte(`"2.5.2.1"`)
	actual, err := profile.MarshalJSON()
	assert.Nil(err)
	assert.Equal(expected, actual)
}

// TestProfileOID_UnMarshalJSON tests the JSON Unmarshalling for OID as profile
func TestProfileOID_UnMarshalJSON(t *testing.T) {
	assert := assert.New(t)
	inputOID := "1.2.3.4"
	profile, err := NewProfile(inputOID)
	assert.Nil(err)
	input := []byte(`"2.5.2.1"`)
	profile.UnmarshalJSON(input)
	expectedOID := "2.5.2.1"
	receivedOID, err := profile.GetProfile()
	assert.Nil(err)
	assert.Equal(expectedOID, receivedOID)
}
