package eat

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestMeasuredComponent(t *testing.T) {
	assert := assert.New(t)
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthForbidden,
	}
	dm, err := decOpt.DecMode()
	assert.Nil(err)
	data := []byte{
		0xA2,             // map(2)
		0x01,             // unsigned(1)
		0x82,             // array(2)
		0x63,             // text(3)
		0x46, 0x6F, 0x6F, // "Foo"
		0x82,                         // array(2)
		0x65,                         // text(5)
		0x31, 0x2E, 0x33, 0x2E, 0x34, // "1.3.4"
		0x01,       // unsigned(1)
		0x02,       // unsigned(2)
		0x82,       // array(2)
		0x01,       // unsigned(1)
		0x58, 0x20, // bytes(32)
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	}

	var mc MeasuredComponent
	assert.Nil(dm.Unmarshal(data, &mc))
	assert.Equal(mc.Id.Name, "Foo")
	assert.Equal(mc.Id.Version.Version, "1.3.4")
}

func TestMyType(t *testing.T) {
	assert := assert.New(t)
	decOpt := cbor.DecOptions{
		IndefLength: cbor.IndefLengthForbidden,
	}
	dm, err := decOpt.DecMode()
	assert.Nil(err)

	data := []byte{
		0xA1,             // map(1)
		0x01,             // unsigned(1)
		0x81,             // array(1)
		0x63,             // text(3)
		0x46, 0x6F, 0x6F, // "Foo"
	}
	type MyType struct {
		Names []string `cbor:"1,keyasint"`
	}
	var test MyType
	assert.Nil(dm.Unmarshal(data, &test))
	assert.Equal(len(test.Names), 1)
	assert.Equal(test.Names[0], "Foo")
}
