package eat

import (
	"encoding/hex"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
)

// A nonce-claim may be single Nonce or an array of two or more.
//
//    nonce-claim = (
//        nonce => nonce-type / [ 2* nonce-type ]
//    )
type Nonces []Nonce

// MarshalCBOR encodes Nonces as a byte string, in case there is only
// one, or an array of byte strings if there are multiple.
func (ns Nonces) MarshalCBOR() ([]byte, error) {
	if len(ns) == 1 {
		return cbor.Marshal(ns[0])
	}

	return cbor.Marshal([]Nonce(ns))
}

// UnmarshalCBOR decodes nonce claim data. This may be a single byte string
// between 8 and 64 bytes long, or an array of two or more such strings.
func (ns *Nonces) UnmarshalCBOR(data []byte) error {
	if isCBORArray(data) {
		return cbor.Unmarshal(data, (*[]Nonce)(ns))
	}

	var n Nonce

	if err := cbor.Unmarshal(data, &n); err != nil {
		return err
	}

	*ns = append(*ns, n)

	return nil
}

// NonceFromHex creates a new Nonce instance from a string containing
// hex-encoded byte values and returns a pointer to it.
func NonceFromHex(text string) (*Nonce, error) {
	value, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return NewNonce(value)
}

// A Nonce is between 8 and 64 bytes
//    nonce-type = bstr .size (8..64)
type Nonce struct {
	value []byte
}

// NewNonce instantiates a new Nonce from the specified data and returns a
// pointer to it.
func NewNonce(data []byte) (*Nonce, error) {
	dlen := len(data)
	if dlen > 64 || dlen < 8 {
		return nil, fmt.Errorf("a nonce must be between 8 and 64 bytes long; found %v", dlen)
	}

	return &Nonce{data}, nil
}

// MarshalCBOR encodes the Nonce a CBOR byte string.
func (n *Nonce) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(n.value)
}

// UnmarshalCBOR decodes a CBOR byte string and uses it as the Nonce value.
func (n *Nonce) UnmarshalCBOR(data []byte) error {
	var value []byte

	if err := cbor.Unmarshal(data, &value); err != nil {
		return err
	}

	vlen := len(value)
	if vlen > 64 || vlen < 8 {
		return fmt.Errorf("a nonce must be between 8 and 64 bytes long; found %v", vlen)
	}

	n.value = value

	return nil
}

func isCBORArray(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	return (data[0] & 0xe0) == 0x80
}
