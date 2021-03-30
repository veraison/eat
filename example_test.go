// Copyright 2020 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0

package eat

import "fmt"

func ExampleEat_ToJSON() {
	nonce := Nonces{}

	if err := nonce.AddHex("0000000000000000"); err != nil {
		panic(err)
	}

	// if required by the use case, add more nonces

	t := Eat{
		Nonce: &nonce,
	}

	j, err := t.ToJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(j))

	// Output: {"nonce":"AAAAAAAAAAA="}
}

func ExampleEat_FromJSON() {
	t := Eat{}

	data := []byte(`{"nonce":"AAAAAAAAAAA="}`)

	if err := t.FromJSON(data); err != nil {
		panic(err)
	}

	if err := t.Nonce.Validate(); err != nil {
		panic(err)
	}

	fmt.Printf("nonces found: %d\n", t.Nonce.Len())
	fmt.Printf("nonce: %x\n", t.Nonce.GetI(0))

	// Output:
	// nonces found: 1
	// nonce: 0000000000000000
}
