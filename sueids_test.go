// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSUEID_JSONMarshal_OK(t *testing.T) {
	expectedUEIDRAND := UEID{
		0x01, // RAND
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef, // 16 bytes
	}
	expectedUEIDEUI := UEID{
		0x02,                               // EUI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 6 bytes
	}
	expectedUEIDIMEI := UEID{
		0x03, // IMEI
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, 0xbe, 0xef,
		0xde, 0xad, 0xbe, 0xef, 0xde, 0xad, // 14 bytes
	}

	tv := []byte(`{
		"sueids": {
			"RAND": "Ad6tvu_erb7v3q2-796tvu8",
			"EUI": "At6tvu_erQ",
			"IMEI": "A96tvu_erb7v3q2-796t"
		}
	}`)

	var eat Eat

	err := json.Unmarshal(tv, &eat)
	assert.Nil(t, err)
	assert.NotNil(t, eat.SUEIDs)
	rand, exists := (*eat.SUEIDs)["RAND"]
	assert.True(t, exists)
	assert.Equal(t, expectedUEIDRAND, rand)
	eui, exists := (*eat.SUEIDs)["EUI"]
	assert.True(t, exists)
	assert.Equal(t, expectedUEIDEUI, eui)
	imei, exists := (*eat.SUEIDs)["IMEI"]
	assert.True(t, exists)
	assert.Equal(t, expectedUEIDIMEI, imei)
}
