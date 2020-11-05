// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEat_ToCBOR(t *testing.T) {
	tv := Eat{}
	expected := []byte{0xa0}

	actual, err := tv.ToCBOR()

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestEat_FromCBOR(t *testing.T) {
	tv := []byte{0xa0}
	expected := Eat{}

	var actual Eat
	err := expected.FromCBOR(tv)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
