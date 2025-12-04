// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"fmt"
	"math/big"

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

func NewKeyFromPublic(pub crypto.PublicKey) (*COSEKey, error) {
	var key COSEKey
	switch vk := pub.(type) {
	case *ecdsa.PublicKey:
		err := key.FromECDSAPublicKey(vk)
		if err != nil {
			return nil, err
		}
		return &key, nil
	case ed25519.PublicKey:
		err := key.FromEd25519PublicKey(vk)
		if err != nil {
			return nil, err
		}
		return &key, nil
	default:
		return nil, fmt.Errorf("invalid key type: %v", pub)
	}
}

func NewKeyFromPrivate(priv crypto.PrivateKey) (*COSEKey, error) {
	var key COSEKey
	switch vk := priv.(type) {
	case *ecdsa.PrivateKey:
		err := key.FromECDSAPrivateKey(vk)
		if err != nil {
			return nil, err
		}
		return &key, nil
	case ed25519.PrivateKey:
		if len(vk) != 64 {
			return nil, fmt.Errorf("invalid length of private key: %d", len(vk))
		}
		err := key.FromEd25519KeyPair(vk[:32], ed25519.PublicKey(vk[32:]))
		if err != nil {
			return nil, err
		}
		return &key, nil
	default:
		return nil, fmt.Errorf("invalid key type: %v", priv)
	}
}

type thumbprintHandler func() ([]byte, error)

type CurveInfo struct {
	COSECurve cose.Curve
	KeySize   int
}

var curveInfos = map[string]CurveInfo{
	"P-256": {COSECurve: cose.CurveP256, KeySize: 32},
	"P-384": {COSECurve: cose.CurveP384, KeySize: 48},
	"P-521": {COSECurve: cose.CurveP521, KeySize: 66},
}

func curveName(c elliptic.Curve) string {
	switch c {
	case elliptic.P256():
		return "P-256"
	case elliptic.P384():
		return "P-384"
	case elliptic.P521():
		return "P-521"
	default:
		return "unknown"
	}
}

func (c *COSEKey) FromECDSAPublicKey(key *ecdsa.PublicKey) error {
	curve, ok := curveInfos[curveName(key.Curve)]
	if !ok {
		return fmt.Errorf("unknown curve for ECDSA: %d", key.Curve)
	}

	x, err := paddedBytes(key.X, curve.KeySize)
	if err != nil {
		return err
	}
	y, err := paddedBytes(key.Y, curve.KeySize)
	if err != nil {
		return err
	}

	c.Type = cose.KeyTypeEC2
	c.Crv = curve.COSECurve
	c.X = x
	c.Y = y
	return nil
}

func (c *COSEKey) FromECDSAPrivateKey(key *ecdsa.PrivateKey) error {
	pub := key.PublicKey
	curve, ok := curveInfos[curveName(pub.Curve)]
	if !ok {
		return fmt.Errorf("unknown curve for ECDSA: %d", key.Curve)
	}

	d, err := paddedBytes(key.D, curve.KeySize)
	if err != nil {
		return err
	}
	c.D = d
	return c.FromECDSAPublicKey(&pub)
}

func (c *COSEKey) FromEd25519PublicKey(pub ed25519.PublicKey) error {
	if pub == nil {
		return fmt.Errorf("invalid key: pub must not be nil")
	}
	if len(pub) != 32 {
		return fmt.Errorf("invalid key length: the length of pub must be 32")
	}
	c.Type = cose.KeyTypeOKP
	c.Crv = cose.CurveEd25519
	c.X = pub
	return nil
}

func (c *COSEKey) FromEd25519KeyPair(priv ed25519.PrivateKey, pub ed25519.PublicKey) error {
	if priv == nil || pub == nil {
		return fmt.Errorf("invalid key: priv and pub must not be nil")
	}
	if len(priv) != 32 || len(pub) != 32 {
		return fmt.Errorf("invalid key length: the length of priv and pub must be 32")
	}
	c.Type = cose.KeyTypeOKP
	c.Crv = cose.CurveEd25519
	c.X = pub
	c.D = priv
	return nil
}

// paddedBytes returns fixed-length bytes from big.Int
func paddedBytes(n *big.Int, size int) ([]byte, error) {
	b := n.Bytes()
	if len(b) > size {
		return nil, fmt.Errorf("integer too large for field size")
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b) // Zero-pad on the left (right-aligned)
	return padded, nil
}

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
