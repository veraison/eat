// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/veraison/go-cose"
)

var (
	xBytesP256 = []byte{
		0x65, 0xed, 0xa5, 0xa1, 0x25, 0x77, 0xc2, 0xba, 0xe8, 0x29, 0x43, 0x7f, 0xe3, 0x38, 0x70, 0x1a,
		0x10, 0xaa, 0xa3, 0x75, 0xe1, 0xbb, 0x5b, 0x5d, 0xe1, 0x08, 0xde, 0x43, 0x9c, 0x08, 0x55, 0x1d,
	}
	yBytesP256 = []byte{
		0x1e, 0x52, 0xed, 0x75, 0x70, 0x11, 0x63, 0xf7, 0xf9, 0xe4, 0x0d, 0xdf, 0x9f, 0x34, 0x1b, 0x3d,
		0xc9, 0xba, 0x86, 0x0a, 0xf7, 0xe0, 0xca, 0x7c, 0xa7, 0xe9, 0xee, 0xcd, 0x00, 0x84, 0xd1, 0x9c,
	}
	kidBytesP256 = []byte{
		0x49, 0x6b, 0xd8, 0xaf, 0xad, 0xf3, 0x07, 0xe5, 0xb0, 0x8c, 0x64, 0xb0, 0x42, 0x1b, 0xf9, 0xdc,
		0x01, 0x52, 0x8a, 0x34, 0x4a, 0x43, 0xbd, 0xa8, 0x8f, 0xad, 0xd1, 0x66, 0x9d, 0xa2, 0x53, 0xec,
	}

	xBytesEd25519 = []byte{
		0xd7, 0x5a, 0x98, 0x01, 0x82, 0xb1, 0x0a, 0xb7, 0xd5, 0x4b, 0xfe, 0xd3, 0xc9, 0x64, 0x07, 0x3a,
		0x0e, 0xe1, 0x72, 0xf3, 0xda, 0xa6, 0x23, 0x25, 0xaf, 0x02, 0x1a, 0x68, 0xf7, 0x07, 0x51, 0x1a,
	}

	// need padding (start from 0x00)
	xLeadingZeroBytes = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	}
	yLeadingZeroBytes = []byte{
		0x00, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07,
		0x18, 0x29, 0x3A, 0x4B, 0x5C, 0x6D, 0x7E, 0x8F,
		0x90, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07,
		0x18, 0x29, 0x3A, 0x4B, 0x5C, 0x6D, 0x7E, 0x8F,
	}
	dLeadingZeroBytes = []byte{
		0x00, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07,
		0x18, 0x29, 0x3A, 0x4B, 0x5C, 0x6D, 0x7E, 0x8F,
		0x90, 0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07,
		0x18, 0x29, 0x3A, 0x4B, 0x5C, 0x6D, 0x7E, 0x8F,
	}
)

func TestCOSEKey_NewKeyFromPublic_NG(t *testing.T) {
	_, err := NewKeyFromPublic(0)
	assert.NotNil(t, err)

	var ecdsaPublicKey ecdsa.PublicKey
	_, err = NewKeyFromPublic(ecdsaPublicKey)
	assert.NotNil(t, err)

	var ed25519PublicKey ed25519.PublicKey
	_, err = NewKeyFromPublic(ed25519PublicKey)
	assert.NotNil(t, err)
}

func TestCOSEKey_NewKeyFromPrivate_NG(t *testing.T) {
	_, err := NewKeyFromPrivate(nil)
	assert.NotNil(t, err)

	var ecdsaPrivateKey ecdsa.PrivateKey
	_, err = NewKeyFromPrivate(ecdsaPrivateKey)
	assert.NotNil(t, err)

	var ed25519PrivateKey ed25519.PrivateKey
	_, err = NewKeyFromPrivate(ed25519PrivateKey)
	assert.NotNil(t, err)
}

