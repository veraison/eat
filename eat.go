// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

//nolint:staticcheck // json.Marshal triggers a warning because go-cose Key.Params uses map[any]any
package eat

import (
	"encoding/json"
)

// Eat is the internal representation of an Entity Attestation Token.
type Eat struct {

	// TODO: implement CWT/JWT and bundle handing here

	claims ClaimsSet
}

// Claims returns a pointer to the Eat's contained ClaimsSet
func (o *Eat) Claims() *ClaimsSet {
	return &o.claims
}

func (o *Eat) UnmarshalCBOR(data []byte) error {
	return dm.Unmarshal(data, &o.claims)
}

func (o *Eat) MarshalCBOR() ([]byte, error) {
	return em.Marshal(&o.claims)
}

func (o *Eat) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.claims)
}

func (o *Eat) MarshalJSON() ([]byte, error) {
	return json.Marshal(&o.claims)
}
