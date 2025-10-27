// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	version                 = "1.3.4"
	versionMultipartNumeric = "1.3.4-beta"
	versionAlphanumeric     = "v1beta2"
	versionDecimal          = "134"
	scheme                  = Multipartnumeric

	// echo "[\"1.3.4\"]" | diag2cbor.rb | xxd -i
	encodedVersion = []byte{
		0x81, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34,
	}
	// echo "[\"1.3.4\",1]" | diag2cbor.rb | xxd -i
	encodedVersionMultipartNumeric = []byte{
		0x82, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34, 0x01,
	}
	// echo "[\"1.3.4\",h'']" | diag2cbor.rb | xxd -i
	encodedVersionByteScheme = []byte{
		0x82, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34, 0x40,
	}
	// echo "[\"1.3.4\",1,\"rc1\"]" | diag2cbor.rb | xxd -i
	encodedVersionLong = []byte{
		0x83, 0x65, 0x31, 0x2e, 0x33, 0x2e, 0x34, 0x01, 0x64, 0x72, 0x63, 0x31,
	}
	// echo "[]" | diag2cbor.rb | xxd -i
	encodedVersionShort = []byte{
		0x80,
	}
	// echo "[]" | diag2cbor.rb | xxd -i
	encodedVersionBroken = []byte{
		0x82,
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
	assert.Nil(t, v.UnmarshalCBOR(encodedVersion))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Nil(t, v.Scheme)

	assert.Nil(t, v.UnmarshalCBOR(encodedVersionMultipartNumeric))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, Multipartnumeric, *v.Scheme)
}

func TestVersion_CBORUnmarshal_NG(t *testing.T) {
	var v Version
	assert.NotNil(t, v.UnmarshalCBOR(encodedVersionByteScheme))
	assert.NotNil(t, v.UnmarshalCBOR(encodedVersionLong))
	assert.NotNil(t, v.UnmarshalCBOR(encodedVersionShort))
	assert.NotNil(t, v.UnmarshalCBOR(encodedVersionBroken))
}

func TestVersion_JSONMarshal_OK(t *testing.T) {
	v := Version{Version: version}
	encoded, err := json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `["1.3.4"]`, string(encoded))

	v = Version{Version: version, Scheme: &scheme}
	encoded, err = json.Marshal(v)
	assert.Nil(t, err)
	assert.Equal(t, `["1.3.4","multipartnumeric"]`, string(encoded))
}

func TestVersion_JSONUnmarshal_OK(t *testing.T) {
	var v Version
	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4"]`)))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Nil(t, v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4","multipartnumeric"]`)))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, Multipartnumeric, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4",1]`)))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, Multipartnumeric, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4-beta",2]`)))
	assert.NotNil(t, v)
	assert.Equal(t, versionMultipartNumeric, v.Version)
	assert.Equal(t, MultipartnumericSuffix, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["v1beta2",3]`)))
	assert.NotNil(t, v)
	assert.Equal(t, versionAlphanumeric, v.Version)
	assert.Equal(t, Alphanumeric, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["134",4]`)))
	assert.NotNil(t, v)
	assert.Equal(t, versionDecimal, v.Version)
	assert.Equal(t, Decimal, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4",16384]`)))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, Semver, *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`["1.3.4",100]`)))
	assert.NotNil(t, v)
	assert.Equal(t, version, v.Version)
	assert.Equal(t, VersionScheme(100), *v.Scheme)

	assert.Nil(t, v.UnmarshalJSON([]byte(`[""]`)))
	assert.NotNil(t, v)
	assert.Equal(t, "", v.Version)
	assert.Nil(t, v.Scheme)
}

func TestVersion_JSONUnmarshal_NG(t *testing.T) {
	var v Version
	assert.NotNil(t, v.UnmarshalJSON([]byte(`'`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`134`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`[]`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`[134]`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`["1.3.4","extra"]`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`["1.3.4",{}]`)))
	assert.NotNil(t, v.UnmarshalJSON([]byte(`["1.3.4",42,"extra"]`)))
}
