// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "fmt"

/*
debug-disable-type = &(
    not-disabled: 0,
    disabled: 1,
    disabled-since-boot: 2,
    permanent-disable: 3,
    full-permanent-disable: 4
)
*/
const (
	DebugNotDisabled = iota
	DebugDisabled
	DebugDisabledSinceBoot
	DebugPermanentDisable
	DebugFullPermanentDisable
)

// Debug models the debug-disable type
type Debug uint

// Set makes sure that the supplied val makes a good Debug claim
func (d *Debug) Set(val uint) error {
	switch val {
	case DebugNotDisabled:
	case DebugDisabled:
	case DebugDisabledSinceBoot:
	case DebugPermanentDisable:
	case DebugFullPermanentDisable:
	default:
		return fmt.Errorf("out of range value %v for Debug type", val)
	}

	*d = Debug(val)

	return nil
}
