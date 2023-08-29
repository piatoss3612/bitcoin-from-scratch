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
	p := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	p.Sub(p, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), nil))
	p.Sub(p, big.NewInt(977))

	bigX, _ := new(big.Int).SetString(gx, 16)
	bigY, _ := new(big.Int).SetString(gy, 16)
	bigN, _ := new(big.Int).SetString(n, 16)

	// y^2 mod p = x^3 + 7 mod p
	left := big.NewInt(0).Exp(bigY, big.NewInt(2), nil)
	left.Mod(left, p)
	right := big.NewInt(0).Exp(bigX, big.NewInt(3), nil)
	right.Add(right, big.NewInt(7))
	right.Mod(right, p)

	if left.Cmp(right) == 0 {
		println("y^2 mod p = x^3 + 7 mod p")
	} else {
		println("y^2 mod p != x^3 + 7 mod p")
	}

	x1, err := ecc.NewFieldElement(bigX, p)
	if err != nil {
		panic(err)
	}

	y1, err := ecc.NewFieldElement(bigY, p)
	if err != nil {
		panic(err)
	}

	zero, err := ecc.NewFieldElement(big.NewInt(0), p)
	if err != nil {
		panic(err)
	}

	seven, err := ecc.NewFieldElement(big.NewInt(7), p)
	if err != nil {
		panic(err)
	}

	G, err := ecc.New(x1, y1, zero, seven)
	if err != nil {
		panic(err)
	}

	nG, err := G.Mul(bigN)
	if err != nil {
		panic(err)
	}

	fmt.Println(nG)
}