func TestCOSEKey_FromECDSA_OK(t *testing.T) {
	pub := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(xLeadingZeroBytes),
		Y:     new(big.Int).SetBytes(yLeadingZeroBytes),
	}
	priv := ecdsa.PrivateKey{
		PublicKey: pub,
		D:         new(big.Int).SetBytes(dLeadingZeroBytes),
	}

	var key COSEKey
	err := key.FromECDSAPrivateKey(&priv)
	assert.Nil(t, err)
	assert.Equal(t, cose.KeyTypeEC2, key.Type)
	assert.Equal(t, cose.CurveP256, key.Crv)
	assert.Equal(t, xLeadingZeroBytes, key.X)
	assert.Equal(t, yLeadingZeroBytes, key.Y)
	assert.Equal(t, dLeadingZeroBytes, key.D)

	keyFromNew, err := NewKeyFromPrivate(&priv)
	assert.Nil(t, err)
	assert.Equal(t, key, *keyFromNew)

	key = COSEKey{}
	err = key.FromECDSAPublicKey(&pub)
	assert.Nil(t, err)
	assert.Equal(t, cose.KeyTypeEC2, key.Type)
	assert.Equal(t, cose.CurveP256, key.Crv)
	assert.Equal(t, xLeadingZeroBytes, key.X)
	assert.Equal(t, yLeadingZeroBytes, key.Y)

	keyFromNew, err = NewKeyFromPublic(&pub)
	assert.Nil(t, err)
	assert.Equal(t, key, *keyFromNew)
}

func TestCOSEKey_FromECDSA_NG(t *testing.T) {
	pub := ecdsa.PublicKey{
		Curve: elliptic.P224(),
		X:     new(big.Int).SetBytes(xLeadingZeroBytes),
		Y:     new(big.Int).SetBytes(yLeadingZeroBytes),
	}
	priv := ecdsa.PrivateKey{
		PublicKey: pub,
		D:         new(big.Int).SetBytes(dLeadingZeroBytes),
	}

	var key COSEKey
	err := key.FromECDSAPublicKey(&pub)
	assert.NotNil(t, err)
	err = key.FromECDSAPrivateKey(&priv)
	assert.NotNil(t, err)

	pub.Curve = elliptic.P256()

	concatenated := append(xLeadingZeroBytes, yLeadingZeroBytes...)
	pub.X = new(big.Int).SetBytes(concatenated)
	err = key.FromECDSAPublicKey(&pub)
	assert.NotNil(t, err)

	pub.X = new(big.Int).SetBytes(xLeadingZeroBytes)
	pub.Y = new(big.Int).SetBytes(concatenated)
	err = key.FromECDSAPublicKey(&pub)
	assert.NotNil(t, err)

	pub.Y = new(big.Int).SetBytes(yLeadingZeroBytes)
	priv.D = new(big.Int).SetBytes(concatenated)
	err = key.FromECDSAPrivateKey(&priv)
	assert.NotNil(t, err)
}

func TestCOSEKey_FromEd25519_OK(t *testing.T) {
	var key COSEKey
	err := key.FromEd25519KeyPair(dLeadingZeroBytes, xLeadingZeroBytes)
	assert.Nil(t, err)
	assert.Equal(t, cose.KeyTypeOKP, key.Type)
	assert.Equal(t, cose.CurveEd25519, key.Crv)
	assert.Equal(t, xLeadingZeroBytes, key.X)
	assert.Equal(t, dLeadingZeroBytes, key.D)

	concatenated := append(dLeadingZeroBytes, xLeadingZeroBytes...)
	priv := ed25519.PrivateKey(concatenated)
	keyFromNew, err := NewKeyFromPrivate(priv)
	assert.Nil(t, err)
	assert.Equal(t, key, *keyFromNew)

	key = COSEKey{}
	pub := ed25519.PublicKey(xLeadingZeroBytes)
	err = key.FromEd25519PublicKey(pub)
	assert.Nil(t, err)
	assert.Equal(t, cose.KeyTypeOKP, key.Type)
	assert.Equal(t, cose.CurveEd25519, key.Crv)
	assert.Equal(t, xLeadingZeroBytes, key.X)

	keyFromNew, err = NewKeyFromPublic(pub)
	assert.Nil(t, err)
	assert.Equal(t, key, *keyFromNew)
}

func TestCOSEKey_FromEd25519_NG(t *testing.T) {
	var key COSEKey
	err := key.FromEd25519PublicKey(nil)
	assert.NotNil(t, err)

	pub := ed25519.PublicKey(xBytesEd25519[:10])
	err = key.FromEd25519PublicKey(pub)
	assert.NotNil(t, err)

	err = key.FromEd25519KeyPair(nil, xBytesEd25519)
	assert.NotNil(t, err)

	err = key.FromEd25519KeyPair(dLeadingZeroBytes[:10], xBytesEd25519)
	assert.NotNil(t, err)
}

