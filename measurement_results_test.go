// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ResultSuccess = ResultType(ResultTypeComparisonSuccess)
	ResultFail    = ResultType(ResultTypeComparisonFail)
	ResultInvalid = ResultType(ResultTypeInvalid)
)

var (
	IndividualResultValidSuccess = IndividualResult{
		ID:   StringOrBinary{"component 1"},
		Type: ResultSuccess,
	}
	IndividualResultValidFail = IndividualResult{
		ID:   StringOrBinary{"component 2"},
		Type: ResultFail,
	}
	MeasurementResultsValid = MeasurementResultsGroup{
		System: "OS",
		Results: []IndividualResult{
			IndividualResultValidSuccess,
			IndividualResultValidFail,
		},
	}
)

func TestMeasurementResultsGroup_MarshalJSON_OK(t *testing.T) {
	t1, err := json.Marshal(MeasurementResultsValid)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["OS",[["component 1","success"],["component 2","fail"]]]`), t1)

	var m1 MeasurementResultsGroup
	err = json.Unmarshal([]byte(`["OS",[["component 1","success"],["component 2","fail"]]]`), &m1)
	assert.Nil(t, err)
	assert.Equal(t, MeasurementResultsValid, m1)
}

func TestMeasurementResultsGroup_MarshalJSON_NG(t *testing.T) {
	var m1 MeasurementResultsGroup
	err := json.Unmarshal([]byte(`["OS"]`), &m1)
	assert.NotNil(t, err)

	var m2 MeasurementResultsGroup
	err = json.Unmarshal([]byte(`"OS"`), &m2)
	assert.NotNil(t, err)

	var m3 MeasurementResultsGroup
	err = json.Unmarshal([]byte(`["OS","component 1"]`), &m3)
	assert.NotNil(t, err)

	var m4 MeasurementResultsGroup
	err = json.Unmarshal([]byte(`[0,"component 1"]`), &m4)
	assert.NotNil(t, err)
}

func TestMeasurementResultsGroup_MarshalCBOR_OK(t *testing.T) {
	// echo '["OS",[["component 1",1],["component 2",2]]]' | diag2cbor.rb | xxd -i
	encodedMeasurementResults := []byte{
		0x82, 0x62, 0x4f, 0x53, 0x82, 0x82, 0x6b, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
		0x6e, 0x65, 0x6e, 0x74, 0x20, 0x31, 0x01, 0x82, 0x6b, 0x63, 0x6f, 0x6d,
		0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x20, 0x32, 0x02,
	}

	t1, err := em.Marshal(MeasurementResultsValid)
	assert.Nil(t, err)
	assert.Equal(t, encodedMeasurementResults, t1)

	var m1 MeasurementResultsGroup
	err = dm.Unmarshal(encodedMeasurementResults, &m1)
	assert.Nil(t, err)
	assert.Equal(t, MeasurementResultsValid, m1)
}

func TestMeasurementResultsGroup_MarshalCBOR_NG(t *testing.T) {
	// 81         # array(1)
	//    62      # text(2)
	//       4f53 # "OS"
	var m1 MeasurementResultsGroup
	err := dm.Unmarshal([]byte{0x81, 62, 0x4f, 0x53}, &m1)
	assert.NotNil(t, err)

	var m2 MeasurementResultsGroup
	err = dm.Unmarshal([]byte{0x62, 0x4f, 0x53}, &m2)
	assert.NotNil(t, err)

	// 82                           # array(2)
	//    62                        # text(2)
	//       4f53                   # "OS"
	//    6b                        # text(11)
	//       636f6d706f6e656e742031 # "component 1"
	var m3 MeasurementResultsGroup
	err = dm.Unmarshal([]byte{
		0x82, 0x62, 0x4f, 0x53, 0x6b, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
		0x6e, 0x65, 0x6e, 0x74, 0x20, 0x31,
	}, &m3)
	assert.NotNil(t, err)

	// 82                           # array(2)
	//    00                        # unsigned(0)
	//    6b                        # text(11)
	//       636f6d706f6e656e742031 # "component 1"
	var m4 MeasurementResultsGroup
	err = dm.Unmarshal([]byte{
		0x82, 0x00, 0x6b, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
		0x6e, 0x74, 0x20, 0x31,
	}, &m4)
	assert.NotNil(t, err)
}

func TestIndividualResult_MarshalJSON_OK(t *testing.T) {
	t1, err := json.Marshal(IndividualResultValidSuccess)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["component 1","success"]`), t1)

	var i1 IndividualResult
	err = json.Unmarshal([]byte(`["component 1","success"]`), &i1)
	assert.Nil(t, err)
	assert.Equal(t, IndividualResultValidSuccess, i1)
}

func TestIndividualResult_MarshalJSON_NG(t *testing.T) {
	e1 := IndividualResult{
		Type: ResultSuccess,
	}
	_, err := json.Marshal(e1)
	assert.NotNil(t, err)

	var i1 IndividualResult
	err = json.Unmarshal([]byte(`["component 1"]`), &i1)
	assert.NotNil(t, err)

	var i2 IndividualResult
	err = json.Unmarshal([]byte(`"component 1"`), &i2)
	assert.NotNil(t, err)

	var i3 IndividualResult
	err = json.Unmarshal([]byte(`["component 1",1]`), &i3)
	assert.NotNil(t, err)

	var i4 IndividualResult
	err = json.Unmarshal([]byte(`[1,"fail"]`), &i4)
	assert.NotNil(t, err)
}

func TestResultType_MarshalJSON_OK(t *testing.T) {
	t1, err := json.Marshal(ResultSuccess)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"success"`), t1)

	var r1 ResultType
	err = json.Unmarshal([]byte(`"success"`), &r1)
	assert.Nil(t, err)
	assert.Equal(t, ResultSuccess, r1)

	t2, err := json.Marshal(ResultFail)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`"fail"`), t2)

	var r2 ResultType
	err = json.Unmarshal([]byte(`"fail"`), &r2)
	assert.Nil(t, err)
	assert.Equal(t, ResultFail, r2)
}

func TestResultType_MarshalJSON_NG(t *testing.T) {
	_, err := json.Marshal(ResultInvalid)
	assert.NotNil(t, err)

	var r1 ResultType
	err = json.Unmarshal([]byte(`"invalid"`), &r1)
	assert.NotNil(t, err)

	var r2 ResultType
	err = json.Unmarshal([]byte(`1`), &r2)
	assert.NotNil(t, err)
}
