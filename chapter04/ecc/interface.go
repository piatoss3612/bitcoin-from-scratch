package ecc

import (
	"fmt"
	"math/big"
)

// 유한체의 원소 인터페이스
type FieldElement interface {
	fmt.Stringer
	Num() *big.Int
	Prime() *big.Int
	Equal(other FieldElement) bool
	NotEqual(other FieldElement) bool
	Add(other FieldElement) (FieldElement, error)
	Sub(other FieldElement) (FieldElement, error)
	Mul(other FieldElement) (FieldElement, error)
	Pow(exp *big.Int) (FieldElement, error)
	Div(other FieldElement) (FieldElement, error)
}

// 서명 검증 인터페이스
type Verifier interface {
	Verify(z []byte, sig Signature) (bool, error)
}

// 서명 인터페이스
type Signer interface {
	Sign(z []byte) (Signature, error)
}

// 타원곡선의 점 인터페이스
type Point interface {
	fmt.Stringer
	Verifier
	X() FieldElement
	Y() FieldElement
	A() FieldElement
	B() FieldElement
	Equal(other Point) bool
	NotEqual(other Point) bool
	Add(other Point) (Point, error)
	Mul(coefficient *big.Int) (Point, error)
	SEC(compressed bool) []byte
	Hash160(compressed bool) []byte
	Address(compressed bool, testnet bool) string
}

// 서명 인터페이스
type Signature interface {
	fmt.Stringer
	R() *big.Int
	S() *big.Int
	DER() []byte
}

// 개인키 인터페이스
type PrivateKey interface {
	fmt.Stringer
	Signer
	Point() Point
	WIF(compressed bool, testnet bool) string
}
