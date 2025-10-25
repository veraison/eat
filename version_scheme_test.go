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

func TestVersionScheme_JSONUnmarshal_OK(t *testing.T) {
	var vs VersionScheme
	assert.Nil(t, vs.UnmarshalJSON([]byte(`"multipartnumeric"`)))
	assert.Equal(t, VersionScheme(Multipartnumeric), vs)
	var vs2 VersionScheme
	assert.Nil(t, vs2.UnmarshalJSON([]byte(`2`)))
	assert.Equal(t, VersionScheme(MultipartnumericSuffix), vs2)
}

func TestVersionScheme_JSONUnmarshal_NG(t *testing.T) {
	var vs VersionScheme
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`"unknown-scheme"`)))
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`1.2`)))
	assert.NotNil(t, vs.UnmarshalJSON([]byte(`'`)))
}
