// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// COSE Algorithm Registry
// https://www.iana.org/assignments/cose#algorithms
const (
	Shake256   int = -45
	Sha512     int = -44
	Sha384     int = -43
	Shake128   int = -18
	Sha512_256 int = -17
	Sha256     int = -16
	Sha256_64  int = -15
	Sha1       int = -14
)

var (
	HashAlgorithmShake256   = MustHashAlgorithmFromInt(Shake256)
	HashAlgorithmSha512     = MustHashAlgorithmFromInt(Sha512)
	HashAlgorithmSha384     = MustHashAlgorithmFromInt(Sha384)
	HashAlgorithmShake128   = MustHashAlgorithmFromInt(Shake128)
	HashAlgorithmSha512_256 = MustHashAlgorithmFromInt(Sha512_256)
	HashAlgorithmSha256     = MustHashAlgorithmFromInt(Sha256)
	HashAlgorithmSha256_64  = MustHashAlgorithmFromInt(Sha256_64)
	HashAlgorithmSha1       = MustHashAlgorithmFromInt(Sha1)

	algToValueLen = map[int]int{
		Shake256:   64,
		Sha512:     64,
		Sha384:     48,
		Shake128:   32,
		Sha512_256: 32,
		Sha256:     32,
		Sha256_64:  8,
		Sha1:       20,
	}

	algToString = map[int]string{
		Shake256:   "shake-256",
		Sha512:     "sha-512",
		Sha384:     "sha-384",
		Shake128:   "shake-128",
		Sha512_256: "sha-512-256",
		Sha256:     "sha-256",
		Sha256_64:  "sha-256-64",
		Sha1:       "sha1",
	}

	stringToAlg = map[string]int{
		"shake-256":   Shake256,
		"sha-512":     Sha512,
		"sha-384":     Sha384,
		"shake-128":   Shake128,
		"sha-512-256": Sha512_256,
		"sha-256":     Sha256,
		"sha-256-64":  Sha256_64,
		"sha1":        Sha1,
	}
)

