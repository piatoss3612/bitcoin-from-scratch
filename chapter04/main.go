package main

import (
	"chapter04/ecc"
	"fmt"
	"math/big"
)

func main() {
	e1 := big.NewInt(5003)

	pvk1, err := ecc.NewS256PrivateKey(e1.Bytes())
	if err != nil {
		panic(err)
	}

	wif1 := pvk1.WIF(true, true)

	fmt.Printf("wif1: %s\n", wif1)

	e2 := big.NewInt(0).Exp(big.NewInt(2021), big.NewInt(5), nil)

	pvk2, err := ecc.NewS256PrivateKey(e2.Bytes())
	if err != nil {
		panic(err)
	}

	wif2 := pvk2.WIF(false, true)

	fmt.Printf("wif2: %s\n", wif2)

	e3, _ := new(big.Int).SetString("54321deadbeef", 16)

	pvk3, err := ecc.NewS256PrivateKey(e3.Bytes())
	if err != nil {
		panic(err)
	}

	wif3 := pvk3.WIF(true, false)

	fmt.Printf("wif3: %s\n", wif3)
}
