// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import cose "github.com/veraison/go-cose"

type KeyConfirmation struct {
	Key *COSEKey `cbor:"1,keyasint,omitempty" json:"jwk,omitempty"`
	// TODO: EncryptedKey (currently go-cose doesn't support COSE_Encrypt0 / COSE_Encrypt)
	Kid           *[]byte `cbor:"3,keyasint,omitempty" json:"kid,omitempty"`
	KeyThumbprint *[]byte `cbor:"5,keyasint,omitempty" json:"jkt,omitempty"`
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
	Type      cose.KeyType   `cbor:"1,keyasint" json:"kty"`
	ID        []byte         `cbor:"2,keyasint,omitempty" json:"kid,omitempty"`
	Algorithm cose.Algorithm `cbor:"3:keyasint,omitempty" json:"alg,omitempty"`
	Ops       []cose.KeyOp   `cbor:"4,keyasint,omitempty" json:"ops,omitempty"`
	BaseIV    []byte         `cbor:"5,keyasint,omitempty"`

	// Additional parameter pairs for OKP and EC2.
	Crv cose.Curve `cbor:"-1,keyasint,omitempty" json:"crv,omitempty"`
	X   []byte     `cbor:"-2,keyasint,omitempty" json:"x,omitempty"`
	Y   []byte     `cbor:"-3,keyasint,omitempty" json:"y,omitempty"`
	D   []byte     `cbor:"-4,keyasint,omitempty" json:"d,omitempty"`
}
