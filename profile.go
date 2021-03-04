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

// NewProfile instantiates a new Profile object from a given input string
// The input string can either be a well formed URI as string
// OR it could be a OID formated as a string example "1.2.3.4"
func NewProfile(urlOrOID string) (*Profile, error) {
	p := Profile{}
	if err := p.SetProfile(urlOrOID); err != nil {
		return nil, err
	}
	return &p, nil
}

// SetProfile sets the given urlOrOID overwriting a previously stored
// value
func (s *Profile) SetProfile(urlOrOID string) error {
	// First attempt to decode input string as a URL
	u, err := url.Parse(urlOrOID)
	if err != nil || !u.IsAbs() {
		s.val, err = decodeOIDfromString(urlOrOID)
		if err != nil {
			return fmt.Errorf("no valid URI or OID supplied as an argument: %w", err)
		}
	} else {
		s.val = u
	}
	return nil
}

// GetProfile returns either a valid URI or OID as strings
// to the caller
func (s Profile) GetProfile() (string, error) {
	switch t := s.val.(type) {
	case *url.URL:
		return t.String(), nil
	case asn1.ObjectIdentifier:
		return t.String(), nil
	default:
		return "", fmt.Errorf("no valid EAT profile")
	}
}

// IsURI checks whether a stored profile is a URI
func (s Profile) IsURI() bool {
	switch s.val.(type) {
	case *url.URL:
		return true
	default:
		return false
	}
}

// IsOID checks whether a stored profile is an OID
func (s Profile) IsOID() bool {
	switch s.val.(type) {
	case asn1.ObjectIdentifier:
		return true
	default:
		return false
	}
}

// DecodeProfileCBOR decodes from a received CBOR data the profile
// as either a URL or a Object Identifier
func (s *Profile) DecodeProfileCBOR(val interface{}) error {
	switch t := val.(type) {
	case string:
		lurl, err := url.Parse(t)
		if err != nil {
			return fmt.Errorf("profile URL parsing failed: %w", err)
		}
		if !lurl.IsAbs() {
			return fmt.Errorf("profile URL not in absolute form: %w", err)
		}
		s.val = lurl
	case []byte:
		var profileOID asn1.ObjectIdentifier
		_, err := asn1.Unmarshal(t, &profileOID)
		if err != nil {
			return fmt.Errorf("malformed profile OID")
		}
		s.val = profileOID
	default:
		return fmt.Errorf("decoding failed malformed CBOR")
	}
	return nil
}

// MarshalCBOR will encode the Profile value either as a CBOR text string(URL),
// or as byte array
func (s Profile) MarshalCBOR() ([]byte, error) {
	switch t := s.val.(type) {
	case *url.URL:
		return em.Marshal(t.String())

	case asn1.ObjectIdentifier:
		var asn1OID []byte
		asn1OID, err := asn1.Marshal(t)
		if err != nil {
			return nil, fmt.Errorf("asn1 encoding failed for OID: %w", err)
		}
		return em.Marshal(asn1OID)
	default:
		return nil, fmt.Errorf("invalid type for EAT profile")
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

func decodeOIDfromString(val string) (asn1.ObjectIdentifier, error) {
	// Attempt to decode OID from received string
	var oid asn1.ObjectIdentifier
	for _, s := range strings.Split(val, ".") {
		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to extract OID from string: %w", err)
		}
		n := int(num)
		oid = append(oid, n)
	}
	return oid, nil
}

// DecodeProfileJSON decodes a valid profile, from the received
// JSON string, mapping it to either a URL or an OID
func (s *Profile) DecodeProfileJSON(val string) error {
	// First attempt to decode profile as a URL
	u, err := url.Parse(val)
	if err != nil || !u.IsAbs() {
		s.val, err = decodeOIDfromString(val)
		if err != nil {
			return fmt.Errorf("json decode of profile failed: %w", err)
		}
	} else {
		s.val = u
	}
	return nil
}

// MarshalJSON encodes the receiver Profile into a JSON string
func (s Profile) MarshalJSON() ([]byte, error) {
	// json interoperability oid -- encoded as a string using the well
	// established dotted-decimal notation (e.g., the text "1.2.250.1").
	switch t := s.val.(type) {
	case *url.URL:
		return json.Marshal(t.String())
	case asn1.ObjectIdentifier:
		return json.Marshal(t.String())
	default:
		return nil, fmt.Errorf("invalid profile type: %T", t)
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
