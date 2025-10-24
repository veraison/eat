// Copyright 2025 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

type Measurement struct {
	_      struct{} `cbor:",toarray"` // TODO: implement Unmarshal.JSON
	Type   int      // coap-content-format, see https://www.iana.org/assignments/core-parameters/core-parameters.xhtml
	Format []byte   // bstr wrapped untagged-coswid, measured-component, ...
}
