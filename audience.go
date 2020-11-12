// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"errors"
)

// In the general case, the "aud" value is an array of case- sensitive strings,
// each containing a StringOrURI value.  In the special case when the JWT has
// one audience, the "aud" value MAY be a single case-sensitive string
// containing a StringOrURI value.
type Audience []StringOrURI

// MarshalCBOR encodes Audience as a StringOrURI, in case there is only
// one, or an array of StringOrURI's if there are multiple.
func (a Audience) MarshalCBOR() ([]byte, error) {
	if len(a) == 1 {
		return em.Marshal(a[0])
	}

	return em.Marshal([]StringOrURI(a))
}

// UnmarshalCBOR decodes audience claim data. This may be a single StringOrURI,
// or an array of such.
func (a *Audience) UnmarshalCBOR(data []byte) error {
	if isCBORArray(data) {
		return dm.Unmarshal(data, (*[]StringOrURI)(a))
	}

	var v StringOrURI

	if err := dm.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = Audience{v}

	return nil
}

// MarshalJSON encodes the receiver Audience as a JSON string
func (a Audience) MarshalJSON() ([]byte, error) {
	if len(a) == 1 {
		return json.Marshal(a[0])
	}

	return nil, errors.New("TODO handle array of audiences")
}

// UnmarshalJSON decodes a JSON string into  the receiver Audience
func (a *Audience) UnmarshalJSON(data []byte) error {
	var v interface{}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch t := v.(type) {
	case string:
		var s StringOrURI
		if err := s.FromString(t); err != nil {
			return err
		}
		*a = Audience{s}
		return nil
	default:
		return errors.New("TODO handle array of nonces")
	}
}
