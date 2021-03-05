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

const (
	// ASN1AbsoluteOIDType represents the Type number of Absolute OID ASN Encoding
	ASN1AbsoluteOIDType = 0x06
	// ASN1LongLenMask is used to mask bit 8 of Length Indicator byte
	ASN1LongLenMask = 0x80
	// ASN1LenBytesMask is used to extract first 7 bits of Length Indicator byte
	ASN1LenBytesMask = 0x7F
	// MaxASN1OIDLen is the constant to adress the maximum OID length handled in system
	MaxASN1OIDLen = 260
	// MinNumOIDArcs represents the minimum required arcs for a valid OID
	MinNumOIDArcs = 3
)

// Profile is either an absolute URI (RFC3986) or an ASN.1 Object Identifier
type Profile struct {
	val interface{}
}

// NewProfile instantiates a Profile object from the given input string
// The string can either be an absolute URI or an ASN.1 Object Identifier
// in dotted-decimal notation
func NewProfile(urlOrOID string) (*Profile, error) {
	p := Profile{}
	if err := p.Set(urlOrOID); err != nil {
		return nil, err
	}
	return &p, nil
}

// Set sets the internal value of the Profile object to the given urlOrOID string
func (s *Profile) Set(urlOrOID string) error {
	return s.decodeProfileFromString(urlOrOID)
}

// Get returns the profile as string (URI or dotted-decimal OID)
func (s Profile) Get() (string, error) {
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
	_, ok := s.val.(*url.URL)
	return ok
}

// IsOID checks whether a stored profile is an OID
func (s Profile) IsOID() bool {
	_, ok := s.val.(asn1.ObjectIdentifier)
	return ok
}

// constructASN1fromVal constructs a TLV ASN1 byte array from an ASN1 value
// supplied as an input, assumption is a meaningful OID byte array of less than 255 bytes
func constructASN1fromVal(val []byte) ([]byte, error) {
	var OID [MaxASN1OIDLen]byte
	asn1OID := OID[:2]
	asn1OID[0] = ASN1AbsoluteOIDType
	if len(val) < 127 {
		asn1OID[1] = byte(len(val))
	} else if len(val) < 256 {
		// extra one byte is sufficient
		asn1OID[1] = 1 // Set to 1 to indicate one more byte carries the length
		asn1OID[1] = asn1OID[1] | ASN1LongLenMask
		asn1OID = append(asn1OID, byte(len(val)))
	} else {
		return nil, fmt.Errorf("excessively large OID not handled")
	}
	asn1OID = append(asn1OID, val...)
	return asn1OID, nil
}

// DecodeProfileCBOR decodes from a received CBOR data the profile
// as either a URL or a Object Identifier
func (s *Profile) decodeProfileCBOR(val interface{}) error {
	switch t := val.(type) {
	case string:
		u, err := url.Parse(t)
		if err != nil {
			return fmt.Errorf("profile URL parsing failed: %w", err)
		}
		if !u.IsAbs() {
			return fmt.Errorf("profile URL not in absolute form: %w", err)
		}
		s.val = u
	case []byte:
		var profileOID asn1.ObjectIdentifier
		val, err := constructASN1fromVal(t)
		if err != nil {
			return fmt.Errorf("could not construct valid ASN1 buffer from ASN1 value: %w", err)
		}
		rest, err := asn1.Unmarshal(val, &profileOID)
		if err != nil {
			return fmt.Errorf("malformed profile OID")
		}
		if len(rest) > 0 {
			return fmt.Errorf("ASN1 Unmarshal failed, as returned leftover %d bytes", len(rest))
		}
		if len(profileOID) < MinNumOIDArcs {
			return fmt.Errorf("CBOR decoding invalid, num arcs: %d < min OID arcs %d", len(profileOID), MinNumOIDArcs)
		}
		s.val = profileOID
	default:
		return fmt.Errorf("decoding failed unexpected type for profile: %T", t)
	}
	return nil
}

// extractASNValue removes Type and Len Bytes to generate a value component of encoded ASN
func extractASNValue(asn1OID []byte) ([]byte, error) {
	if asn1OID[0] != ASN1AbsoluteOIDType {
		return nil, fmt.Errorf("not a valid absoulute ASN1OID")
	}
	// offset to default TL bytes
	byteOffset := 2
	if asn1OID[1]&ASN1LongLenMask != 0 {
		byteOffset = byteOffset + int(asn1OID[1]&ASN1LenBytesMask)
	}
	return asn1OID[byteOffset:], nil
}

// MarshalCBOR encodes the Profile object as a CBOR text string (if it is a URL),
// or as CBOR byte string (if it is an ASN.1 OID)
func (s Profile) MarshalCBOR() ([]byte, error) {
	switch t := s.val.(type) {
	case *url.URL:
		return em.Marshal(t.String())

	case asn1.ObjectIdentifier:
		var asn1OID []byte
		asn1OID, err := asn1.Marshal(t)
		if err != nil {
			return nil, fmt.Errorf("ASN.1 encoding failed for OID: %w", err)
		}
		asn1OIDval, err := extractASNValue(asn1OID)
		if err != nil {
			return nil, fmt.Errorf("ASN.1 value extraction failed for OID: %w", err)
		}
		return em.Marshal(asn1OIDval)
	default:
		return nil, fmt.Errorf("invalid type for EAT profile")
	}
}

// UnmarshalCBOR attempts to initialize the Profile from the presented
// CBOR data. The data must be a text string, representing a URL
// or a byte array representing an Object Identifier
func (s *Profile) UnmarshalCBOR(data []byte) error {
	var val interface{}
	if len(data) == 0 {
		return fmt.Errorf("decoding of CBOR data failed: zero length data buffer")
	}
	if err := dm.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("CBOR decoding of profile failed: %w", err)
	}
	if err := s.decodeProfileCBOR(val); err != nil {
		return fmt.Errorf("invalid profile data: %w", err)
	}
	return nil
}

func decodeOIDfromString(val string) (asn1.ObjectIdentifier, error) {
	// Attempt to decode OID from received string
	var oid asn1.ObjectIdentifier
	if val == "" {
		return nil, fmt.Errorf("no valid OID")
	}

	for _, s := range strings.Split(val, ".") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to extract OID from string: %w", err)
		}
		oid = append(oid, n)
	}
	if len(oid) < MinNumOIDArcs {
		return nil, fmt.Errorf("invalid OID, num arcs: %d < min OID arcs %d", len(oid), MinNumOIDArcs)
	}
	return oid, nil
}

// decodeProfileFromString attempts to decode received string as a URL
// if it fails then attempts to decode as an OID.
func (s *Profile) decodeProfileFromString(val string) error {
	// First attempt to decode profile as a URL
	u, err := url.Parse(val)
	if err != nil || !u.IsAbs() {
		val, err := decodeOIDfromString(val)
		if err != nil {
			return fmt.Errorf("profile decode failed no valid URL or OID: %w", err)
		}
		s.val = val
	} else {
		s.val = u
	}
	return nil
}

// DecodeProfileJSON decodes a valid profile, from the received
// JSON string, mapping it to either a URL or an OID
func (s *Profile) decodeProfileJSON(val string) error {
	return s.decodeProfileFromString(val)
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

	if err := s.decodeProfileJSON(v); err != nil {
		return fmt.Errorf("ecoding of profile failed: %w", err)
	}
	return nil
}
