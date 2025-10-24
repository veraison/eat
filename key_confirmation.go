// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "github.com/veraison/go-cose"

type KeyConfirmation struct {
	Key *cose.Key `cbor:"1,keyasint,omitempty" json:"jwk,omitempty"`
	// XXX: Correct? are there any appropriate type for COSE_Encrypt / COSE_Encrypt0?
	EncryptedCoseKey *[]byte `cbor:"2,keyasint,omitempty" json:"jwe,omitempty"`
	Kid              *[]byte `cbor:"3,keyasint,omitempty" json:"kid,omitempty"`
	KeyThumbprint    *[]byte `cbor:"5,keyasint,omitempty" json:"jkt,omitempty"`
}

/*
NOTE: supports only OKP and EC2 key

	COSE_Key = {
	    1 => tstr / int,          ; kty
	    ? 2 => bstr,              ; kid
	    ? 3 => tstr / int,        ; alg
	    ? 4 => [+ (tstr / int) ], ; key_ops
	    ? 5 => bstr,              ; Base IV
	    * label => values
	}
*/
type COSEKey struct {
	// Type identifies the family of keys for this structure, and thus,
	// which of the key-type-specific parameters need to be set.
	Type cose.KeyType `cbor:"1,keyasint" json:"kty"`
	// ID is the identification value matched to the kid in the message.
	ID []byte `cbor:"2,keyasint,omitempty" json:"kid,omitempty"`
	// Algorithm is used to restrict the algorithm that is used with the
	// key. If it is set, the application MUST verify that it matches the
	// algorithm for which the Key is being used.
	Algorithm cose.Algorithm `cbor:"3:keyasint,omitempty" json:"alg,omitempty"`
	// Ops can be set to restrict the set of operations that the Key is used for.
	Ops []cose.KeyOp `cbor:"4,keyasint,omitempty" json:"ops,omitempty"`
	// BaseIV is the Base IV to be xor-ed with Partial IVs.
	BaseIV []byte `cbor:"5,keyasint,omitempty"`

	// Any additional parameter (label,value) pairs.
	Crv cose.Curve `cbor:"-1,keyasint,omitempty" json:"crv,omitempty"`
	X   []byte     `cbor:"-2,keyasint,omitempty" json:"x,omitempty"`
	Y   []byte     `cbor:"-3,keyasint,omitempty" json:"y,omitempty"`
	D   []byte     `cbor:"-4,keyasint,omitmepty" json:"d,omitempty"`
}
