// Copyright 2021 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Profile is either a well formed RFC3986 compliant URI(uri  *url.URL).
// or an ObjectIdentifier  ([]int)
type Profile struct {
	val interface{}
}

// DecodeProfileCBOR decodes from a received CBOR data the profile
// as either a URL or a Object Identifier
func (s Profile) DecodeProfileCBOR(val interface{}) error {
	switch val.(type) {
	case string:
		var value string
		value = val.(string)
		if len(value) == 0 {
			return fmt.Errorf("no valid URL for profile")
		}
		lurl, err := url.Parse(value)
		if err != nil || !lurl.IsAbs() {
			return fmt.Errorf("profile URL parsing failed: %w", err)
		}
		s.val = lurl
	case []byte:
		var value []byte
		value = val.([]byte)
		var profileOID asn1.ObjectIdentifier
		_, err := asn1.Unmarshal(value, &profileOID)
		if err != nil {
			return fmt.Errorf("malformed profile OID detetced")
		}
		s.val = profileOID
	}
	return nil
}

// MarshalCBOR will encode the Profile value either as a CBOR text string(URL),
// or as byte array
func (s Profile) MarshalCBOR() ([]byte, error) {
	switch s.val.(type) {
	case *url.URL:
		var uri *url.URL
		uri = s.val.(*url.URL)
		return em.Marshal(uri.String())

	case asn1.ObjectIdentifier:
		var asn1OID []byte
		asn1OID, err := asn1.Marshal(s.val)
		if err != nil {
			return nil, fmt.Errorf("asn1 encoding failed for OID: %w", err)
		}
		return em.Marshal(asn1OID)
	default:
		return nil, fmt.Errorf("invalid type for eat profile")
	}
}

// UnmarshalCBOR attempts to initialize the Profile from the presented
// CBOR data. The data must be a text string, representing a URL
// or a byte array representing a Object Identifier
func (s *Profile) UnmarshalCBOR(data []byte) error {
	var val interface{}
	if len(data) == 0 {
		return fmt.Errorf("decoding of CBOR data failed: zero length data buffer")
	}
	if err := dm.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("decoding of CBOR data failed: %w", err)
	}
	if err := s.DecodeProfileCBOR(val); err != nil {
		return fmt.Errorf("invalid profile data: %w", err)
	}
	return nil
}

// DecodeProfileJSON decodes a valid profile, from the received
// JSON string, mapping it to either a URL or an OID
func (s Profile) DecodeProfileJSON(val string) error {
	// First attempt to decode profile as a URL
	u, err := url.Parse(val)
	if err != nil || !u.IsAbs() {
		// Now attempt to decode the same as OID
		var oid asn1.ObjectIdentifier
		for _, s := range strings.Split(val, ".") {
			num, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return fmt.Errorf("json decoding failed for profile: %w", err)
			}
			n := int(num)
			oid = append(oid, n)
		}
		s.val = oid
	} else {
		s.val = u
	}
	return nil
}

// MarshalJSON encodes the receiver Profile into a JSON string
func (s Profile) MarshalJSON() ([]byte, error) {
	// json interoperability oid -- encoded as a string using the well
	// established dotted-decimal notation (e.g., the text "1.2.250.1").
	switch s.val.(type) {
	case *url.URL:
		var uri *url.URL
		uri = s.val.(*url.URL)
		return json.Marshal(uri.String())
	case asn1.ObjectIdentifier:
		oid := s.val.(asn1.ObjectIdentifier)
		return json.Marshal(oid.String())
	default:
		return nil, fmt.Errorf("invalid profile type")
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
		return fmt.Errorf("failed to unMarshal JSON profile: %w", err)
	}
	return nil
}
