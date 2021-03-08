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
	NegativeOID    = "-1.2.3.-4"
)

// TestProfile_GetSet_Basic_URL_OK tests the valid setting of Profile as URL string
func TestProfile_GetSet_Basic_URL_OK(t *testing.T) {
	inputURL := "https://samplewebsite.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	expectedURL := inputURL
	actualURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, actualURL)
	assert.True(t, profile.IsURI())

	inputURL = "https://samplenewwebsite.co.uk"
	err = profile.Set(inputURL)
	assert.Nil(t, err)
	expectedURL = inputURL
	actualURL, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, actualURL)
	assert.True(t, profile.IsURI())
}

// TestProfile_GetSet_Basic_URL_NOK tests the invalid setting of Profile as URL string
func TestProfile_GetSet_Basic_URL_NOK(t *testing.T) {
	profile := &Profile{}

	// Negative test cases, below
	inputURL := InvalidURL
	expectedError := `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `failed to extract OID from string: `
	expectedError += `strconv.Atoi: parsing "abcd": invalid syntax`
	err := profile.Set(inputURL)
	assert.EqualError(t, err, expectedError)

	inputURL = EmptyURL
	expectedError = `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `no valid OID`
	err = profile.Set(inputURL)
	assert.EqualError(t, err, expectedError)
}

// TestProfile_GetSet_Basic_OID_OK tests the valid case of saving and retrieval of OID as profile
func TestProfile_GetSet_Basic_OID_OK(t *testing.T) {
	inputOID := "1.2.3.4"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)

	expectedOID := inputOID
	actualOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)
	assert.True(t, profile.IsOID())

	inputOID = "24.43.27.88"
	err = profile.Set(inputOID)
	assert.Nil(t, err)
	expectedOID = inputOID
	actualOID, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)
	assert.True(t, profile.IsOID())
}

// TestProfile_GetSet_Basic_OID_NOK tests the invalid case of saving of OID as profile
func TestProfile_GetSet_Basic_OID_NOK(t *testing.T) {
	profile := &Profile{}
	// Add a test for OID less than minimum arcs
	inputOID := "56.78"
	expectedError := `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `invalid OID: `
	expectedError += `got 2 arcs, expecting at least 3`
	err := profile.Set(inputOID)
	assert.EqualError(t, err, expectedError)

	inputOID = InvalidOID
	expectedError = `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `failed to extract OID from string: `
	expectedError += `strconv.Atoi: parsing "xxxx": invalid syntax`
	err = profile.Set(inputOID)
	assert.EqualError(t, err, expectedError)

	inputOID = EmptyOID
	expectedError = `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `no valid OID`
	err = profile.Set(inputOID)
	assert.EqualError(t, err, expectedError)

	inputOID = NegativeOID
	expectedError = `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `invalid OID: negative subidentifier -1 not allowed`
	err = profile.Set(inputOID)
	assert.EqualError(t, err, expectedError)
}

// TestProfile_MarshalCBOR_URL_OK tests the valid case of CBOR marshaling of profile been set as an URI
func TestProfile_MarshalCBOR_URL_OK(t *testing.T) {
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

	actual, err := profile.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

// TestProfile_MarshalCBOR_URL_NOK tests the invalid case of CBOR marshaling of profile been set as an URI
func TestProfile_MarshalCBOR_URL_NOK(t *testing.T) {
	profile := &Profile{}
	// (corrupted len byte)
	data := []byte{
		0x72, 0x88, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}
	expectedError := `CBOR decoding of profile failed: `
	expectedError += `cbor: `
	expectedError += `invalid UTF-8 string`
	err := profile.UnmarshalCBOR(data)
	assert.EqualError(t, err, expectedError)

	// (malformed profile as CBOR )
	data = []byte{
		0x6B, 0x74, 0x65, 0x78, 0x74, 0x20, 0x73, 0x74, 0x72,
		0x69, 0x6E, 0x67,
	}
	expectedError = `invalid profile data: `
	expectedError += `profile URL not in absolute form`
	err = profile.UnmarshalCBOR(data)
	assert.EqualError(t, err, expectedError)
}

// TestProfile_UnmarshalCBOR_URL_OK tests the valid case of CBOR UnMarshaling of profile claim set as URL
func TestProfile_UnmarshalCBOR_URL_OK(t *testing.T) {
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	//    72                                   # text(18)
	//          687474703a2f2f6578616d706c652e636f6d # "http://example.com"
	data := []byte{
		0x72, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65,
		0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d,
	}
	expectedURL := inputURL
	err = profile.UnmarshalCBOR(data)
	assert.Nil(t, err)
	actualURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, actualURL)
}

