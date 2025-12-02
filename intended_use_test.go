// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	GenericUse = IntendedUse(IntendedUseGeneric)
	InvalidUse = IntendedUse(IntendedUseInvalid)
)

func TestIntendedUse_MarshalJSON_OK(t *testing.T) {
	t1, err := json.Marshal(GenericUse)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"generic"`), t1)

	var i1 IntendedUse
	err = json.Unmarshal([]byte(`"generic"`), &i1)
	assert.Nil(t, err)
	assert.Equal(t, GenericUse, i1)
}

func TestIntendedUse_MarshalJSON_NG(t *testing.T) {
	_, err := json.Marshal(InvalidUse)
	assert.NotNil(t, err)

	var i1 IntendedUse
	err = json.Unmarshal([]byte(`"invalid"`), &i1)
	assert.NotNil(t, err)

	var i2 IntendedUse
	err = json.Unmarshal([]byte(`1`), &i2)
	assert.NotNil(t, err)
}
