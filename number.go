// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"math"
)

// Number models a CDDL number, i.e.:
//   uint / nint / float16 / float32 / float64
type Number float64

// MarshalCBOR encodes a Number using the smallest possible encoding.
// Note that choosing the smallest fitting float variant is decided
// by the encoding mode defined in cbor.go which is configured to use
// ShortestFloat == ShortestFloat16.  What remains to be done here is
// intercepting the uint / nint cases and dispatch them to the default
// int marshaler.
func (n Number) MarshalCBOR() ([]byte, error) {
	f := float64(n)

	if math.Trunc(f) == f {
		return em.Marshal(int(f))
	}

	return em.Marshal(f)
}
