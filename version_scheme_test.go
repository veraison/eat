// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionScheme_CBORMarshal_OK(t *testing.T) {
	vs := VersionScheme(Multipartnumeric)
	encoded, err := vs.MarshalCBOR()
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x01}, encoded)
}

func TestVersionScheme_CBORUnmarshal_OK(t *testing.T) {
	var vs VersionScheme
	assert.Nil(t, vs.UnmarshalCBOR([]byte{0x01}))
	assert.Equal(t, VersionScheme(Multipartnumeric), vs)
}

func TestVersionScheme_CBORUnmarshal_NG(t *testing.T) {
	var vs VersionScheme
	assert.Nil(t, vs.UnmarshalCBOR([]byte{0x05}))
	assert.Equal(t, VersionScheme(5), vs)
	loc, ok := versionSchemeToString[vs]
	assert.NotNil(t, ok)
	assert.Equal(t, "", loc)

	assert.NotNil(t, vs.UnmarshalCBOR([]byte{0x41}))
}

func TestVersionScheme_JSONMarshal_OK(t *testing.T) {
	vs := VersionScheme(Multipartnumeric)
	encoded, err := vs.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"multipartnumeric"`), encoded)
}

func TestVersionScheme_JSONMarshal_NG(t *testing.T) {
	vs := VersionScheme(5)
	_, err := vs.MarshalJSON()
	assert.NotNil(t, err)
}

func TestVersionScheme_JSONUnmarshal_OK(t *testing.T) {
	var vs VersionScheme
	assert.Nil(t, vs.UnmarshalJSON([]byte(`"multipartnumeric"`)))
	assert.Equal(t, VersionScheme(Multipartnumeric), vs)
	assert.Nil(t, vs.UnmarshalJSON([]byte(`2`)))
	assert.Equal(t, VersionScheme(MultipartnumericSuffix), vs)
	assert.Nil(t, vs.UnmarshalJSON([]byte(`3`)))
	assert.Equal(t, VersionScheme(Alphanumeric), vs)
	assert.Nil(t, vs.UnmarshalJSON([]byte(`4`)))
	assert.Equal(t, VersionScheme(Decimal), vs)
	assert.Nil(t, vs.UnmarshalJSON([]byte(`16384`)))
	assert.Equal(t, VersionScheme(Semver), vs)
	assert.Nil(t, vs.UnmarshalJSON([]byte(`5`)))
	assert.Equal(t, VersionScheme(5), vs)
	assert.Equal(t, "", versionSchemeToString[vs])
	assert.Nil(t, vs.UnmarshalJSON([]byte(`-1`)))
	assert.Equal(t, VersionScheme(-1), vs)
	assert.Equal(t, "", versionSchemeToString[vs])
	assert.Nil(t, vs.UnmarshalJSON([]byte(`0`)))
	assert.Equal(t, VersionScheme(0), vs)
	assert.Equal(t, "", versionSchemeToString[vs])
}

func TestVersionScheme_JSONUnmarshal_NG(t *testing.T) {
	var vs VersionScheme
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`"unknown-scheme"`)))
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`1.2`)))
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`'`)))
	// exceeds IEEE 754 integer range, and it will be capped with 2^53
	assert.Nil(t, vs.UnmarshalJSON([]byte(`9007199254740993`))) // 2^53 + 1
	assert.NotEqual(t, VersionScheme(9007199254740993), vs)
}
