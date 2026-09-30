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
			DebugEnabledValue,
			nil,
		},
		{
			"disabled",
			DebugDisabledValue,
			nil,
		},
		{
			"disabled-since-boot",
			DebugDisabledSinceBootValue,
			nil,
		},
		{
			"permanent-disable",
			DebugPermanentDisableValue,
			nil,
		},
		{
			"full-permanent-disable",
			DebugFullPermanentDisableValue,
			nil,
		},
		{
			"out of range value",
			5,
			errors.New("invalid Debug value: 5"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d, err := DebugFromUint(test.tv)
			if test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, err)
			} else {
				assert.Equal(t, test.tv, d.Uint())
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
		DebugEnabledValue:              {[]byte{0x00}, `"enabled"`},
		DebugDisabledValue:             {[]byte{0x01}, `"disabled"`},
		DebugDisabledSinceBootValue:    {[]byte{0x02}, `"disabled-since-boot"`},
		DebugPermanentDisableValue:     {[]byte{0x03}, `"disabled-permanently"`},
		DebugFullPermanentDisableValue: {[]byte{0x04}, `"disabled-fully-and-permanently"`},
	}

	for codepoint, expected := range tests {
		d, err := DebugFromUint(codepoint)
		assert.NoError(t, err)

		actual, err := em.Marshal(d)
		assert.Nil(t, err)
		assert.Equal(t, expected.CBOR, actual)

		actual, err = json.Marshal(d)
		assert.Nil(t, err)
		assert.JSONEq(t, expected.JSON, string(actual))
	}
}
