// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

const DetachedEatBundleTag = 602

// DetachedEatBundle is a message that conveys an EAT (a NestedToken) plus
// detached claims sets setcured by that EAT.
type DetachedEatBundle struct {
	token  *NestedToken
	claims map[string]*ClaimsSet
}

// NewDetachedEatBundle returns a pointer to a new bundle that contains
// the provided NestedToken as its main token.
func NewDetachedEatBundle(mainToken *NestedToken) *DetachedEatBundle {
	if mainToken == nil {
		panic("main token is nil")
	}

	return &DetachedEatBundle{token: mainToken}
}

// MainToken returns a pointer to this bundle's main token.
func (o *DetachedEatBundle) MainToken() *NestedToken {
	return o.token
}

// DetachedClaimSets returns a map containing this bundle's detached claims
// sets.
func (o *DetachedEatBundle) DetachedClaimSets() map[string]*ClaimsSet {
	return o.claims
}

// SetDetachedClaimSets sets this bundle's detached claims sets to the provided
// map. A pointer to the bundle is returned to allow chaining.
func (o *DetachedEatBundle) SetDetachedClaimSets(dcs map[string]*ClaimsSet) *DetachedEatBundle {
	if len(dcs) == 0 {
		panic("empty detached claims sets map")
	}

	o.claims = dcs
	return o
}

// AddDetachedClaimSet adds the provided ClaimsSet to the bundle under the
// specified name. A pointer to the bundle is returned to allow chaining.
func (o *DetachedEatBundle) AddDetachedClaimSet(name string, claims *ClaimsSet) *DetachedEatBundle {
	if claims == nil {
		panic("nil claims")
	}

	if o.claims == nil {
		o.claims = make(map[string]*ClaimsSet)
	}

	o.claims[name] = claims
	return o
}

// ToTaggedMessageCBOR is like MarshalCBOR but wraps the resulting CBOR in tag 602.
func (o *DetachedEatBundle) ToTaggedMessageCBOR() ([]byte, error) {
	encodedBundle, err := o.MarshalCBOR()
	if err != nil {
		return nil, err
	}

	toEncode := cbor.RawTag{
		Number:  DetachedEatBundleTag,
		Content: encodedBundle,
	}

	return em.Marshal(toEncode)
}

func (o *DetachedEatBundle) MarshalCBOR() ([]byte, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	encodedClaimSets := make(map[string][]byte)
	for label, claimSet := range o.claims {
		encodedClaimSet, err := em.Marshal(claimSet)
		if err != nil {
			return nil, fmt.Errorf("claims set %q: %w", label, err)
		}

		encodedClaimSets[label] = encodedClaimSet
	}

	toEncode := [2]any{o.token, encodedClaimSets}
	return em.Marshal(toEncode)
}

func (o *DetachedEatBundle) UnmarshalCBOR(data []byte) error {
	if len(data) == 0 {
		return errors.New("truncated input")
	}

	majorType := (data[0] & 0xe0) >> 5
	if majorType == 6 {
		var tag cbor.RawTag
		if err := dm.Unmarshal(data, &tag); err != nil {
			return err
		}

		if tag.Number != DetachedEatBundleTag {
			return fmt.Errorf("invalid tag %d for DetachedEatBundle (must be %d)",
				tag.Number,
				DetachedEatBundleTag,
			)
		}

		data = tag.Content
	} else if majorType != 4 {
		return fmt.Errorf("invalid major type %d for DetachedEatBundle (must be 4 or 6)", majorType)
	}

	var outer [2]cbor.RawMessage

	if err := dm.Unmarshal(data, &outer); err != nil {
		return err
	}

	var nestedToken NestedToken
	if err := nestedToken.UnmarshalCBOR(outer[0]); err != nil {
		return fmt.Errorf("main token: %w", err)
	}

	o.token = &nestedToken

	encodedClaimSets := make(map[string][]byte)
	if err := dm.Unmarshal(outer[1], &encodedClaimSets); err != nil {
		return fmt.Errorf("detached claims sets: %w", err)
	}

	for label, encodedClaimSet := range encodedClaimSets {
		var claimSet ClaimsSet
		if err := dm.Unmarshal(encodedClaimSet, &claimSet); err != nil {
			return fmt.Errorf("claims set %q: %w", label, err)
		}

		o.AddDetachedClaimSet(label, &claimSet)
	}

	return o.validate()
}

func (o *DetachedEatBundle) MarshalJSON() ([]byte, error) {
	if err := o.validate(); err != nil {
		return nil, err
	}

	encodedClaimSets := make(map[string]string)
	for label, claimSet := range o.claims {
		jsonClaimSet, err := json.Marshal(claimSet)
		if err != nil {
			return nil, fmt.Errorf("claims set %q: %w", label, err)
		}

		encodedClaimSets[label] = base64.RawURLEncoding.EncodeToString(jsonClaimSet)
	}

	toEncode := [2]any{o.token, encodedClaimSets}
	return json.Marshal(toEncode)
}

func (o *DetachedEatBundle) UnmarshalJSON(data []byte) error {
	var outer [2]json.RawMessage

	if err := json.Unmarshal(data, &outer); err != nil {
		return err
	}

	var nestedToken NestedToken
	if err := nestedToken.UnmarshalJSON(outer[0]); err != nil {
		return fmt.Errorf("main token: %w", err)
	}

	o.token = &nestedToken

	encodedClaimsSets := make(map[string]string)
	if err := json.Unmarshal(outer[1], &encodedClaimsSets); err != nil {
		return fmt.Errorf("detached claims sets: %w", err)
	}

	for label, encodedClaimsSet := range encodedClaimsSets {
		jsonClaimsSet, err := base64.RawURLEncoding.DecodeString(encodedClaimsSet)
		if err != nil {
			return fmt.Errorf("claims set %q base64: %w", label, err)
		}

		var claimsSet ClaimsSet
		if err := json.Unmarshal(jsonClaimsSet, &claimsSet); err != nil {
			return fmt.Errorf("claims set %q: %w", label, err)
		}

		o.AddDetachedClaimSet(label, &claimsSet)
	}

	return o.validate()
}

func (o *DetachedEatBundle) validate() error {
	if o.token == nil {
		return errors.New("nil main token")
	}

	if len(o.claims) == 0 {
		return errors.New("no detached claims sets")
	}

	return nil
}
