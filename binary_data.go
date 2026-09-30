// Copyright 2026 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/base64"
	"encoding/json"
)

// BinaryData is a []byte that encodes as base64url_nopad in JSON (normally a
// []byte is encoded as base64std).
type BinaryData []byte

func (o BinaryData) MarshalJSON() ([]byte, error) {
	encoded := base64.RawURLEncoding.EncodeToString(o)
	return json.Marshal(encoded)
}

func (o *BinaryData) UnmarshalJSON(data []byte) error {
	var encoded string
	var err error

	if err = json.Unmarshal(data, &encoded); err != nil { // nolint:gocritic
		return err
	}

	*o, err = base64.RawURLEncoding.DecodeString(encoded)
	return err
}
