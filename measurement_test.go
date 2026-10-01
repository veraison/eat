// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

var (
	measurementType   = 258
	measurementFormat = []byte{
		0xa4, 0x00, 0x63, 0x66, 0x6f, 0x6f, 0x0c, 0x01, 0x01, 0x63, 0x62, 0x61,
		0x72, 0x02, 0xa2, 0x18, 0x1f, 0x63, 0x62, 0x61, 0x7a, 0x18, 0x21, 0x82,
		0x01, 0x02,
	}

	// echo "[258, << {0:\"foo\",12:1,1:\"bar\",2:{31:\"baz\",33:[1,2]}} >> ]" | diag2cbor.rb | xxd -i
	encodedMeasurement = []byte{
		0x82, 0x19, 0x01, 0x02, 0x58, 0x1a, 0xa4, 0x00, 0x63, 0x66, 0x6f, 0x6f,
		0x0c, 0x01, 0x01, 0x63, 0x62, 0x61, 0x72, 0x02, 0xa2, 0x18, 0x1f, 0x63,
		0x62, 0x61, 0x7a, 0x18, 0x21, 0x82, 0x01, 0x02,
	}
)

func TestMeasurement_CBORMarshal_OK(t *testing.T) {
	m := Measurement{Type: measurementType, Format: measurementFormat}

	encoded, err := em.Marshal(m)
	assert.Nil(t, err)
	assert.Equal(t, encodedMeasurement, encoded)
}

func TestMeasurement_CBORUnmarshal_OK(t *testing.T) {
	var m Measurement
	assert.Nil(t, cbor.Unmarshal(encodedMeasurement, &m))
	assert.NotNil(t, m)
	assert.Equal(t, measurementType, m.Type)
	assert.Equal(t, measurementFormat, m.Format)
}
