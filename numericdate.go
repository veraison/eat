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

// UnmarshalJSON populates the receiver NumericDate by interpreting the
// supplied data as UTF-8-encoded Unix timestamp.
func (nd *NumericDate) UnmarshalJSON(data []byte) error {
	s := string(data)

	i, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	*nd = NumericDate(time.Unix(int64(i), 0))

	return nil
}

// MarshalCBOR unwraps the receiver NumericDate exposing its underlying Time
// and dispatches it to our custom encoding mode, which automatically adds the
// required tag (TimeTag == cbor.EncTagRequired)
func (nd NumericDate) MarshalCBOR() ([]byte, error) {
	t := time.Time(nd)

	return em.Marshal(t)
}

// UnmarshalCBOR decodes the data into a Time and sets the receiver NumericDate
// to the decoded value.
func (nd *NumericDate) UnmarshalCBOR(data []byte) error {
	var t time.Time
	if err := dm.Unmarshal(data, &t); err != nil {
		return err
	}

	*nd = NumericDate(t)

	return nil
}
