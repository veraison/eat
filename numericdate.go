// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"strconv"
	"time"
)

// NumericDate models RFC7519 NumericDate, i.e., the number of seconds since
// UNIX epoch
type NumericDate time.Time

// MarshalJSON unwraps the receiver NumericDate exposing its underlying Time and
// converts it to UNIX time.
func (nd NumericDate) MarshalJSON() ([]byte, error) {
	t := time.Time(nd)

	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

// MarshalCBOR unwraps the receiver NumericDate exposing its underlying Time
// and dispatches it to our custom encoding mode, which automatically adds the
// required tag (TimeTag == cbor.EncTagRequired)
func (nd NumericDate) MarshalCBOR() ([]byte, error) {
	t := time.Time(nd)

	return em.Marshal(t)
}
