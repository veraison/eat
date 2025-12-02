// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

// $$Claims-Set-Claims //= ( intended-use-label => intended-use-type )

// intended-use-type = generic /
//                     registration /
//                     provisioning /
//                     csr /
//                     pop

// generic      = JC< "generic",      1 >
// registration = JC< "registration", 2 >
// provisioning = JC< "provisioning", 3 >
// csr          = JC< "csr",          4 >
// pop          = JC< "pop",          5 >

const (
	IntendedUseInvalid = iota
	IntendedUseGeneric
	IntendedUseRegistration
	IntendedUseProvisioning
	IntendedUseCSR // Certificate Issuance
	IntendedUsePoP // Proof of Possession
)

var intendedUseToString = map[IntendedUse]string{
	IntendedUseGeneric:      "generic",
	IntendedUseRegistration: "registration",
	IntendedUseProvisioning: "provisioning",
	IntendedUseCSR:          "csr",
	IntendedUsePoP:          "pop",
}

var stringToIntendedUse = map[string]IntendedUse{
	"generic":      IntendedUseGeneric,
	"registration": IntendedUseRegistration,
	"provisioning": IntendedUseProvisioning,
	"csr":          IntendedUseCSR,
	"pop":          IntendedUsePoP,
}

type IntendedUse uint

func (i IntendedUse) MarshalJSON() ([]byte, error) {
	s, ok := intendedUseToString[i]
	if !ok {
		return nil, fmt.Errorf("out of range value %v for intended-use-type", i)
	}
	return json.Marshal(s)
}

func (i *IntendedUse) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("expected string: %w", err)
	}

	v, ok := stringToIntendedUse[t]
	if !ok {
		return fmt.Errorf("invalid intended-use-type string %q", v)
	}

	*i = v
	return nil
}
