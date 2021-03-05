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
	// MaxASN1OIDLen is the maximum OID length accepted by the implementation
	MaxASN1OIDLen = 255
	// MinNumOIDArcs represents the minimum required arcs for a valid OID
	MinNumOIDArcs = 3
)

const (
	// asn1AbsolueOIDType is the type of an ASN.1 Object Identifier
	asn1AbsolueOIDType = 0x06
	// asn1LongLenMask is used to mask bit 8 of the length byte
	asn1LongLenMask = 0x80
	// asn1LenBytesMask is used to extract the first 7 bits from the length byte
	asn1LenBytesMask = 0x7F
)

// Profile is either an absolute URI (RFC3986) or an ASN.1 Object Identifier
type Profile struct {
	val interface{}
}

// NewProfile instantiates a Profile object from the given input string
// The string can either be an absolute URI or an ASN.1 Object Identifier
// in dotted-decimal notation. Relative Object Identifiers (e.g., .1.1.29) are
// not accepted.
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

// constructASN1fromVal constructs a TLV ASN.1 byte array from an ASN.1 value
// supplied as an input.  The supplied OID byte array must be upto
// MaxASN1OIDLen bytes.
func constructASN1fromVal(val []byte) ([]byte, error) {
	const maxTLOffset = 3
	var OID [MaxASN1OIDLen + maxTLOffset]byte
	asn1OID := OID[:2]
	asn1OID[0] = asn1AbsolueOIDType
	if len(val) < 127 {
		asn1OID[1] = byte(len(val))
	} else if len(val) <= MaxASN1OIDLen {
		// extra one byte is sufficient
		asn1OID[1] = 1 // Set to 1 to indicate one more byte carries the length
		asn1OID[1] |= asn1LongLenMask
		asn1OID = append(asn1OID, byte(len(val)))
	} else {
		return nil, fmt.Errorf("OIDs greater than %d bytes are not accepted", MaxASN1OIDLen)
	}
	asn1OID = append(asn1OID, val...)
	return asn1OID, nil
}

// decodeProfile decodes from a received CBOR data the profile
// as either a URL or a Object Identifier
func (s *Profile) decodeProfile(val interface{}) error {
	switch t := val.(type) {
	case string:
		u, err := url.Parse(t)
		if err != nil {
			return fmt.Errorf("profile URL parsing failed: %w", err)
		}
		if !u.IsAbs() {
			return fmt.Errorf("profile URL not in absolute form")
		}
		s.val = u
	case []byte:
		var profileOID asn1.ObjectIdentifier
		val, err := constructASN1fromVal(t)
		if err != nil {
			return fmt.Errorf("could not construct valid ASN.1 buffer from ASN.1 value: %w", err)
		}
		rest, err := asn1.Unmarshal(val, &profileOID)
		if err != nil {
			return fmt.Errorf("malformed profile OID")
		}
		if len(rest) > 0 {
			return fmt.Errorf("ASN.1 Unmarshal returned with %d leftover bytes", len(rest))
		}
		if len(profileOID) < MinNumOIDArcs {
			return fmt.Errorf("ASN.1 OID decoding failed: got %d arcs, expecting at least %d", len(profileOID), MinNumOIDArcs)
		}
		s.val = profileOID
	default:
		return fmt.Errorf("decoding failed: unexpected type for profile: %T", t)
	}
	return nil
}

// extractASNValue extracts the value component from the supplied ASN.1 OID
func extractASNValue(asn1OID []byte) ([]byte, error) {
	if asn1OID[0] != asn1AbsolueOIDType {
		return nil, fmt.Errorf("the supplied value is not an ASN.1 OID")
	}
	// offset to default TL bytes
	byteOffset := 2
	if asn1OID[1]&asn1LongLenMask != 0 {
		byteOffset += int(asn1OID[1] & asn1LenBytesMask)
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
// CBOR data. The data must be a CBOR text string, representing a URL
// or a CBOR byte string representing an Object Identifier
func (s *Profile) UnmarshalCBOR(data []byte) error {
	var val interface{}
	if len(data) == 0 {
		return fmt.Errorf("decoding of CBOR data failed: zero length data buffer")
	}
	if err := dm.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("CBOR decoding of profile failed: %w", err)
	}
	if err := s.decodeProfile(val); err != nil {
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
		return nil, fmt.Errorf("invalid OID: got %d arcs, expecting at least %d", len(oid), MinNumOIDArcs)
	}
	return oid, nil
}

// decodeProfileFromString attempts to decode the supplied string as a URL,
// if that fails, it then attempts to decode it as an OID.
func (s *Profile) decodeProfileFromString(val string) error {
	// First attempt to decode profile as a URL
	u, err := url.Parse(val)
	if err != nil || !u.IsAbs() {
		val, err := decodeOIDfromString(val)
		if err != nil {
			return fmt.Errorf("profile string must be an absolute URL or an ASN.1 OID: %w", err)
		}
		s.val = val
	} else {
		s.val = u
	}
	return nil
}

// decodeProfileJSON attempts at extracting an absolute URI or ASN.1 OID
// from the supplied JSON string
func (s *Profile) decodeProfileJSON(val string) error {
	return s.decodeProfileFromString(val)
}

// MarshalJSON encodes the receiver Profile into a JSON string
func (s Profile) MarshalJSON() ([]byte, error) {
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
		return fmt.Errorf("encoding of profile failed: %w", err)
	}
	return nil
}
