// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"fmt"
)

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
	// DebugNotDisabled is asserted if any debug facility, even manufacturer
	// hardware diagnostics, is currently enabled
	DebugNotDisabled = iota

	// DebugDisabled indicates all debug facilities are currently disabled. It
	// may be possible to enable them in the future, and it may also be possible
	// that they were enabled in the past after the target device/sub-system
	// booted/started, but they are currently disabled.
	DebugDisabled

	// DebugDisabledSinceBoot indicates all debug facilities are currently
	// disabled and have been so since the target device/sub-system
	// booted/started.
	DebugDisabledSinceBoot

	// DebugPermanentDisable indicates all non-manufacturer facilities are
	// permanently disabled such that no end user or developer cannot enable
	// them. Only the manufacturer indicated in the OEMID claim can enable them.
	// This also indicates that all debug facilities are currently disabled and
	// have been so since boot/start.
	DebugPermanentDisable

	// DebugFullPermanentDisable indicates that all debug capabilities for the
	// target device/sub-module are permanently disabled.
	DebugFullPermanentDisable
)

// Debug models the debug-disable type
type Debug uint

// Validate makes sure that the receiver is a valid Debug claim
func (d Debug) Validate() error {
	switch d {
	case DebugNotDisabled:
	case DebugDisabled:
	case DebugDisabledSinceBoot:
	case DebugPermanentDisable:
	case DebugFullPermanentDisable:
	default:
		return fmt.Errorf("out of range value %v for Debug type", d)
	}

	return nil
}
