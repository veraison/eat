// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"fmt"
)

// $$Claims-Set-Claims //= (
//     measurement-results-label =>
//         [ + measurement-results-group ] )

// measurement-results-group = [
//     measurement-system: tstr,
//     measurement-results: [ + individual-result ]
// ]

// individual-result = [
//     result-id:  tstr / binary-data,
//     result:     result-type,
// ]

// result-type = comparison-success /
//               comparison-fail /
//               comparison-not-run /
//               measurement-absent

// comparison-success       = JC< "success",       1 >
// comparison-fail          = JC< "fail",          2 >
// comparison-not-run       = JC< "not-run",       3 >
// measurement-absent       = JC< "absent",        4 >

type MeasurementResultsGroup struct {
	_       struct{}           `cbor:",toarray"`
	System  string             `cbor:"0,keyasint"`
	Results []IndividualResult `cbor:"1,keyasint"`
}

func (m MeasurementResultsGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{m.System, m.Results})
}

func (m *MeasurementResultsGroup) UnmarshalJSON(data []byte) error {
	var tmp [2]json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("expected JSON array: %w", err)
	}

	if err := json.Unmarshal(tmp[0], &m.System); err != nil {
		return fmt.Errorf("invalid measurement-system: %w", err)
	}
	if err := json.Unmarshal(tmp[1], &m.Results); err != nil {
		return fmt.Errorf("invalid measurement-results: %w", err)
	}

	return nil
}

type IndividualResult struct {
	_    struct{}       `cbor:",toarray"`
	ID   StringOrBinary `cbor:"1,keyasint"`
	Type ResultType     `cbor:"2,keyasint"`
}

func (i IndividualResult) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]interface{}{i.ID, i.Type})
}

func (i *IndividualResult) UnmarshalJSON(data []byte) error {
	var tmp [2]json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("expected JSON array: %w", err)
	}

	if err := json.Unmarshal(tmp[0], &i.ID); err != nil {
		return fmt.Errorf("invalid result-id: %w", err)
	}
	if err := json.Unmarshal(tmp[1], &i.Type); err != nil {
		return fmt.Errorf("invalid result-type: %w", err)
	}

	return nil
}

const (
	ResultTypeInvalid = iota
	ResultTypeComparisonSuccess
	ResultTypeComparisonFail
	ResultTypeComparisonNotRun
	ResultTypeMeasurementAbsent
)

var resultTypeToString = map[ResultType]string{
	ResultTypeComparisonSuccess: "success",
	ResultTypeComparisonFail:    "fail",
	ResultTypeComparisonNotRun:  "not-run",
	ResultTypeMeasurementAbsent: "absent",
}

var stringToResultType = map[string]ResultType{
	"success": ResultTypeComparisonSuccess,
	"fail":    ResultTypeComparisonFail,
	"not-run": ResultTypeComparisonNotRun,
	"absent":  ResultTypeMeasurementAbsent,
}

type ResultType uint

func (r ResultType) MarshalJSON() ([]byte, error) {
	s, ok := resultTypeToString[r]
	if !ok {
		return nil, fmt.Errorf("out of range value %v for result-type", r)
	}
	return json.Marshal(s)
}

func (r *ResultType) UnmarshalJSON(data []byte) error {
	var t string
	if err := json.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("expected string: %w", err)
	}

	v, ok := stringToResultType[t]
	if !ok {
		return fmt.Errorf("invalid result-type string %q", v)
	}

	*r = v
	return nil
}
