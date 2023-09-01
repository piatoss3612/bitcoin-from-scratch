package main

import (
	"chapter04/ecc"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	e1 := big.NewInt(5000)
	pvk1, err := ecc.NewS256PrivateKey(e1.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(pvk1.Point().SEC()))

	e2 := big.NewInt(0).Exp(big.NewInt(2018), big.NewInt(5), nil)
	pvk2, err := ecc.NewS256PrivateKey(e2.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(pvk2.Point().SEC()))

	e3, _ := big.NewInt(0).SetString("deadbeef12345", 16)
	pvk3, err := ecc.NewS256PrivateKey(e3.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(pvk3.Point().SEC()))
}
