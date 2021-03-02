// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"strconv"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"encoding/asn1"
	cbor "github.com/fxamacker/cbor/v2"
)

// Profile is either a well formed RFC3986 compliant URI(uri  *url.URL).
// or an ObjectIdentifier  ([]int)
type Profile struct {
	val interface{}
}

// DecodeProfileCBOR decodes from a received CBOR data the profile
// as either a URL or a Object Identifier
func (s Profile) DecodeProfileCBOR(val interface{}) error
{
	switch t := val.(type) {
	case string:
		if len(val) == 0 {
			return fmt.Errorf("No valid URL for profile")
		}
		s.val, err := url.Parse(s.val)
		if err != nil {
			return fmt.Errorf("Failed to parse profile URL: %w", err)
		}
	case []byte:
		if len(val) == 0 {
			return fmt.Errorf("No valid OID for profile")
		}
		_, err := asn1.Unmarshal(val, &s.val)
		if (err != nil) {
			return fmt.Errorf("MalFormed profile OID detetced")
		}
	}
} 
// MarshalCBOR will encode the Profile value either as a CBOR text string(URL),
// or as byte array 
func (s Profile) MarshalCBOR() ([]byte, error) {
	switch t := s.val.(type) {
	case *url.URL:
		var uri *url.URL
		uri = s.val.(*url.URL)
		return em.Marshal(uri.String())

	case asn1.ObjectIdentifier:
		var asn1_oid []byte

		asn1_oid, err := asn1.Marshal(s.val)
		if err != nil{
			return nil, fmt.Errorf("ASN1 Encoding failed for OID", err)
		}
		return em.Marshal(asn1_oid)
	}
	
}

// UnmarshalCBOR attempts to initialize the Profile from the presented
// CBOR data. The data must be a text string, representing a URL
// or a byte array representing a Object Identifier
func (s *Profile) UnmarshalCBOR(data []byte) error {
	var val interface{}
	if len(data) == 0 {
		return nil
	}
	if err := dm.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("Unmarshal of CBOR data failed: %w", err)
	}
	if err := s.DecodeProfileCBOR(val); err != nil {
		return  fmt.Errorf("Invalid profile data decoded: %w", err) 
	}
	return nil
}

// DecodeProfileJSON decodes a valid profile, from the received 
// JSON string, mapping it to either a URL or an OID
func (s Profile) DecodeProfileJSON(val string) error{
	// First attempt to decode profile as a URL
	u, err := url.Parse(val)
	if err != nil || u.IsAbs() {
		// Now attempt to decode the same as OID
		var oid asn1.ObjectIdentifier
		for _, s := range strings.Split(val, "."){
			num, err := strconv.ParseInt(s, 10, 32)
			if err != nil{
				return ftm.Errorf("JSON Decoding Failed for profile: %w", err)
			}
			n := int32(num)
			oid = append(oid, n)
		}
		s.val = oid
	}
	else {
		s.val = val
	}
	return nil
}
// MarshalJSON encodes the receiver Profile into a JSON string
func (s Profile) MarshalJSON() ([]byte, error) {
	// json interoperability oid -- encoded as a string using the well
	// established dotted-decimal notation (e.g., the text "1.2.250.1").
	switch t := s.val.(type) {
	case *url.URL:
		var uri *url.URL
		uri = s.val.(*url.URL)
		return json.Marshal(uri.String())
	case asn1.ObjectIdentifier:
		oid := s.val.(asn1.ObjectIdentifier)
		return json.Marshal(oid.String())
	default:
		return nil, fmt.Errorf("Invalid Profile Type")
	}
}

// UnmarshalJSON attempts at decoding the supplied JSON data into the receiver
// Profile
func (s *Profile) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if err := s.DecodeProfileJSON(v); err != nil {
		return fmt.Errorf("Failed to UnMarshal JSON Profile", err)
	}
	return nil
}
