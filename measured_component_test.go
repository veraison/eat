package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedHashValue = B64Url{
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	}
)

func TestDigest_MarshalJSON_OK(t *testing.T) {
	data := []byte(`[1,"3q2-796tvu_erb7v3q2-796tvu_erb7v3q2-796tvu8"]`)

	var digest Digest
	assert.Nil(t, json.Unmarshal(data, &digest))
	assert.Equal(t, 1, digest.Alg)
	assert.Equal(t, expectedHashValue, digest.Val)

	encoded, err := json.Marshal(digest)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestDigest_MarshalCBOR_OK(t *testing.T) {
	data := []byte{
		0x82,       // array(2)
		0x01,       // unsigned(1)
		0x58, 0x20, // bytes(32)
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	}

	var digest Digest
	assert.Nil(t, dm.Unmarshal(data, &digest))
	assert.Equal(t, 1, digest.Alg)
	assert.Equal(t, expectedHashValue, digest.Val)

	encoded, err := em.Marshal(digest)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestMeasuredComponent_MarshalJSON_OK(t *testing.T) {
	data := []byte(`{"id":["Foo",["1.3.4","multipartnumeric"]],"measurement":[1,"3q2-796tvu_erb7v3q2-796tvu_erb7v3q2-796tvu8"]}`)

	var mc MeasuredComponent
	assert.Nil(t, json.Unmarshal(data, &mc))
	assert.Equal(t, "Foo", mc.Id.Name)
	assert.Equal(t, "1.3.4", mc.Id.Version.Version)
	assert.Equal(t, 1, mc.Measurement.Alg)
	assert.Equal(t, expectedHashValue, mc.Measurement.Val)

	encoded, err := json.Marshal(mc)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}

func TestMeasuredComponent_MarshalCBOR_OK(t *testing.T) {
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
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
		0xDE, 0xAD, 0xBE, 0xEF, 0xDE, 0xAD, 0xBE, 0xEF,
	}

	var mc MeasuredComponent
	assert.Nil(t, dm.Unmarshal(data, &mc))
	assert.Equal(t, "Foo", mc.Id.Name)
	assert.Equal(t, "1.3.4", mc.Id.Version.Version)
	assert.Equal(t, 1, mc.Measurement.Alg)
	assert.Equal(t, expectedHashValue, mc.Measurement.Val)

	encoded, err := em.Marshal(mc)
	assert.Nil(t, err)
	assert.Equal(t, data, encoded)
}
