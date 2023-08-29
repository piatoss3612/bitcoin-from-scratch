package main

import (
	"chapter03/ecc"
	"fmt"
	"math/big"
)

func main() {
	// 서명 생성
	bigZ, _ := new(big.Int).SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	bigR, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	bigS, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	bigPx, _ := new(big.Int).SetString("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
	bigPy, _ := new(big.Int).SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)

	bigN, _ := new(big.Int).SetString(ecc.N, 16)

	z, err := ecc.NewFieldElement(bigZ, bigN)
	if err != nil {
		panic(err)
	}

	r, err := ecc.NewFieldElement(bigR, bigN)
	if err != nil {
		panic(err)
	}

	s, err := ecc.NewFieldElement(bigS, bigN)
	if err != nil {
		panic(err)
	}

	sig := ecc.NewS256Signature(r, s)

	px, err := ecc.NewS256Field(bigPx)
	if err != nil {
		panic(err)
	}

	py, err := ecc.NewS256Field(bigPy)
	if err != nil {
		panic(err)
	}

	point, err := ecc.NewS256Point(px, py)
	if err != nil {
		panic(err)
	}

	// 서명 검증
	ok, err := point.Verify(z, sig)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ok: %t\n", ok)
}
