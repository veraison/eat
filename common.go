// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

func isCBORTag(data []byte) bool {
	return (data[0] & 0xe0) == 0xc0
}

func isCBORTextString(data []byte) bool {
	return (data[0] & 0xe0) == 0x60
}

func isCBORByteString(data []byte) bool {
	return (data[0] & 0xe0) == 0x40
}

func isCBORArray(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	return (data[0] & 0xe0) == 0x80
}
