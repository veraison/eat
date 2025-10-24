// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
)

type Version struct {
	_       struct{} `cbor:",toarray"` // TODO: implement UnmarshalJSON
	Version string
	Scheme  *VersionScheme
}

type VersionScheme int

const (
	Multipartnumeric VersionScheme = iota + 1
	MultipartnumericSuffix
	Alphanumeric
	Decimal
	Semver = 16384
)

var versionSchemeToString = map[VersionScheme]string{
	Multipartnumeric:       "multipartnumeric",
	MultipartnumericSuffix: "multipartnumeric-suffix",
	Alphanumeric:           "alphanumeric",
	Decimal:                "decimal",
	Semver:                 "semver",
}

var stringToVersionScheme = map[string]VersionScheme{
	"multipartnumeric":        Multipartnumeric,
	"multipartnumeric-suffix": MultipartnumericSuffix,
	"alphanumeric":            Alphanumeric,
	"decimal":                 Decimal,
	"semver":                  Semver,
}

func (vs VersionScheme) MarshalJSON() ([]byte, error) {
	s, ok := versionSchemeToString[vs]
	if !ok {
		return nil, fmt.Errorf("invalid VersionScheme: %d", vs)
	}
	return json.Marshal(s)
}

func (vs *VersionScheme) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	loc, ok := stringToVersionScheme[s]
	if !ok {
		return fmt.Errorf("invalid VersionScheme string: %s", s)
	}
	*vs = loc
	return nil
}

func (vs VersionScheme) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(int(vs))
}

func (vs *VersionScheme) UnmarshalCBOR(data []byte) error {
	var i int
	if err := cbor.Unmarshal(data, &i); err != nil {
		return err
	}
	*vs = VersionScheme(i)
	return nil
}
