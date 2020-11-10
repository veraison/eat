package eat

import (
	"encoding/hex"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
)

type Nonces []Nonce

func (ns Nonces) MarshalCBOR() ([]byte, error) {
	if len(ns) == 1 {
		return cbor.Marshal(ns[0])
	}

	return cbor.Marshal([]Nonce(ns))
}

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

func NonceFromHex(text string) (*Nonce, error) {
	value, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return NewNonce(value)
}

type Nonce struct {
	value []byte
}

func NewNonce(data []byte) (*Nonce, error) {
	dlen := len(data)
	if dlen > 64 || dlen < 8 {
		return nil, fmt.Errorf("a nonce must be between 8 and 64 bytes long; found %v", dlen)
	}

	return &Nonce{data}, nil
}

func (n *Nonce) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(n.value)
}

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
