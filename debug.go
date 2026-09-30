// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

const (
	// DebugEnabledValue is asserted if any debug facility, even manufacturer
	// hardware diagnostics, is currently enabled
	DebugEnabledValue = iota

	// DebugDisabledValue indicates all debug facilities are currently disabled. It
	// may be possible to enable them in the future, and it may also be possible
	// that they were enabled in the past after the target device/sub-system
	// booted/started, but they are currently disabled.
	DebugDisabledValue

	// DebugDisabledSinceBootValue indicates all debug facilities are currently
	// disabled and have been so since the target device/sub-system
	// booted/started.
	DebugDisabledSinceBootValue

	// DebugPermanentDisableValue indicates all non-manufacturer facilities are
	// permanently disabled such that no end user or developer cannot enable
	// them. Only the manufacturer indicated in the OEMID claim can enable them.
	// This also indicates that all debug facilities are currently disabled and
	// have been so since boot/start.
	DebugPermanentDisableValue

	// DebugFullPermanentDisableValue indicates that all debug capabilities for the
	// target device/sub-module are permanently disabled.
	DebugFullPermanentDisableValue
)

var (

	// DebugEnabled is asserted if any debug facility, even manufacturer
	// hardware diagnostics, is currently enabled
	DebugEnabled = Debug{DebugEnabledValue}

	// DebugDisabled indicates all debug facilities are currently disabled. It
	// may be possible to enable them in the future, and it may also be possible
	// that they were enabled in the past after the target device/sub-system
	// booted/started, but they are currently disabled.
	DebugDisabled = Debug{DebugDisabledValue}

	// DebugDisabledSinceBoot indicates all debug facilities are currently
	// disabled and have been so since the target device/sub-system
	// booted/started.
	DebugDisabledSinceBoot = Debug{DebugDisabledSinceBootValue}

	// DebugPermanentDisable indicates all non-manufacturer facilities are
	// permanently disabled such that no end user or developer cannot enable
	// them. Only the manufacturer indicated in the OEMID claim can enable them.
	// This also indicates that all debug facilities are currently disabled and
	// have been so since boot/start.
	DebugPermanentDisable = Debug{DebugPermanentDisableValue}

	// DebugFullPermanentDisable indicates that all debug capabilities for the
	// target device/sub-module are permanently disabled.
	DebugFullPermanentDisable = Debug{DebugFullPermanentDisableValue}

	debugNameMap = map[uint]string{
		DebugEnabledValue:              "enabled",
		DebugDisabledValue:             "disabled",
		DebugDisabledSinceBootValue:    "disabled-since-boot",
		DebugPermanentDisableValue:     "disabled-permanently",
		DebugFullPermanentDisableValue: "disabled-fully-and-permanently",
	}

	debugValueMap = map[string]uint{
		"enabled":                        DebugEnabledValue,
		"disabled":                       DebugDisabledValue,
		"disabled-since-boot":            DebugDisabledSinceBootValue,
		"disabled-permanently":           DebugPermanentDisableValue,
		"disabled-fully-and-permanently": DebugFullPermanentDisableValue,
	}
)

// DebugFromString converts the provided string to a pointer to the
// corresponding Debug value. An error is returned if the provided string does
// not correspond to a Debug value.
func DebugFromString(text string) (*Debug, error) {
	value, ok := debugValueMap[text]
	if !ok {
		return nil, fmt.Errorf("invalid Debug value: %q", text)
	}

	return &Debug{value}, nil
}

// DebugFromString converts the provided uint to a pointer to the corresponding
// Debug value. An error is returned if the provided uint does not correspond
// to a Debug value.
func DebugFromUint(value uint) (*Debug, error) {
	_, ok := debugNameMap[value]
	if !ok {
		return nil, fmt.Errorf("invalid Debug value: %d", value)
	}

	return &Debug{value}, nil
}

// Debug models the debug-disable type
type Debug struct {
	value uint
}

// Uint returns the uint corresponding to this Debug's value.
func (o *Debug) Uint() uint {
	return o.value
}

// String returns the string corresponding to this Debug's value.
func (o *Debug) String() string {
	return debugNameMap[o.value]
}

func (o *Debug) MarshalCBOR() ([]byte, error) {
	return em.Marshal(o.value)
}

func (o *Debug) UnmarshalCBOR(data []byte) error {
	var value uint
	if err := dm.Unmarshal(data, &value); err != nil {
		return err
	}

	_, ok := debugNameMap[value]
	if !ok {
		return fmt.Errorf("invalid Debug value: %d", value)
	}

	o.value = value

	return nil
}

func (o *Debug) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *Debug) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	value, ok := debugValueMap[text]
	if !ok {
		return fmt.Errorf("invalid Debug value: %q", text)
	}

	o.value = value

	return nil
}
