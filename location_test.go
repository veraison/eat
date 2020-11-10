// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLocation_Marshal(t *testing.T) {
	one := Number(1)
	ts := NumericDate(time.Date(2020, time.November, 10, 0, 0, 0, 0, time.UTC))
	age := uint(969)

	tests := []struct {
		name         string
		testVector   Location
		expectedCBOR []byte
		expectedJSON string
	}{
		{
			"mandatory fields only (float64)",
			Location{
				Latitude:  12.34,
				Longitude: 56.78,
			},
			/*
			   a2                     # map(2)
			      01                  # unsigned(1)
			      fb 4028ae147ae147ae # primitive(4623136420479977390)
			      02                  # unsigned(2)
			      fb 404c63d70a3d70a4 # primitive(4633187891898314916)
			*/
			[]byte{
				0xa2, 0x01, 0xfb, 0x40, 0x28, 0xae, 0x14, 0x7a, 0xe1, 0x47,
				0xae, 0x02, 0xfb, 0x40, 0x4c, 0x63, 0xd7, 0x0a, 0x3d, 0x70,
				0xa4,
			},
			`{"lat":12.34,"long":56.78}`,
		},
		{
			"mandatory fields only (float16)",
			Location{
				Latitude:  1.9306640625,
				Longitude: 2.00390625,
			},
			/*
			   a2         # map(2)
			      01      # unsigned(1)
			      f9 3fb9 # primitive(16313)
			      02      # unsigned(2)
			      f9 4002 # primitive(16386)
			*/
			[]byte{
				0xa2, 0x01, 0xf9, 0x3f, 0xb9, 0x02, 0xf9, 0x40, 0x02,
			},
			`{"lat":1.9306640625,"long":2.00390625}`,
		},
		{
			"mandatory fields only (float32)",
			Location{
				Latitude:  3.1414999961853027,
				Longitude: -12.119999885559082,
			},
			/*
			   a2             # map(2)
			      01          # unsigned(1)
			      fa 40490e56 # primitive(1078529622)
			      02          # unsigned(2)
			      fa c141eb85 # primitive(3242322821)
			*/
			[]byte{
				0xa2, 0x01, 0xfa, 0x40, 0x49, 0x0e, 0x56, 0x02, 0xfa, 0xc1,
				0x41, 0xeb, 0x85,
			},
			`{"lat":3.1414999961853027,"long":-12.119999885559082}`,
		},
		{
			"mandatory fields only (uint / nint)",
			Location{
				Latitude:  3,
				Longitude: -12,
			},
			/*
			   a2    # map(2)
			      01 # unsigned(1)
			      03 # unsigned(3)
			      02 # unsigned(2)
			      2b # negative(11)
			*/
			[]byte{
				0xa2, 0x01, 0x03, 0x02, 0x2b,
			},
			`{"lat":3,"long":-12}`,
		},
		{
			"all fields",
			Location{
				Latitude:         3,
				Longitude:        -12.1,
				Altitude:         &one,
				Accuracy:         &one,
				AltitudeAccuracy: &one,
				Heading:          &one,
				Speed:            &one,
				Timestamp:        &ts,
				Age:              &age,
			},
			/*
			   a9                     # map(9)
			      01                  # unsigned(1)
			      03                  # unsigned(3)
			      02                  # unsigned(2)
			      fb c028333333333333 # primitive(13846373349345932083)
			      03                  # unsigned(3)
			      01                  # unsigned(1)
			      04                  # unsigned(4)
			      01                  # unsigned(1)
			      05                  # unsigned(5)
			      01                  # unsigned(1)
			      06                  # unsigned(6)
			      01                  # unsigned(1)
			      07                  # unsigned(7)
			      01                  # unsigned(1)
			      08                  # unsigned(8)
			      c1                  # tag(1)
			         1a 5fa9d800      # unsigned(1604966400)
			      09                  # unsigned(9)
			      19 03c9             # unsigned(969)
			*/
			[]byte{
				0xa9, 0x01, 0x03, 0x02, 0xfb, 0xc0, 0x28, 0x33, 0x33, 0x33,
				0x33, 0x33, 0x33, 0x03, 0x01, 0x04, 0x01, 0x05, 0x01, 0x06,
				0x01, 0x07, 0x01, 0x08, 0xc1, 0x1a, 0x5f, 0xa9, 0xd8, 0x00,
				0x09, 0x19, 0x03, 0xc9,
			},
			`{"lat":3,"long":-12.1,"alt":1,"accry":1,"alt-accry":1,"heading":1,"speed":1,"timestamp":1604966400,"age":969}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := em.Marshal(test.testVector)
			assert.Nil(t, err)
			assert.Equal(t, test.expectedCBOR, actual)

			actual, err = json.Marshal(test.testVector)
			assert.Nil(t, err)
			assert.JSONEq(t, test.expectedJSON, string(actual))
		})
	}
}

func TestLocation_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		tvCBOR   []byte
		tvJSON   string
		expected Location
	}{
		{
			"mandatory fields only (uint)",
			/*
			   a2       # map(2)
			      00    # unsigned(0)
			      0c    # unsigned(12)
			      01    # unsigned(1)
			      18 38 # unsigned(56)
			*/
			[]byte{
				0xa2, 0x01, 0x0c, 0x02, 0x18, 0x38,
			},
			`{"lat":12,"long":56}`,
			Location{
				Latitude:  12,
				Longitude: 56,
			},
		},
		{
			"mandatory fields only (nint)",
			/*
			   a2       # map(2)
			      01    # unsigned(1)
			      2b    # negative(11)
			      02    # unsigned(2)
			      38 37 # negative(55)
			*/
			[]byte{
				0xa2, 0x01, 0x2b, 0x02, 0x38, 0x37,
			},
			`{"lat":-12,"long":-56}`,
			Location{
				Latitude:  -12,
				Longitude: -56,
			},
		},
		{
			"mandatory fields only (float16)",
			/*
			   a2         # map(2)
			      01      # unsigned(1)
			      f9 3c00 # primitive(15360)
			      02      # unsigned(2)
			      f9 4000 # primitive(16384)
			*/
			[]byte{
				0xa2, 0x01, 0xf9, 0x3c, 0x00, 0x02, 0xfb, 0x3f, 0xb9, 0x99,
				0x99, 0x99, 0x99, 0x99, 0x9a,
			},
			`{"lat":1.0,"long":0.1}`,
			Location{
				Latitude:  1.0,
				Longitude: 0.1,
			},
		},
		{
			"mandatory fields only (float32)",
			/*
			   a2             # map(2)
			      01          # unsigned(1)
			      fa 40490e56 # primitive(1078529622)
			      02          # unsigned(2)
			      fa c141eb85 # primitive(3242322821)
			*/
			[]byte{
				0xa2, 0x01, 0xfa, 0x40, 0x49, 0x0e, 0x56, 0x02, 0xfa, 0xc1,
				0x41, 0xeb, 0x85,
			},
			`{"lat":3.1414999961853027,"long":-12.119999885559082}`,
			Location{
				Latitude:  3.1414999961853027,
				Longitude: -12.119999885559082,
			},
		},
		{
			"mandatory fields only (float64)",
			/*
			   a2                     # map(2)
			      01                  # unsigned(1)
			      fb 3ff3333333333333 # primitive(4608083138725491507)
			      02                  # unsigned(2)
			      fb 400b333333333333 # primitive(4614838538166547251)
			*/
			[]byte{
				0xa2, 0x01, 0xfb, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x33, 0x02, 0xfb, 0x40, 0x0b, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x33,
			},
			`{"lat":1.2,"long":3.4}`,
			Location{
				Latitude:  1.2,
				Longitude: 3.4,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Location{}
			err := dm.Unmarshal(test.tvCBOR, &actual)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)

			actual = Location{}
			err = json.Unmarshal([]byte(test.tvJSON), &actual)
			assert.Nil(t, err)
			assert.Equal(t, test.expected, actual)
		})
	}
}