// TestProfile_MarshalCBOR_OID_OK tests the valid CBOR marshaling of profile set as OID
func TestProfile_MarshalCBOR_OID_OK(t *testing.T) {
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

//  TestProfile_UnmarshalCBOR_OID_OK tests the valid CBOR unmarshaling of profile decoded as OID
func TestProfile_UnmarshalCBOR_OID_OK(t *testing.T) {
	inputOID := "2.5.2.8192"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)

	input := []byte{
		0x44, 0x55, 0x02, 0xC0, 0x00,
	}
	expectedOID := inputOID
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	actualOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)

	// Section 3.1 of https://tools.ietf.org/html/draft-ietf-cbor-tags-oid-05#section-3
	inputOID = "2.16.840.1.101.3.4.2.1"

	// CBOR Input
	input = []byte{
		0x49, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
	}
	expectedOID = inputOID
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	actualOID, err = profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)
}

// TestProfile_MarshalJSON_URL_OK tests the JSON Marshaling for a known URL
func TestProfile_MarshalJSON_URL_OK(t *testing.T) {
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

// TestProfile_UnmarshalJSON_URL_OK tests the Unmarshaling of a JSON value as URL string
func TestProfile_UnmarshalJSON_URL_OK(t *testing.T) {
	inputURL := "http://example.com"
	profile, err := NewProfile(inputURL)
	assert.Nil(t, err)

	data := []byte{
		0x22, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}
	expectedURL := inputURL
	err = profile.UnmarshalJSON(data)
	assert.Nil(t, err)
	actualURL, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedURL, actualURL)
}

// TestProfile_UnmarshalJSON_URL_NOK tests the invalid case of Unmarshaling of a JSON value as URL string
func TestProfile_UnmarshalJSON_URL_NOK(t *testing.T) {
	profile := &Profile{}
	// Corrupted Header
	data := []byte{
		0x10, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x10,
	}

	expectedError := `invalid character '\x10' looking for beginning of value`
	err := profile.UnmarshalJSON(data)
	assert.EqualError(t, err, expectedError)

	// Invalid Input data as string
	data = []byte{
		0x22, 0x88, 0x78, 0x74, 0x70, 0x3a, 0x2f, 0x2f, 0x65, 0x78,
		0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x22,
	}
	expectedError = `encoding of profile failed: `
	expectedError += `profile string must be an absolute URL or an ASN.1 OID: `
	expectedError += `failed to extract OID from string: `
	expectedError += `strconv.Atoi: parsing "�xtp://example": invalid syntax`
	err = profile.UnmarshalJSON(data)
	assert.EqualError(t, err, expectedError)
}

// TestProfile_MarshalJSON_OID_OK validates the JSON Marshaling of OID as profile
func TestProfile_MarshalJSON_OID_OK(t *testing.T) {
	inputOID := "2.5.2.1"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	expected := []byte(`"2.5.2.1"`)
	actual, err := profile.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

// TestProfile_UnmarshalJSON_OID_OK tests the JSON unmarshaling for OID as profile
func TestProfile_UnmarshalJSON_OID_OK(t *testing.T) {
	inputOID := "1.2.3.4"
	profile, err := NewProfile(inputOID)
	assert.Nil(t, err)
	input := []byte(`"2.5.2.1"`)
	err = profile.UnmarshalJSON(input)
	assert.Nil(t, err)
	expectedOID := "2.5.2.1"
	actualOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)

	// Partial OID test
	inputOID = ".2.3.4"
	_, err = NewProfile(inputOID)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to extract OID from string")
}

// TestProfile_RoundTrip_CBOR_Long_OID_OK, tests for a round trip encode/decode of
func TestProfile_RoundTrip_CBOR_Long_OID_OK(t *testing.T) {
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

	// Test Unmarshal CBOR
	input := actual
	expectedOID := inputOID
	err = profile.UnmarshalCBOR(input)
	assert.Nil(t, err)
	actualOID, err := profile.Get()
	assert.Nil(t, err)
	assert.Equal(t, expectedOID, actualOID)
}
