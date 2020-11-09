// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug_Set(t *testing.T) {
	tests := []struct {
		name        string
		tv          uint
		expectedErr error
	}{
		{
			"not-disabled",
			0,
			nil,
		},
		{
			"disabled",
			1,
			nil,
		},
		{
			"disabled-since-boot",
			2,
			nil,
		},
		{
			"permanent-disable",
			3,
			nil,
		},
		{
			"full-permanent-disable",
			4,
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
			var d Debug
			err := d.Set(test.tv)
			if test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, err)
			} else {
				assert.Equal(t, test.tv, uint(d))
			}
		})
	}
}
