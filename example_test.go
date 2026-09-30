// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "fmt"

func ExampleEat_MarshalJSON() {
	nonce := Nonce{}

	if err := nonce.AddHex("0000000000000000"); err != nil {
		panic(err)
	}

	// if required by the use case, add more nonces

	t := Eat{
		ClaimsSet{
			Nonce: &nonce,
		},
	}

	j, err := t.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))

	// Output: {"eat_nonce":"AAAAAAAAAAA"}
}

func ExampleEat_UnmarshalJSON() {
	t := Eat{}

	data := []byte(`{"eat_nonce":"AAAAAAAAAAA"}`)

	if err := t.UnmarshalJSON(data); err != nil {
		panic(err)
	}

	fmt.Printf("nonces found: %d\n", t.Claims().Nonce.Len())
	fmt.Printf("nonce: %x\n", t.Claims().Nonce.Get(0))

	// Output:
	// nonces found: 1
	// nonce: 0000000000000000
}