// MustHashAlgorithmFromInt returns the HashAlgorithm corresponding to the
// provided int (as HashAlgorithmFromInt). It panics if the provided int does
// not correspond to an algorithm.
func MustHashAlgorithmFromInt(val int) HashAlgorithm {
	ret, err := HashAlgorithmFromInt(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// HashAlgorithmFromInt returns the HashAlgorithm corresponding to the provided
// int. An error is returned is the provided int does not correspond to an
// algorithm inside the COSE Algorithm Registry:
// https://www.iana.org/assignments/cose#algorithms
func HashAlgorithmFromInt(val int) (HashAlgorithm, error) {
	_, ok := algToString[val]
	if ok {
		return HashAlgorithm{val}, nil
	}

	return HashAlgorithm{}, fmt.Errorf("invalid hash algorithm: %d", val)
}

// MustHashAlgorithmFromString converts an informal algorithm name to the
// corresponding HashAlgorithm (as HashAlgorithmFromString). It panics on
// error.
func MustHashAlgorithmFromString(val string) HashAlgorithm {
	ret, err := HashAlgorithmFromString(val)
	if err != nil {
		panic(err)
	}

	return ret
}

// HashAlgorithmFromInt returns the HashAlgorithm corresponding to the provided
// string. At the time of writing, no algorithm inside the COSE algorithm
// register has a string identifier. Currently, this function returns
// algorithms corresponding the the informal name (as obtained by calling
// HashAlgorithm.String()). An error is returned is the provided string does
// not correspond to an algorithm.
func HashAlgorithmFromString(val string) (HashAlgorithm, error) {
	i, ok := stringToAlg[val]
	if ok {
		return HashAlgorithm{i}, nil
	}

	return HashAlgorithm{}, fmt.Errorf("invalid hash algorithm: %q", val)
}

// HashAlgorithmFromAny converts the specfied value into corresponding
// HashAlgorithm. If the value is an int, int64, or float64 (without a
// fractional part), it behaves as HashAlgorithmFromInt; if the value is a
// string, it behaves as HashAlgorithmFromString; an error is returned in all
// other cases.
func HashAlgorithmFromAny(val any) (HashAlgorithm, error) {
	switch t := val.(type) {
	case int:
		return HashAlgorithmFromInt(t)
	case int64:
		return HashAlgorithmFromInt(int(t))
	case float64:
		return HashAlgorithmFromInt(int(t))
	case string:
		return HashAlgorithmFromString(t)
	default:
		return HashAlgorithm{}, fmt.Errorf("invalid hash algorithm: %v(%T)", t, t)
	}
}

// HashAlgorithm identifies the hashing algorithm used by a
// DetachedSubmodDigest.
type HashAlgorithm struct {
	val any
}

// IsString true if the underlying identifier is a string. (Note: all current
// COSE Algorithm Regiter entries are identified by integer values, this always
// returns false)
func (o HashAlgorithm) IsString() bool {
	_, ok := o.val.(string)
	return ok
}

// IsInt true if the underlying identifier is an int. (Note: all current COSE
// Algorithm Regiter entries are identified by integer values, this always
// returns true)
func (o HashAlgorithm) IsInt() bool {
	_, ok := o.val.(int)
	return ok
}

// String returns the string representation of the HashAlgorithm. If the
// underlying identifier is a string (currently, there are none), that
// idnetifier is returned, otherwise an informal name corresponding to the
// identifier is returned.
func (o HashAlgorithm) String() string {
	switch t := o.val.(type) {
	case string:
		return t
	case int:
		text, ok := algToString[t]
		if !ok {
			text = fmt.Sprintf("%d", t)
		}

		return text
	default:
		return ""
	}
}

// Int returns the integer identifier corresponding to the HashAlgorithm. If
// the underlying identifier is a string (currently, there are none), 0 is
// returned.
func (o HashAlgorithm) Int() int {
	switch t := o.val.(type) {
	case int:
		return t
	case string:
		val := stringToAlg[t]
		return val
	default:
		return 0
	}
}

func (o HashAlgorithm) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.val)
}

func (o *HashAlgorithm) UnmarshalCBOR(data []byte) error { // nolint:dupl
	if len(data) == 0 {
		return errors.New("buffer too short")
	}

	majorType := (data[0] & 0xe0) >> 5
	switch majorType {
	case 0, 1:
		var val int
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	case 3:
		var val string
		if err := dm.Unmarshal(data, &val); err != nil {
			return err
		}

		o.val = val
		return nil
	default:
		return fmt.Errorf("unexpected CBOR major type for DigestAlgID: %d", majorType)
	}
}

func (o HashAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.val)
}

func (o *HashAlgorithm) UnmarshalJSON(data []byte) error {
	var val any
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	switch t := val.(type) {
	case float64:
		alg, err := HashAlgorithmFromInt(int(t))
		if err != nil {
			return err
		}

		*o = alg
	default:
		return fmt.Errorf("unexpected algorithm value: %v(%T)", t, t)
	}

	return nil
}

// DetachedSubmodDigest implements Detached-Submodule-Digest as defined in
// https://www.rfc-editor.org/rfc/rfc9711.html#name-submods-submodules-claim
// Detached-Submodule-Digest contains a digest of a Claim-Set detached from an
// EAT.
type DetachedSubmodDigest struct {
	_         struct{} `cbor:",toarray"`
	Algorithm HashAlgorithm
	Value     []byte
}

// DetachedSubmodDigestFromString returns a DetachedSubmodDigest costructed
// from the provided string formatted as
// "<algorithm-name>;<base64-encoded-digest>". This is an inviser of
// DetachedSubmodDigest.String().
func DetachedSubmodDigestFromString(val string) (*DetachedSubmodDigest, error) {
	parts := strings.Split(val, ";")
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected exactly two ;-separated parts, got %q", val)
	}

	value, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("val: %w", err)
	}

	alg, err := HashAlgorithmFromString(parts[0])
	if err != nil {
		return nil, err
	}

	return NewDetachedSubmodDigest(alg, value)
}

