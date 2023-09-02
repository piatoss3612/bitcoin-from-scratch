package ecc

import (
	"bytes"
	"math/big"
	"strings"
	"unsafe"
)

// 바이트 슬라이스를 big.Int로 변환하는 함수
func BytesToBigInt(b []byte) *big.Int {
	return big.NewInt(0).SetBytes(b)
}

// 문자열을 바이트 슬라이스로 변환하는 함수
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// 두 원소값이 같은지 확인하는 함수
func sameBN(x, y *big.Int) bool {
	return x.Cmp(y) == 0
}

// 두 원소값을 더하는 함수
func addBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Add(x, y), mod)
}

// 두 원소값을 빼는 함수
func subBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Sub(x, y), mod)
}

// 두 원소값을 곱하는 함수
func mulBN(x, y, mod *big.Int) *big.Int {
	return big.NewInt(0).Mod(big.NewInt(0).Mul(x, y), mod)
}

// 원소값을 거듭제곱하는 함수
func powBN(x, exp, mod *big.Int) *big.Int {
	return big.NewInt(0).Exp(x, exp, mod)
}

// 원소값의 역원을 구하는 함수
func invBN(x, mod *big.Int) *big.Int {
	return big.NewInt(0).ModInverse(x, mod)
}

func sqrtBN(x, mod *big.Int) *big.Int {
	return big.NewInt(0).ModSqrt(x, mod)
}

// 무한원점인지 확인하는 함수
func isInfinity(x, y FieldElement) bool {
	return x == nil && y == nil
}

// 타원곡선 위에 있는지 확인하는 함수
func isOnCurve(x, y, a, b FieldElement) bool {
	prime := x.Prime()

	// y^2 == x^3 + ax + b
	left := powBN(y.Num(), big.NewInt(2), prime)
	right := addBN(
		addBN(
			powBN(x.Num(), big.NewInt(3), prime),
			mulBN(a.Num(), x.Num(), prime),
			prime,
		),
		b.Num(),
		prime,
	)

	return sameBN(left, right)
}

// 두 점이 서로 역원인지 확인하는 함수
func areInverse(x1, x2, y1, y2 FieldElement) bool {
	return x1.Equal(x2) && y1.NotEqual(y2)
}

// 두 타원곡선이 같은지 확인하는 함수
func sameCurve(a1, b1, a2, b2 FieldElement) bool {
	return a1.Equal(a2) && b1.Equal(b2)
}

// 두 점이 같은지 확인하는 함수
func samePoint(x1, y1, x2, y2 FieldElement) bool {
	return x1.Equal(x2) && y1.Equal(y2)
}

// num이 0보다 크거나 같고 prime보다 작은지 확인하는 함수
func inRange(num, prime *big.Int) bool {
	return num.Cmp(big.NewInt(0)) != -1 && num.Cmp(prime) == -1
}

var base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func EncodeBase58(s []byte) string {
	// 앞에 0x00이 몇개 있는지 확인
	count := 0
	for _, c := range s {
		if c == 0x00 {
			count++
		} else {
			break
		}
	}

	n := big.NewInt(0).SetBytes(s)              // 바이트 슬라이스를 big.Int로 변환
	prefix := bytes.Repeat([]byte{0x01}, count) // 0x00이 몇개 있는지에 따라 0x01을 count만큼 반복하여 prefix를 만듬

	result := strings.Builder{}

	// 58로 나눈 나머지를 base58Alphabet에서 찾아서 문자열을 만듦
	for n.Cmp(big.NewInt(0)) == 1 {
		mod := big.NewInt(0)
		n.DivMod(n, big.NewInt(58), mod)
		result.WriteByte(base58Alphabet[mod.Int64()])
	}

	// prefix를 붙임
	result.WriteString(BytesToString(prefix))

	// 문자열을 뒤집음
	resultStrBytes := StringToBytes(result.String())

	for i, j := 0, len(resultStrBytes)-1; i < j; i, j = i+1, j-1 {
		resultStrBytes[i], resultStrBytes[j] = resultStrBytes[j], resultStrBytes[i]
	}

	return BytesToString(resultStrBytes)
}
