// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Submod is the type of a submod: either a raw EAT (wrapped in a Sign1 CWT), or
// a map of EAT claims
type Submod struct {
	value any
}

func (o *Submod) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

func (o *Submod) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("buffer too short")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 2, 3:
		// bstr or tstr -> CBOR-Nested-Token
		var token NestedToken
		if err := token.UnmarshalCBOR(data); err != nil {
			return err
		}

		o.value = &token
	case 4:
		// array -> Detached-Submodule-Digest
		var digest DetachedSubmodDigest
		if err := dm.Unmarshal(data, &digest); err != nil {
			return err
		}

		o.value = &digest
	case 5:
		// map -> Claims-Set
		var claimsSet ClaimsSet
		if err := claimsSet.UnmarshalCBOR(data); err != nil {
			return err
		}

		o.value = &claimsSet
	default:
		return fmt.Errorf("unexpected CBOR major type for submod: %d", majorType)
	}

	return nil
}

func (o *Submod) MarshalJSON() ([]byte, error) {
	toMarshal := o.value

	// JSON encoding does not allow a Submod to be a DetachedSubmodDigest, so we
	// need to wrap it inside a NestedToken
	if digest, ok := toMarshal.(*DetachedSubmodDigest); ok {
		data, err := digest.MarshalJSON()
		if err != nil {
			return nil, err
		}

		toMarshal = &NestedToken{
			Type: NestedTokenDigest,
			Data: data,
		}
	}

	return json.Marshal(toMarshal)
}

func (o *Submod) UnmarshalJSON(data []byte) error {
	text := strings.TrimSpace(string(data))
	if text == "" {
		return errors.New("truncated input")
	}

	switch text[0] {
	case '[':
		// array -> JSON-Selector -> Nested-Token
		var token NestedToken
		if err := token.UnmarshalJSON(data); err != nil {
			return err
		}

		o.value = &token
	case '{':
		// object -> Claims-Set
		var claimsSet ClaimsSet
		if err := claimsSet.UnmarshalJSON(data); err != nil {
			return err
		}

		o.value = &claimsSet
	default:
		return fmt.Errorf("invalid JSON type (expected object or array): %q", text)
	}

	return nil
}

// Submods is a collection of named submodule claims. Each submodule is either
// a ClaimsSet, a NestedToken, or a DetachedSubmodDigest.
type Submods map[string]*Submod

// NewSubmods returns a pointer to an empty Submods:
func NewSubmods() *Submods {
	ret := Submods(make(map[string]*Submod))
	return &ret
}

// AddClaimsSet adds the provided *ClaimsSet to the Submods under the specified
// name, overwritting any existing submod under that name. Pointer to the
// Submods is returned to allow chaing with other adds. Panics if the provided
// *ClaimsSet is nil.
func (o *Submods) AddClaimsSet(name string, claimsSet *ClaimsSet) *Submods {
	if claimsSet == nil {
		panic("nil claims set")
	}

	(*o)[name] = &Submod{claimsSet}
	return o
}

// AddDigest adds the provided *DetachedSubmodDigest to the Submods under the
// specified name, overwritting any existing submod under that name. Pointer to
// the Submods is returned to allow chaing with other adds. Panics if the
// provided *DetachedSubmodDigest is nil or is invalid.
func (o *Submods) AddDigest(name string, digest *DetachedSubmodDigest) *Submods {
	if digest == nil {
		panic("nil detached submod digest")
	}

	if err := digest.Validate(); err != nil {
		panic(err)
	}

	(*o)[name] = &Submod{digest}
	return o
}

// AddClaimsSet adds the provided *NestedToken to the Submods under the
// specified name, overwritting any existing submod under that name. Pointer to
// the Submods is returned to allow chaing with other adds. Panics if the
// provided *NestedToken is nil.
func (o *Submods) AddNestedToken(name string, token *NestedToken) *Submods {
	if token == nil {
		panic("nil nested token")
	}

	(*o)[name] = &Submod{token}
	return o
}

// Get retrieves a submod's value by name. If name is not in the collection,
// nil is returned.
func (o Submods) Get(name string) any {
	submod, ok := o[name]
	if !ok {
		return nil
	}

	return submod.value
}

// GetClaimsSet returns the *ClaimsSet associated with the provided name. If the
// name is not in Submods or the associated submod is not a *ClaimsSet, an
// error is returned.
func (o Submods) GetClaimsSet(name string) (*ClaimsSet, error) {
	value := o.Get(name)
	if value == nil {
		return nil, fmt.Errorf("no submod named %q", name)
	}

	claimsSet, ok := value.(*ClaimsSet)
	if !ok {
		return nil, fmt.Errorf("submod %q is not a claims set", name)
	}

	return claimsSet, nil
}

// GetDigest returns the *DetachedSubmodDigest associated with the provided
// name. If the name is not in Submods or the associated submod is not a
// *DetachedSubmodDigest, an error is returned.
func (o Submods) GetDigest(name string) (*DetachedSubmodDigest, error) {
	value := o.Get(name)
	if value == nil {
		return nil, fmt.Errorf("no submod named %q", name)
	}

	digest, ok := value.(*DetachedSubmodDigest)
	if !ok {
		return nil, fmt.Errorf("submod %q is not a detached submod digest", name)
	}

	return digest, nil
}

// GetNestedToken returns the *DetachedNestedToken associated with the provided
// name. If the name is not in Submods or the associated submod is not a
// *DetachedNestedToken, an error is returned.
func (o Submods) GetNestedToken(name string) (*NestedToken, error) {
	value := o.Get(name)
	if value == nil {
		return nil, fmt.Errorf("no submod named %q", name)
	}

	token, ok := value.(*NestedToken)
	if !ok {
		return nil, fmt.Errorf("submod %q is not a nested token", name)
	}

	return token, nil
}
