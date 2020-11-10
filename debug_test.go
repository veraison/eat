// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug_Validate(t *testing.T) {
	tests := []struct {
		name        string
		tv          uint
		expectedErr error
	}{
		{
			"not-disabled",
			DebugNotDisabled,
			nil,
		},
		{
			"disabled",
			DebugDisabled,
			nil,
		},
		{
			"disabled-since-boot",
			DebugDisabledSinceBoot,
			nil,
		},
		{
			"permanent-disable",
			DebugPermanentDisable,
			nil,
		},
		{
			"full-permanent-disable",
			DebugFullPermanentDisable,
			nil,
		},
		{
			"out of range value",
			5,
			errors.New("out of range value 5 for Debug type"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := Debug(test.tv)
			err := d.Validate()
			if test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, err)
			} else {
				assert.Equal(t, test.tv, uint(d))
			}
		})
	}
}

func TestDebug_Marshal(t *testing.T) {
	type Expected struct {
		CBOR []byte
		JSON string
	}

	tests := map[uint]Expected{
		DebugNotDisabled:          {[]byte{0x00}, `0`},
		DebugDisabled:             {[]byte{0x01}, `1`},
		DebugDisabledSinceBoot:    {[]byte{0x02}, `2`},
		DebugPermanentDisable:     {[]byte{0x03}, `3`},
		DebugFullPermanentDisable: {[]byte{0x04}, `4`},
		5:                         {[]byte{0x05}, `5`},
	}

	for codepoint, expected := range tests {
		d := Debug(codepoint)

		actual, err := em.Marshal(d)
		assert.Nil(t, err)
		assert.Equal(t, expected.CBOR, actual)

		actual, err = json.Marshal(d)
		assert.Nil(t, err)
		assert.JSONEq(t, expected.JSON, string(actual))
	}
}
