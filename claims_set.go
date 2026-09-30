// Copyright 2020-2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ClaimsSet is the internal representation of a EAT Claim-Set
type ClaimsSet struct {
	Nonce *Nonce `cbor:"10,keyasint,omitempty" json:"eat_nonce,omitempty"`
	UEID  *UEID  `cbor:"256,keyasint,omitempty" json:"ueid,omitempty"`
	// TODO: support SUEIDs
	// TODO: support oemid-pem = int type
	OemID           *OemID      `cbor:"258,keyasint,omitempty" json:"oemid,omitempty"`
	HardwareModel   *BinaryData `cbor:"259,keyasint,omitempty" json:"hwmodel,omitempty"`
	HardwareVersion *Version    `cbor:"260,keyasint,omitempty" json:"hwversion,omitempty"`
	Uptime          *uint       `cbor:"261,keyasint,omitempty" json:"uptime,omitempty"`
	OemBoot         *bool       `cbor:"262,keyasint,omitempty" json:"oemboot,omitempty"`
	DebugStatus     *Debug      `cbor:"263,keyasint,omitempty" json:"dbgstat,omitempty"`
	Location        *Location   `cbor:"264,keyasint,omitempty" json:"location,omitempty"`
	Profile         *Profile    `cbor:"265,keyasint,omitempty" json:"eat-profile,omitempty"`
	Submods         *Submods    `cbor:"266,keyasint,omitempty" json:"submods,omitempty"`
	BootCount       *uint       `cbor:"267,keyasint,omitempty" json:"bootcount,omitempty"`
	BootSeed        *BinaryData `cbor:"268,keyasint,omitempty" json:"bootseed,omitempty"`
	// TODO: DLOAs
	SoftwareName    *StringOrURI   `cbor:"270,keyasint,omitempty" json:"swname,omitempty"`
	SoftwareVersion *Version       `cbor:"271,keyasint,omitempty" json:"swversion,omitempty"`
	Manifests       *[]Manifest    `cbor:"272,keyasint,omitempty" json:"manifests,omitempty"`
	Measurements    *[]Measurement `cbor:"273,keyasint,omitempty" json:"measurements,omitempty"`
	// TODO: MeasrementResults
	// TODO: IntendedUse

	// CWT/JWT claims
	Issuer     *string          `cbor:"1,keyasint,omitempty" json:"iss,omitempty"`
	Subject    *string          `cbor:"2,keyasint,omitempty" json:"sub,omitempty"`
	Audience   *Audience        `cbor:"3,keyasint,omitempty" json:"aud,omitempty"`
	Expiration *NumericDate     `cbor:"4,keyasint,omitempty" json:"exp,omitempty"`
	NotBefore  *NumericDate     `cbor:"5,keyasint,omitempty" json:"nbf,omitempty"`
	IssuedAt   *NumericDate     `cbor:"6,keyasint,omitempty" json:"iat,omitempty"`
	CwtID      *BinaryData      `cbor:"7,keyasint,omitempty" json:"cti,omitempty"`
	Cnf        *KeyConfirmation `cbor:"8,keyasint,omitempty" json:"cnf,omitempty"`

	privateClaims map[IntOrString]any `cbor:"-" json:"-"`
}

// NewClaimsSet returns a pointer to a new empty ClaimsSet
func NewClaimsSet() *ClaimsSet {
	return &ClaimsSet{
		privateClaims: make(map[IntOrString]any),
	}
}

// ToEat returns a new *Eat that wraps this ClaimsSet.
func (o *ClaimsSet) ToEat() *Eat {
	return &Eat{*o}
}

// SetPrivate sets the value of the specified private claims key. If the
// specified key clashes with one of the standard EAT claims, an error is
// returned.
func (o *ClaimsSet) SetPrivate(key IntOrString, value any) error {
	var standardClaims map[IntOrString]any
	if key.IsInt() {
		standardClaims = getValueMapCBOR(o)
	} else {
		standardClaims = getValueMapJSON(o)
	}

	for standardKey := range standardClaims {
		if key == standardKey {
			return fmt.Errorf("cannot set private claim: key %q is a standard claim", key.String())
		}
	}

	o.privateClaims[key] = value

	return nil
}

// GetPrivate returns the value of the specified private claims key. If the key
// has not been set, nil is returned.
func (o *ClaimsSet) GetPrivate(key IntOrString) any {
	return o.privateClaims[key]
}

func (o *ClaimsSet) MarshalCBOR() ([]byte, error) {
	toEncode := getValueMapCBOR(o)
	for key, value := range toEncode {
		if reflect.ValueOf(value).IsNil() {
			delete(toEncode, key)
		}
	}

	for privateLabel, privateValue := range o.privateClaims {
		toEncode[privateLabel] = privateValue
	}

	return em.Marshal(toEncode)
}

func (o *ClaimsSet) UnmarshalCBOR(data []byte) error {
	if err := dm.Unmarshal(data, (*claimSetUnmarshalHelper)(o)); err != nil {
		return err
	}

	var claims map[IntOrString]any
	if err := dm.Unmarshal(data, &claims); err != nil {
		return err
	}

	known := getValueMapCBOR(o)
	for label := range claims {
		if _, ok := known[label]; ok {
			delete(claims, label)
		}
	}

	o.privateClaims = claims

	return nil
}

func (o *ClaimsSet) MarshalJSON() ([]byte, error) {
	toEncode := make(map[string]any)
	claimsMap := getValueMapJSON(o)
	for key, value := range claimsMap {
		if !reflect.ValueOf(value).IsNil() {
			toEncode[key.String()] = value
		}
	}

	for privateLabel, privateValue := range o.privateClaims {
		toEncode[privateLabel.String()] = privateValue
	}

	return json.Marshal(toEncode)
}

func (o *ClaimsSet) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*claimSetUnmarshalHelper)(o)); err != nil {
		return err
	}

	var rawClaims map[string]any
	if err := json.Unmarshal(data, &rawClaims); err != nil {
		return err
	}

	known := getValueMapJSON(o)
	claims := make(map[IntOrString]any)
	for rawLabel, value := range rawClaims {
		label := IntOrStringFromString(rawLabel)
		if _, ok := known[label]; !ok {
			claims[label] = value
		}
	}

	o.privateClaims = claims

	return nil
}

type claimSetUnmarshalHelper ClaimsSet

type taggedValue struct {
	Tag   string
	Value any
}

func getValueMapJSON(source any) map[IntOrString]any {
	ret := make(map[IntOrString]any)

	for _, tv := range getTaggedValues(source, "json") {
		ret[IntOrString{tv.Tag}] = tv.Value
	}

	return ret
}

func getValueMapCBOR(source any) map[IntOrString]any {
	ret := make(map[IntOrString]any)

	for _, tv := range getTaggedValues(source, "cbor") {
		i, err := strconv.Atoi(tv.Tag)
		if err != nil {
			continue
		}

		ret[IntOrString{i}] = tv.Value
	}

	return ret
}

func getTaggedValues(source any, tagType string) []taggedValue {
	var ret []taggedValue

	structType := reflect.TypeOf(source)
	structVal := reflect.ValueOf(source)
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
		structVal = structVal.Elem()
	}

	for i := 0; i < structType.NumField(); i++ {
		tag := structType.Field(i).Tag

		tagText, ok := tag.Lookup(tagType)
		if !ok {
			continue
		}

		labelText := strings.Split(tagText, ",")[0]
		if labelText == "-" {
			continue
		}

		ret = append(ret, taggedValue{
			Tag:   labelText,
			Value: structVal.Field(i).Interface(),
		})
	}

	return ret
}
