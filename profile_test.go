package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	InvalidProfile = 100
	InvalidURL     = "abcd"
	EmptyURL       = ""
	InvalidOID     = "xxxx"
	EmptyOID       = ""
)

// TestProfile_GetSet_Basic_URL tests the basic setting of Profile as URL string
func TestProfile_GetSet_Basic_URL(t *testing.T) {
	inputURL := "https://samplewebsite.com"

	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	expectedURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, inputURL)
	assert.True(t, profile.IsURI())

	inputURL = "https://samplenewwebsite.co.uk"
	err = profile.Set(inputURL)
	assert.Nil(t, err)

	expectedURL, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, inputURL)
	assert.True(t, profile.IsURI())

	// Negative test cases, below
	inputURL = InvalidURL
	err = profile.Set(inputURL)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "profile decode failed no valid URL or OID")

	inputURL = EmptyURL
	err = profile.Set(inputURL)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "profile decode failed no valid URL or OID")
}

// TestProfile_GetSet_Basic_OID tests the saving and retrieval of OID as profile
func TestProfile_GetSet_Basic_OID(t *testing.T) {
	inputOID := "1.2.3.4"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)

	expectedOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, inputOID)
	assert.True(t, profile.IsOID())

	inputOID = "24.43.27.88"
	err = profile.Set(inputOID)
	assert.Nil(t, err)

	expectedOID, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, inputOID)
	assert.True(t, profile.IsOID())

	// Negative test cases below
	// Add a test for OID less than minimum arcs
	inputOID = "56.78"
	err = profile.Set(inputOID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid OID, num arcs: 2 < min OID arcs 3")

	inputOID = InvalidOID
	err = profile.Set(inputOID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to extract OID from string")

	inputOID = EmptyOID
	err = profile.Set(inputOID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no valid OID")
}

// TestProfileURI_MarshalCBOROK tests the CBOR marshaling of profile been set as an URI
func TestProfile_MarshalCBOR_URL(t *testing.T) {
	urlText := "http://example.com"
	profile, err := NewProfile(urlText)
	assert.Nil(t, err)

	// URI is set as a non tagged string within CBOR payload
	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	expected := []byte{
		0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	// Negative Test Cases
	actual, err := profile.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// (corrupted len byte)
	data := []byte{
		0x72, 0x88, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}
	err = profile.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid UTF-8 string")

	// (malformed profile as CBOR )
	data = []byte{
		0x6B, 0x74, 0x65, 0x78, 0x74, 0x20, 0x73, 0x74, 0x72,
		0x69, 0x6E, 0x67,
	}
	err = profile.UnmarshalCBOR(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid profile data")
}

// TestProfile_UnMarshalCBOR_WithURL tests the CBOR UnMarshaling of profile claim set as URL
func TestProfile_UnMarshalCBOR_URL(t *testing.T) {
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data := []byte{
		0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}

	err = profile.UnmarshalCBOR(data)
	assert.Nil(t, err)
	recvURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, inputURL, recvURL)
}

// TestProfile_MarshalCBOR_OID tests the CBOR marshaling of profile set as OID
func TestProfile_MarshalCBOR_OID(t *testing.T) {
	inputOID := "2.5.2.8192"

	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	expected := []byte{
		0x44, 0x55, 0x02, 0xC0, 0x00,
	}
	actual, err := profile.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// Sample test from section 3.1: https://tools.ietf.org/html/draft-ietf-cbor-tags-oid-05#section-3
	inputOID = "2.16.840.1.101.3.4.2.1"
	expected = []byte{
		0x49, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
	}
	err = profile.Set(inputOID)
	assert.Nil(t, err)
	actual, err = profile.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

//  TestProfile_UnMarshalCBOR_OID tests the CBOR unmarshaling of profile as decoded as OID
func TestProfile_UnMarshalCBOR_OID(t *testing.T) {
	inputOID := "2.5.2.8192"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)

	input := []byte{
		0x44, 0x55, 0x02, 0xC0, 0x00,
	}
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	expectedOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, inputOID)

	// Section 3.1 of https://tools.ietf.org/html/draft-ietf-cbor-tags-oid-05#section-3
	inputOID = "2.16.840.1.101.3.4.2.1"

	// CBOR Input
	input = []byte{
		0x49, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
	}
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	expectedOID, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, inputOID)
}

// TestProfile_MarshalJSON_URL tests the JSON Marshaling for a known URL
func TestProfile_MarshalJSON_URL(t *testing.T) {
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	expected := []byte{
		0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}

	actual, err := profile.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

// TestProfile_UnMarshalJSON_URL tests the UnMarshaling of a JSON value as URL string
func TestProfile_UnMarshalJSON_URL(t *testing.T) {
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	data := []byte{
		0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}

	err = profile.UnmarshalJSON(data)
	assert.Nil(t, err)
	recvURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, inputURL, recvURL)

	// Negative Test Cases
	// Corrupted Header
	data = []byte{
		0x10, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x10,
	}
	err = profile.UnmarshalJSON(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid character")

	// Invalid Input data as string
	data = []byte{
		0x22, 0x88, 0x78, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}
	err = profile.UnmarshalJSON(data)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to extract OID")

}

// TestProfile_MarshalJSON_OID validates the JSON Marshaling of OID as profile
func TestProfile_MarshalJSON_OID(t *testing.T) {
	inputOID := "2.5.2.1"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	expected := []byte(`"2.5.2.1"`)
	actual, err := profile.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

// TestProfile_UnMarshalJSON_OID tests the JSON Unmarshaling for OID as profile
func TestProfile_UnMarshalJSON_OID(t *testing.T) {
	inputOID := "1.2.3.4"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	input := []byte(`"2.5.2.1"`)
	err = profile.UnmarshalJSON(input)
	assert.Nil(t, err)
	expectedOID := "2.5.2.1"
	receivedOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, receivedOID)

	// Partial OID test
	inputOID = ".2.3.4"
	_, err = NewProfile(inputOID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to extract OID from string")
}

// TestProfile_RoundTrip_CBOR_Long_OID, tests for a round trip encode/decode of
func TestProfile_RoundTrip_CBOR_Long_OID(t *testing.T) {
	inputOID := "1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1.1.3.6.1.4.1.2706.123.1.2.1"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	expected := []byte{
		0x58, 0xa7, 0x2b, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1,
		0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b,
		0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1,
		0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1,
		0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1,
		0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1,
		0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1,
		0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b,
		0x1, 0x2, 0x1, 0x1, 0x3, 0x6, 0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1, 0x1, 0x3, 0x6,
		0x1, 0x4, 0x1, 0x95, 0x12, 0x7b, 0x1, 0x2, 0x1,
	}
	actual, err := profile.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	// Test UnMarshal CBOR
	input := actual
	expectedOID := inputOID
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	receivedOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, receivedOID)
}
