package main

import (
	"chapter04/ecc"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	e1 := big.NewInt(5001)
	pvk1, err := ecc.NewS256PrivateKey(e1.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(pvk1.Point())

	sec := pvk1.Point().SEC(true)

	fmt.Println(hex.EncodeToString(sec))

	parsed, err := ecc.Parse(sec)
	if err != nil {
		panic(err)
	}

	fmt.Println(parsed)
}
