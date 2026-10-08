// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/fxamacker/cbor/v2"
)

const (
	CWTTaggedMessageTag    uint64 = 61
	BUNDLETaggedMessageTag uint64 = 602
)

// NestedTokenType enumerates the valid types of the nested contents of a NestedToken
type NestedTokenType int8

const (
	NestedTokenInvalid NestedTokenType = iota
	NestedTokenCBOR
	NestedTokenJWT
	NestedTokenBundle
	NestedTokenDigest
)

func (o NestedTokenType) String() string {
	switch o {
	case NestedTokenCBOR:
		return "CBOR" // nolint:goconst
	case NestedTokenJWT:
		return "JWT"
	case NestedTokenBundle:
		return "BUNDLE"
	case NestedTokenDigest:
		return "DIGEST"
	default:
		panic(fmt.Sprintf("invalid token type: %d", int8(o)))
	}
}

// NestedTokenTypeFromString returns a NestedToknType corresponding to the
// provided string. An error is returned if the provided string does not
// correspond to a type.
func NestedTokenTypeFromString(typ string) (NestedTokenType, error) {
	typeMap := map[string]NestedTokenType{
		"CBOR":   NestedTokenCBOR,
		"JWT":    NestedTokenJWT,
		"BUNDLE": NestedTokenBundle,
		"DIGEST": NestedTokenDigest,
	}

	ret, ok := typeMap[typ]
	if ok {
		return ret, nil
	}

	return NestedTokenInvalid, fmt.Errorf("unknown nested token type %q", typ)
}

// NestedToken implements Nested-Token submod type as defined in
// https://www.rfc-editor.org/rfc/rfc9711.html#name-submods-submodules-claim
// A Nested-Token represents another EAT token nested as a submod inside an
// outer EAT token. The encoding of the inner token need not match the outer
// token.
type NestedToken struct {
	// Type identifes the format of Data
	Type NestedTokenType
	// Data contains the actual nested contents. Its format is determined by Type.
	Data []byte
}

// NewNestedToken returns a new NestedToken from the provided type and data. An
// error is returned if either the type or the data are invalid, or if they
// don't match.
func NewNestedToken(typ NestedTokenType, data []byte) (*NestedToken, error) {

	if typ == NestedTokenCBOR {
		if err := validateInnerCBORTags(data); err != nil {
			return nil, fmt.Errorf("cbor: %w", err)
		}
	} else {
		identifiedType, err := identifyInnerJSONType(data)
		if err != nil {
			return nil, err
		}

		if typ != identifiedType {
			return nil, fmt.Errorf("provided type is %q, but data appears to contain %q",
				typ.String(), identifiedType.String())
		}
	}

	return &NestedToken{
		Type: typ,
		Data: data,
	}, nil
}

// UnwrapDigest returns a pointer to a DetachedSubmodDigest constructed from
// the NestedToken's data. An error is returned if the token's type is not
// NestedTokenDigest or if the data contained invalid DetachedSubmodDigest
// JSON.
func (o *NestedToken) UnwrapDigest() (*DetachedSubmodDigest, error) {
	if o.Type != NestedTokenDigest {
		return nil, fmt.Errorf("token contains %s, not DIGEST", o.Type.String())
	}

	var digest DetachedSubmodDigest
	if err := digest.UnmarshalJSON(o.Data); err != nil {
		return nil, err
	}

	if err := digest.Validate(); err != nil {
		return nil, err
	}

	return &digest, nil
}

// UnwrapBundle returns a pointer to a DetachedEatBundle constructed from the
// NestedToken's data. An error is returned if the token's type is not
// NestedTokenBundle or if the data contained invalid DetachedEatBundle JSON.
func (o *NestedToken) UnwrapBundle() (*DetachedEatBundle, error) {
	if o.Type != NestedTokenBundle {
		return nil, fmt.Errorf("token contains %s, not BUNDLE", o.Type.String())
	}

	var bundle DetachedEatBundle
	if err := bundle.UnmarshalJSON(o.Data); err != nil {
		return nil, err
	}

	return &bundle, nil
}

func (o *NestedToken) MarshalCBOR() ([]byte, error) {
	switch o.Type {
	case NestedTokenCBOR:
		if err := validateInnerCBORTags(o.Data); err != nil {
			return nil, fmt.Errorf("inner CBOR: %w", err)
		}

		return em.Marshal(o.Data)
	case NestedTokenJWT, NestedTokenDigest, NestedTokenBundle:
		identifiedType, err := identifyInnerJSONType(o.Data)
		if err != nil {
			return nil, err
		}

		if o.Type != identifiedType {
			return nil, fmt.Errorf("Type is %q, but Data appears to contain %q", // nolint:staticcheck
				o.Type.String(), identifiedType.String())
		}

		return em.Marshal(string(o.Data))
	default:
		return nil, fmt.Errorf("invalid token type: %d", o.Type)
	}
}