// MustNewDetachedSubmodDigestIntAlg creates a new DetachedSubmodDigest
// as MustNewDetachedSubmodDigestIntAlg, except it panics on error.
func MustNewDetachedSubmodDigestIntAlg(algID int, value []byte) *DetachedSubmodDigest {
	ret, err := NewDetachedSubmodDigestIntAlg(algID, value)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewDetachedSubmodDigestIntAlg returns a new DetatchedSubmodDigest
// constructed from the provided int aglorithm idnetifier and []byte data. An
// error is returned if the identifier does not correspond to an algorithm from
// the COSE Algorithm Registry, or if the length of data is inconsitent with
// the specified algorithm.
func NewDetachedSubmodDigestIntAlg(algID int, value []byte) (*DetachedSubmodDigest, error) {
	alg, err := HashAlgorithmFromInt(algID)
	if err != nil {
		return nil, err
	}

	return NewDetachedSubmodDigest(alg, value)
}

// MustNewDetachedSubmodDigest creates a new DetachedSubmodDigest as
// NewDetachedSubmodDigest, except it panics on error.
func MustNewDetachedSubmodDigest(alg HashAlgorithm, value []byte) *DetachedSubmodDigest {
	ret, err := NewDetachedSubmodDigest(alg, value)
	if err != nil {
		panic(err)
	}

	return ret
}

// NewDetachedSubmodDigest returns a new DetachedSubmodDigest constructed from
// the provided hash algorithm and []byte value. An error is returned if the
// length of the value is not consistent with the specified algorithm.
func NewDetachedSubmodDigest(alg HashAlgorithm, value []byte) (*DetachedSubmodDigest, error) {
	ret := &DetachedSubmodDigest{Algorithm: alg, Value: value}
	if err := ret.Validate(); err != nil {
		return nil, err
	}
	return &DetachedSubmodDigest{Algorithm: alg, Value: value}, nil
}

// String returns a string representation of the DetachedSubmodDigest. This
// representation is in the form "<algorithm-name>;<base64-encoded-digest>"
// (consistent with https://datatracker.ietf.org/doc/rfc6920).
func (o DetachedSubmodDigest) String() string {
	return o.Algorithm.String() + ";" + base64.RawURLEncoding.EncodeToString(o.Value)
}

// Validate returns an error if either the algorithm of the value have not been
// initialized or if the length of the value does not match that expected by
// the algorithm.
func (o DetachedSubmodDigest) Validate() error {
	if len(o.Value) == 0 {
		return errors.New("zero length value")
	}

	if o.Algorithm.Int() == 0 && !o.Algorithm.IsString() {
		return errors.New("zero algorithm")
	}

	wantLen, ok := algToValueLen[o.Algorithm.Int()]
	if ok {
		gotLen := len(o.Value)
		if wantLen != gotLen {
			return fmt.Errorf(
				"length mismatch for hash algorithm %s: want %d bytes, got %d",
				o.Algorithm.String(), wantLen, gotLen,
			)
		}
	}

	return nil
}

func (o DetachedSubmodDigest) MarshalJSON() ([]byte, error) {
	toMarshal := [2]any{
		o.Algorithm,
		base64.RawURLEncoding.EncodeToString(o.Value),
	}

	return json.Marshal(toMarshal)
}

func (o *DetachedSubmodDigest) UnmarshalJSON(data []byte) error {
	var decoded []any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	if len(decoded) != 2 {
		return fmt.Errorf("expected array with two elements, got %v", decoded)
	}

	alg, err := HashAlgorithmFromAny(decoded[0])
	if err != nil {
		return fmt.Errorf("alg: %w", err)
	}

	switch t := decoded[1].(type) {
	case string:
		bytes, err := base64.RawURLEncoding.DecodeString(t)
		if err != nil {
			return fmt.Errorf("val: %w", err)
		}

		o.Algorithm = alg
		o.Value = bytes

		return nil
	default:
		return fmt.Errorf("invalid val: %v(%T)", t, t)
	}

}
