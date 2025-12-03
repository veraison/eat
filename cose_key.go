// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"crypto"
	"fmt"

	cose "github.com/veraison/go-cose"
)

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

type thumbprintHandler func() ([]byte, error)

//nolint:gocritic
func (k COSEKey) Thumbprint(hash crypto.Hash) ([]byte, error) {
	if !hash.Available() {
		return nil, fmt.Errorf("unsupported hash function: %d", hash)
	}

	handlers := map[cose.KeyType]thumbprintHandler{
		cose.KeyTypeOKP: k.calcOKPThumbprint,
		cose.KeyTypeEC2: k.calcEC2Thumbprint,
	}
	handler, ok := handlers[k.Type]
	if !ok {
		return nil, fmt.Errorf("unsupported kty for : %d", k.Type)
	}

	toBeHashedData, err := handler()
	if err != nil {
		return nil, err
	}

	h := hash.New()
	h.Write(toBeHashedData)
	return h.Sum(nil), nil
}

//nolint:gocritic
func (k COSEKey) calcOKPThumbprint() ([]byte, error) {
	if k.Type != cose.KeyTypeOKP {
		return nil, fmt.Errorf("invalid struct for OKP key: Type must be OKP (1)")
	}
	if k.X == nil {
		return nil, fmt.Errorf("invalid struct for OKP key: X must not be nil")
	}

	// NOTE: since Go's big.Int.Bytes() returns variable-length byte,
	// the following code not always pass.
	// Be sure to use (c COSEKey) FromEd25519PublicKey(key *ed25519.PublicKey)
	//
	// pub, priv, err := ed25519.GenerateKey(rand.Reader)
	// assert.Equal(t, 32, len(key.X.Bytes())
	switch k.Crv {
	case cose.CurveEd25519:
		if len(k.X) != 32 {
			return nil, fmt.Errorf("invalid struct for OKP Ed25519 key: the length of X must be 32")
		}
	case cose.CurveEd448:
		if len(k.X) != 57 {
			return nil, fmt.Errorf("invalid struct for OKP Ed448 key: the length of X must be 57")
		}
	case cose.CurveX25519:
		if len(k.X) != 32 {
			return nil, fmt.Errorf("invalid struct for OKP X25519 key: the length of X must be 32")
		}
	case cose.CurveX448:
		if len(k.X) != 56 {
			return nil, fmt.Errorf("invalid struct for OKP X448 key: the length of X must be 56")
		}
	default:
		return nil, fmt.Errorf("invalid curve: %d", k.Crv)
	}

	m := make(map[int]interface{})
	m[1] = k.Type
	m[-1] = k.Crv
	m[-2] = k.X
	return em.Marshal(m)
}

//nolint:gocritic
func (k COSEKey) calcEC2Thumbprint() ([]byte, error) {
	if k.Type != cose.KeyTypeEC2 {
		return nil, fmt.Errorf("invalid struct for EC2 key: Type must be EC2 (2)")
	}
	if k.X == nil || k.Y == nil {
		return nil, fmt.Errorf("invalid struct for EC2 key: X and Y must not be nil")
	}

	// NOTE: since Go's big.Int.Bytes() returns variable-length byte,
	// the following code not always pass.
	// Be sure to use (c COSEKey) FromECDSAPublicKey(key *ecdsa.PublicKey)
	//
	// key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	// assert.Equal(t, 66, len(key.X.Bytes())
	switch k.Crv {
	case cose.CurveP256:
		if len(k.X) != 32 || len(k.Y) != 32 {
			return nil, fmt.Errorf("invalid struct for EC2 P-256 key: the length of X and Y must be 32")
		}
	case cose.CurveP384:
		if len(k.X) != 48 || len(k.Y) != 48 {
			return nil, fmt.Errorf("invalid struct for EC2 P-384 key: the length of X and Y must be 48")
		}
	case cose.CurveP521:
		if len(k.X) != 66 || len(k.Y) != 66 {
			return nil, fmt.Errorf("invalid struct for EC2 P-521 key: the length of X and Y must be 66")
		}
	default:
		return nil, fmt.Errorf("invalid curve: %d", k.Crv)
	}

	m := make(map[int]interface{})
	m[1] = k.Type
	m[-1] = k.Crv
	m[-2] = k.X
	m[-3] = k.Y
	return em.Marshal(m)
}