func (o *NestedToken) UnmarshalCBOR(data []byte) error {
	if len(data) < 1 {
		return errors.New("truncated input")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 2:
		// bstr -> CBOR
		var decoded []byte
		if err := dm.Unmarshal(data, &decoded); err != nil {
			return err
		}

		if err := validateInnerCBORTags(decoded); err != nil {
			return err
		}

		o.Type = NestedTokenCBOR
		o.Data = decoded

		return nil
	case 3:
		// tstr -> JSON
		var encodedInner string
		if err := dm.Unmarshal(data, &encodedInner); err != nil {
			return err
		}

		typ, err := identifyInnerJSONType([]byte(encodedInner))
		if err != nil {
			return err
		}

		o.Type = typ
		o.Data = []byte(encodedInner)

		return nil
	default:
		return fmt.Errorf("unexpected major type %d", majorType)
	}
}

func (o *NestedToken) MarshalJSON() ([]byte, error) {
	var dataToEncode any
	switch o.Type {
	case NestedTokenCBOR:
		if err := validateInnerCBORTags(o.Data); err != nil {
			return nil, fmt.Errorf("inner CBOR: %w", err)
		}

		dataToEncode = base64.RawURLEncoding.EncodeToString(o.Data)
	case NestedTokenJWT, NestedTokenDigest, NestedTokenBundle:
		identifiedType, err := identifyInnerJSONType(o.Data)
		if err != nil {
			return nil, err
		}

		if identifiedType != o.Type {
			return nil, fmt.Errorf("Type is %q, but Data appears to contain %q", // nolint:staticcheck
				o.Type.String(), identifiedType.String())
		}

		dataToEncode = json.RawMessage(o.Data)
	default:
		return nil, fmt.Errorf("invalid token type: %d", o.Type)
	}

	toEncode := [2]any{o.Type.String(), dataToEncode}

	return json.Marshal(toEncode)
}

func (o *NestedToken) UnmarshalJSON(data []byte) error {
	var raw [2]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	var decodedType string
	if err := json.Unmarshal(raw[0], &decodedType); err != nil {
		return fmt.Errorf("selector type: %w", err)
	}

	if decodedType == "CBOR" {
		var encodedData string
		if err := json.Unmarshal(raw[1], &encodedData); err != nil {
			return err
		}

		decodedData, err := base64.RawURLEncoding.DecodeString(encodedData)
		if err != nil {
			return err
		}

		o.Type = NestedTokenCBOR
		o.Data = decodedData
	} else {
		identifiedType, err := identifyInnerJSONType(raw[1])
		if err != nil {
			return err
		}

		if decodedType != identifiedType.String() {
			return fmt.Errorf("selector type is %q but the nested-token appears to be %q",
				decodedType, identifiedType.String())
		}

		o.Type = identifiedType
		o.Data = raw[1]
	}

	return nil
}

func identifyInnerJSONType(data []byte) (NestedTokenType, error) {
	text := strings.TrimSpace(string(data))

	if len(text) < 2 {
		return NestedTokenInvalid, errors.New("malformed nested-token: must be at least 2 characters")
	}

	switch text[0] {
	case '[':
		// array -> Detached-Submodule-Digest or Detached-EAT-Bundle

		// both Digest and Bundle are two-element arrays
		var decoded [2]json.RawMessage

		if err := json.Unmarshal([]byte(text), &decoded); err != nil {
			return NestedTokenInvalid, fmt.Errorf("invalid digest or bundle: %w", err)
		}

		// if the first element is an int, then it is an algorithm id, and the
		// value is a Digest
		var test int
		if err := json.Unmarshal(decoded[0], &test); err == nil {
			return NestedTokenDigest, nil
		}

		return NestedTokenBundle, nil
	case '"':
		// string -> JWT (as we're assuming not CBOR)
		if text[len(text)-1] != '"' {
			return NestedTokenInvalid, errors.New(`no closing " detected`)
		}

		parts := strings.Split(text[1:len(text)-1], ".")
		if len(parts) != 3 {
			return NestedTokenInvalid, errors.New("invalid JWT")
		}

		for _, part := range parts {
			if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
				return NestedTokenInvalid, fmt.Errorf("invalid JWT: %w", err)
			}
		}

		return NestedTokenJWT, nil
	default:
		return NestedTokenInvalid, errors.New("malformed nested-token")
	}
}

func validateInnerCBORTags(data []byte) error {
	var tag cbor.RawTag
	if err := cbor.Unmarshal(data, &tag); err != nil {
		return err
	}

	if tag.Number != BUNDLETaggedMessageTag && tag.Number != CWTTaggedMessageTag {
		return fmt.Errorf("unexpected tag %d", tag.Number)
	}

	return nil
}
