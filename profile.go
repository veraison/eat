// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"encoding/asn1"
	cbor "github.com/fxamacker/cbor/v2"
)

// Profile is either a well formed RFC3986 compliant URI.
// or an OID as a Byte String
type Profile struct {
	val interface{}
}

// FromString initializes the Profile value from the specified string,
// overwriting any existing value. If the value contains a colon (":"), then an
// attempt will be made to parse it as a URI (see RFC7519, section 2),
// otherwise, the value is assumed to be a non-URI string.
func (s *Profile) FromString(value string) error {
	if strings.Contains(value, ":") {
		u, err := url.Parse(value)
		if err != nil {
			return err
		}
		s.val = u
	} else {
		s.val = value
	}
	return nil
}

// FromURL initializes the Profile value from the specified url.URL,
// overwriting any existing value.
func (s *Profile) FromURL(value *url.URL) {
	s.val = value
}

// FromOID initializes the Profile value from the specified OID Value,
// overwriting any existing value
func (s *Profile) FromOID(value []byte) {
	s.val = &value
}

// IsURI returns true if the underlying value is a URI and not a oid
func (s Profile) IsURI() bool {
	return (s.val.(type) == string)
}

// IsOID returns true if the underlying value  is an OID and not a URI
func (s Profile) IsOID() bool {
	return (s.val.(type) == []byte)
}

// Validate checks for a valid profile
func (s Profile) Validate() error
{
	switch t := s.val.(type) {
	case string:
		if len(s.val) == 0 {
			return fmt.Errorf("No valid URL for profile")
		}
		_, err := url.Parse(s.val)
		if err != nil {
			return fmt.Errorf("Failed to parse profile URL: %w", err)
		}
	case []byte:
		if len(s.val) == 0 {
			return fmt.Errorf("No valid OID for profile")
		}
		var unmarshaldata interface{}
		_, err := asn1.Unmarshal(s.val, unmarshaldata)
		if (err != nil) {
			return fmt.Errorf("MalFormed profile OID detetced")
		}
	}
} 
// MarshalCBOR will encode the Profile value either as a CBOR text string(URL),
// or as byte array 
func (s Profile) MarshalCBOR() ([]byte, error) {
	if err := s.Validate() err != nil {
		return []byte{}, err 
	}
	return em.Marshal(s.val)
}

// UnmarshalCBOR attempts to initializes the Profile from the presented
// CBOR data. The data must be a text string, possibly wrapped in a Tag with
// the value 32 (URI). See RFC7049, Section 2.4.4.3.
func (s *Profile) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if err := dm.Unmarshal(data, &s.val); err != nil {
		return fmt.Errorf("Unmarshal of CBOR data failed: %w", err)
	}
	if err := s.Validate() err != nil {
		return  fmt.Errorf("Invalid profile data decoded: %w", err) 
	}
	return nil
}

// MarshalJSON encodes the receiver Profile into a JSON string
func (s Profile) MarshalJSON() ([]byte, error) {
	// json interoperability oid -- encoded as a string using the well
	// established dotted-decimal notation (e.g., the text "1.2.250.1").
	return json.Marshal(s.String())
}

// UnmarshalJSON attempts at decoding the supplied JSON data into the receiver
// Profile
func (s *Profile) UnmarshalJSON(data []byte) error {
	var v string

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if err := s.FromString(v); err != nil {
		return err
	}

	return nil
}
