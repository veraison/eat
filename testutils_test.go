// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	_ "embed"
	"encoding/hex"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/claims_set_RFC9711_A.1.1_simple_TEE_attestation.cbor
	rfcExampleClaimsSetSimpleTeeCBOR []byte

	//go:embed testdata/claims_set_RFC9711_A.1.2_board_and_device_submods.cbor
	rfcExampleClaimsSetBoardAndDeviceSubmodsCBOR []byte

	//go:embed testdata/claims_set_RFC9711_A.1.3_attestation_hw_block.cbor
	rfcExampleClaimsSetAttestationHwBlockCBOR []byte

	//go:embed testdata/claims_set_RFC9711_A.1.4_key_store_attestation.cbor
	rfcExampleClaimsSetKeyStoreAttestCBOR []byte

	//go:embed testdata/claims_set_RFC9711_A.1.5_sw_measurements_of_an_IoT_device.cbor
	rfcExampleClaimsSetSwMeasuresIotDeviceCBOR []byte

	//go:embed testdata/claims_set_RFC9711_A.1.6_attestation_results.json
	rfcExampleClaimsSetAttestResultJSON []byte

	//go:embed testdata/claims_set_RFC9711_A.1.7_token_with_submods.json
	rfcExampleClaimsSetTokenWithSubmodsJSON []byte

	//go:embed testdata/cwt_RFC9711_A.2.1_basic_CWT.cbor
	rfcExampleCwtBasicCBOR []byte // nolint:unused

	//go:embed testdata/claims_set_RFC9711_A.2.2_detached_digest_submod.cbor
	rfcExampleClaimsSetDetachedDigestSubmodCBOR []byte

	//go:embed testdata/bundle_RFC9711_A.2.2_detached_EAT_bundle.cbor
	rfcExampleDetachedEatBundleCBOR []byte

	//go:embed testdata/bundle_RFC9711_A.2.3_detached_EAT_bundle.json
	rfcExampleDetachedEatBundleJSON []byte
)

// Ptr returns the pointer to the specified value. This is useful when specifying
// literal values for fields that pointers to types for which directly taking a
// pointer of a literal is not possible (e.g. `Field: &"foo"` is not valid, do
// `Field: Ptr("foo")` instead).
func Ptr[T any](v T) *T {
	return &v
}

// MustHexDecode takes a string containing a sequence of hex byte values (no
// separator and no leading "0x") and returns a []byte containing the
// corresponding bytes. If the string has incorrect format, t.Fail() is
// invoked. If t is nil, this function panics instead.
func MustHexDecode(t *testing.T, s string) []byte {
	// allow long hex string to be split over multiple lines (with soft or hard
	// tab indentation)
	m := regexp.MustCompile("[ \t\n]")
	s = m.ReplaceAllString(s, "")

	data, err := hex.DecodeString(s)
	if t != nil {
		require.Nil(t, err)
	} else if err != nil {
		panic(err)
	}
	return data
}
