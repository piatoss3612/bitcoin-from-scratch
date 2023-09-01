package main

import (
	"chapter03/ecc"
	"chapter03/hash256"
	"fmt"
	"math/big"
)

func main() {
	//sigTest1()
	//sigTest2()
	sigTest3()
}

func sigTest1() {
	// 서명 생성
	z, _ := new(big.Int).SetString("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", 16)
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	bigPx, _ := new(big.Int).SetString("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574", 16)
	bigPy, _ := new(big.Int).SetString("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4", 16)

	sig := ecc.NewS256Signature(r, s)

	px, err := ecc.NewS256FieldElement(bigPx)
	if err != nil {
		panic(err)
	}

	py, err := ecc.NewS256FieldElement(bigPy)
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

func sigTest2() {
	// 두번째 서명 생성
	e := new(big.Int).SetBytes(hash256.New([]byte("my secret")))
	z := new(big.Int).SetBytes(hash256.New([]byte("my message")))

	k := big.NewInt(1234567890)

	temp, err := ecc.G.Mul(k)
	if err != nil {
		panic(err)
	}

	r := temp.X().Num()

	kInv := big.NewInt(0).ModInverse(k, ecc.N)

	re := big.NewInt(0).Mod(big.NewInt(0).Mul(r, e), ecc.N)

	zre := big.NewInt(0).Mod(big.NewInt(0).Add(z, re), ecc.N)

	s := big.NewInt(0).Mod(big.NewInt(0).Mul(kInv, zre), ecc.N)

	point, err := ecc.G.Mul(e)
	if err != nil {
		panic(err)
	}

	fmt.Println(point)
	fmt.Println(z.Text(16))

	sig := ecc.NewS256Signature(r, s)
	fmt.Println(sig)
}

func sigTest3() {
	e := big.NewInt(12345)
	z := hash256.New([]byte("Programming Bitcoin!"))

	pvk, err := ecc.NewS256PrivateKey(e.Bytes())
	if err != nil {
		panic(err)
	}

	sig, err := pvk.Sign(z)
	if err != nil {
		panic(err)
	}

	fmt.Println(sig)
}