func TestCOSEKey_ThumbprintEC2_OK(t *testing.T) {
	// from section 6 of RFC 9679
	// {
	//     / kty set to EC2 = Elliptic Curve Keys /
	//     1: 2,
	//     / crv set to P-256 /
	//     -1: 1,
	//     / public key: x-coordinate /
	//     -2: h'65eda5a12577c2bae829437fe338701a10aaa375e1bb5b5de108de439c08551d',
	//     / public key: y-coordinate /
	//     -3: h'1e52ed75701163f7f9e40ddf9f341b3dc9ba860af7e0ca7ca7e9eecd0084d19c',
	//     / kid is bstr, not used in COSE Key Thumbprint /
	//     2: h'496bd8afadf307e5b08c64b0421bf9dc01528a344a43bda88fadd1669da253ec'
	// }
	key := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveP256,
		X:    xBytesP256,
		Y:    yBytesP256,
		ID:   kidBytesP256,
	}
	expected := []byte{
		0x49, 0x6b, 0xd8, 0xaf, 0xad, 0xf3, 0x07, 0xe5, 0xb0, 0x8c, 0x64, 0xb0, 0x42, 0x1b, 0xf9, 0xdc,
		0x01, 0x52, 0x8a, 0x34, 0x4a, 0x43, 0xbd, 0xa8, 0x8f, 0xad, 0xd1, 0x66, 0x9d, 0xa2, 0x53, 0xec,
	}
	thumbprint, err := key.Thumbprint(crypto.SHA256)
	assert.Nil(t, err)
	assert.Equal(t, expected, thumbprint)
}

func TestCOSEKey_ThumbprintEC2_NG(t *testing.T) {
	missingType := COSEKey{
		Crv: cose.CurveP256,
		X:   xBytesP256,
		Y:   yBytesP256,
	}
	_, err := missingType.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)
	_, err = missingType.Thumbprint(0)
	assert.NotNil(t, err)

	shortXP256 := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveP256,
		X:    xBytesP256[:31],
		Y:    yBytesP256,
	}
	_, err = shortXP256.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	shortXP384 := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveP384,
		X:    xBytesP256,
		Y:    yBytesP256,
	}
	_, err = shortXP384.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	shortXP521 := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveP521,
		X:    xBytesP256,
		Y:    yBytesP256,
	}
	_, err = shortXP521.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	missingY := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveP256,
		X:    xBytesP256,
		ID:   kidBytesP256,
	}
	_, err = missingY.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	invalidCrv := COSEKey{
		Type: cose.KeyTypeEC2,
		Crv:  cose.CurveEd25519,
		X:    xBytesP256,
		Y:    yBytesP256,
		ID:   kidBytesP256,
	}
	_, err = invalidCrv.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)
}

func TestCOSEKey_ThumbprintOKP_OK(t *testing.T) {
	key := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveEd25519,
		X:    xBytesEd25519,
	}
	_, err := key.Thumbprint(crypto.SHA256)
	assert.Nil(t, err)
	// no test vector found to be compared
}

func TestCOSEKey_ThumbprintOKP_NG(t *testing.T) {
	shortXEd25519 := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveEd25519,
		X:    xBytesEd25519[:31],
	}
	_, err := shortXEd25519.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	shortXEd448 := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveEd448,
		X:    xBytesEd25519,
	}
	_, err = shortXEd448.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	shortXX25519 := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveX25519,
		X:    xBytesEd25519[:31],
	}
	_, err = shortXX25519.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	shortXX448 := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveX448,
		X:    xBytesEd25519,
	}
	_, err = shortXX448.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	missingX := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveEd25519,
	}
	_, err = missingX.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)

	invalidCrv := COSEKey{
		Type: cose.KeyTypeOKP,
		Crv:  cose.CurveP256,
		X:    xBytesEd25519,
	}
	_, err = invalidCrv.Thumbprint(crypto.SHA256)
	assert.NotNil(t, err)
}

func TestCOSEKey_curveName(t *testing.T) {
	assert.Equal(t, "P-384", curveName(elliptic.P384()))
	assert.Equal(t, "P-521", curveName(elliptic.P521()))
	assert.Equal(t, "unknown", curveName(elliptic.P224()))
}

func TestCOSEKey_paddedBytes(t *testing.T) {
	xBigInt := new(big.Int).SetBytes(xBytesP256)
	short, err := paddedBytes(xBigInt, 32)
	assert.Nil(t, err)
	assert.Equal(t, 32, len(short))

	_, err = paddedBytes(xBigInt, 10)
	assert.NotNil(t, err)
}
