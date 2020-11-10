// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecLevel_Validate(t *testing.T) {
	assert := assert.New(t)

	var s1 SecurityLevel // default value
	assert.Nil(s1.Validate())

	s2 := SecurityLevel(SecLevelHardware)
	assert.Nil(s2.Validate())

	s3 := SecurityLevel(1337)
	assert.EqualError(s3.Validate(), "out of range value 1337 for SecurityLevel type")
}
