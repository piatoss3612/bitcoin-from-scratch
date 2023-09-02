package main

import (
	"chapter04/ecc"
	"fmt"
	"math/big"
)

func main() {
	e1 := big.NewInt(5002)

	pvk1, err := ecc.NewS256PrivateKey(e1.Bytes())
	if err != nil {
		panic(err)
	}

	addr1 := pvk1.Point().Address(false, true)

	fmt.Printf("addr1: %s\n", addr1)

	e2 := big.NewInt(0).Exp(big.NewInt(2020), big.NewInt(5), nil)

	pvk2, err := ecc.NewS256PrivateKey(e2.Bytes())
	if err != nil {
		panic(err)
	}

	addr2 := pvk2.Point().Address(true, true)

	fmt.Printf("addr2: %s\n", addr2)

	e3, _ := new(big.Int).SetString("12345deadbeef", 16)

	pvk3, err := ecc.NewS256PrivateKey(e3.Bytes())
	if err != nil {
		panic(err)
	}

	addr3 := pvk3.Point().Address(true, false)

	fmt.Printf("addr3: %s\n", addr3)
}
