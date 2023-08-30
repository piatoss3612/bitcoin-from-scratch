package ecc

import "math/big"

var (
	A = 0 // secp256k1 타원곡선의 a 계수
	B = 7 // secp256k1 타원곡선의 b 계수
	G Point
	N *big.Int
	P *big.Int // 2^256 - 2^32 - 977
)

func init() {
	bigGx, _ := new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	bigGy, _ := new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)
	n, _ := new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	gx, _ := NewS256FieldElement(bigGx)
	gy, _ := NewS256FieldElement(bigGy)
	g, _ := NewS256Point(gx, gy)

	G = g
	N = n
	P = big.NewInt(0).Sub(
		big.NewInt(0).Sub(
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), nil)),
		big.NewInt(977))
}
