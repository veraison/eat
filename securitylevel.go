// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "fmt"

/*
security-level-type = &(

	unrestricted: 1,
	restricted: 2,
	secure-restricted: 3,
	hardware: 4

)
*/
const (
	// There is some expectation that implementer will protect the
	// attestation signing keys at this level.  Otherwise, the EAT
	// provides no meaningful security assurance.
	SecLevelUnrestricted = iota

	// Enties at this level should not be general-purpose operating
	// environments that host features such as app download systems, web
	// browsers, and complex productivity applications. It is akin to the
	// Secure Restricted level (see below) without the security orientation.
	// E.g. a Wi-Fi subsystem, an IoT camera, or a sensor device.
	SecLevelRestricted

	// Entities at this level must meet the criteria defined by FIDO
	// Allowed Restricted Operating Environments [1]. Examples include TEE's
	// and schemes using virtualization-based security. Like the FIDO
	// security goal, security at this level is aimed defending well
	// against large-scale network / remote attacks against the device.
	//
	// [1] https://fidoalliance.org/specs/fido-security-requirements-v1.0-fd-20170524/
	// fido-authenticator-allowed-restricted-operating-environments-list_20170524.pdf
	SecLevelSecureRestricted

	// Entities at this level must include substantial defense against
	// physical or electrical attacks against the device itself. It is
	// assumed any potential attacker has captured the device and can
	// disassemble it. Examples include TPMs and Secure Elements.
	SecLevelHardware
)

// Security Level claim type
type SecurityLevel uint

// Validate SecLevel to make sure is with thin bounds allowed by the spec.
func (s SecurityLevel) Validate() error {
	switch s {
	case SecLevelUnrestricted:
	case SecLevelRestricted:
	case SecLevelSecureRestricted:
	case SecLevelHardware:
	default:
		return fmt.Errorf("out of range value %v for SecurityLevel type", s)
	}

	return nil
}
