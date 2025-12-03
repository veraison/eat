// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

type KeyConfirmation struct {
	Key *COSEKey `cbor:"1,keyasint,omitempty" json:"jwk,omitempty"`
	// TODO: EncryptedKey (currently go-cose doesn't support COSE_Encrypt0 / COSE_Encrypt)
	Kid           *[]byte `cbor:"3,keyasint,omitempty" json:"kid,omitempty"`
	KeyThumbprint *[]byte `cbor:"5,keyasint,omitempty" json:"jkt,omitempty"`
}
