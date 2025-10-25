// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

var (
	version = "1.3.4"
	scheme  = Multipartnumeric

	encodedVersion = []byte{
		0x81, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34,
	}
	// echo "[\"1.3.4\",1]" | diag2cbor.rb | xxd -i
	encodedVersionMultipartNumeric = []byte{
		0x82, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34, 0x01,
	}
)

func TestVersion_CBORMarshal_OK(t *testing.T) {
	v := Version{Version: version}
	encoded, err := em.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, encodedVersion, encoded)

	v = Version{Version: version, Scheme: &scheme}
	encoded, err = em.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, encodedVersionMultipartNumeric, encoded)
}

func TestVersion_CBORUnmarshal_OK(t *testing.T) {
	var v Version
	assert.Nil(t, cbor.Unmarshal(encodedVersion, &v))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Nil(t, v.Scheme)

	assert.Nil(t, cbor.Unmarshal(encodedVersionMultipartNumeric, &v))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, Multipartnumeric, *v.Scheme)
}
