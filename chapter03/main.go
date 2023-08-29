package main

import (
	"chapter03/ecc"
	"fmt"
	"math/big"
)

func main() {
	gx := "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	gy := "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	n := "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"

	bigX, _ := new(big.Int).SetString(gx, 16)
	bigY, _ := new(big.Int).SetString(gy, 16)
	bigN, _ := new(big.Int).SetString(n, 16)

	x1, err := ecc.NewS256Field(bigX)
	if err != nil {
		panic(err)
	}

	y1, err := ecc.NewS256Field(bigY)
	if err != nil {
		panic(err)
	}

	G, err := ecc.NewS256Point(x1, y1)
	if err != nil {
		panic(err)
	}

	nG, err := G.Mul(bigN)
	if err != nil {
		panic(err)
	}

	fmt.Println(nG)
}
