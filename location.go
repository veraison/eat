// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

/*
=======
location-type = {
    latitude => number,
    longitude => number,
    ? altitude => number,
    ? accuracy => number,
    ? altitude-accuracy => number,
    ? heading => number,
    ? speed => number,
    ? timestamp => time-int,
    ? age => uint
}

latitude = 1
longitude = 2
altitude = 3
accuracy = 4
altitude-accuracy = 5
heading = 6
speed = 7
timestamp = 8
age = 9

age /= "age"
timestamp /= "timestamp"

latitude /= "lat"
longitude /= "long"
altitude /= "alt"
accuracy /= "accry"
altitude-accuracy /= "alt-accry"
heading /= "heading"
speed /= "speed"
*/

// Location models the location claim
type Location struct {
	Latitude         Number       `cbor:"1,keyasint" json:"lat"`
	Longitude        Number       `cbor:"2,keyasint" json:"long"`
	Altitude         *Number      `cbor:"3,keyasint,omitempty" json:"alt,omitempty"`
	Accuracy         *Number      `cbor:"4,keyasint,omitempty" json:"accry,omitempty"`
	AltitudeAccuracy *Number      `cbor:"5,keyasint,omitempty" json:"alt-accry,omitempty"`
	Heading          *Number      `cbor:"6,keyasint,omitempty" json:"heading,omitempty"`
	Speed            *Number      `cbor:"7,keyasint,omitempty" json:"speed,omitempty"`
	Timestamp        *NumericDate `cbor:"8,keyasint,omitempty" json:"timestamp,omitempty"`
	Age              *uint        `cbor:"9,keyasint,omitempty" json:"age,omitempty"`
}
