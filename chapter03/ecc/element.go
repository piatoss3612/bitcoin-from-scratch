package ecc

import (
	"fmt"
	"math/big"
)

// 유한체의 원소를 나타내는 구조체
type fieldElement struct {
	num   *big.Int
	prime *big.Int
}

// 유한체의 원소를 생성하는 함수
func NewFieldElement(num, prime *big.Int) (FieldElement, error) {
	// num과 prime이 nil인 경우 유한체의 원소를 생성할 수 없음
	if num == nil || prime == nil {
		return nil, fmt.Errorf("num or prime cannot be nil")
	}

	// 유한체의 원소가 0보다 작거나 위수보다 크거나 같은 경우 유한체의 원소를 생성할 수 없음
	if !inRange(num, prime) {
		return nil, fmt.Errorf("num %s not in field range 0 to %s", num, prime.Sub(prime, big.NewInt(1)))
	}

	return &fieldElement{num, prime}, nil
}

// 유한체의 원소값을 반환하는 함수
func (f fieldElement) Num() *big.Int {
	return f.num
}

// 유한체의 위수를 반환하는 함수
func (f fieldElement) Prime() *big.Int {
	return f.prime
}

// 유한체의 원소를 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (f fieldElement) String() string {
	return fmt.Sprintf("FieldElement_%s(%s)", f.prime, f.num)
}

// 유한체의 원소가 같은 유한체에 속하고 같은 값을 가지는지 확인하는 함수
func (f fieldElement) Equal(other FieldElement) bool {
	return f.equal(other)
}

// 유한체의 원소가 다른 유한체에 속하거나 다른 값을 가지는지 확인하는 함수
func (f fieldElement) NotEqual(other FieldElement) bool {
	return !f.equal(other)
}

// 유한체의 원소가 같은 유한체에 속하고 같은 값을 가지는지 확인하는 내부 함수
func (f fieldElement) equal(other FieldElement) bool {
	return sameBN(f.num, other.Num()) && sameBN(f.prime, other.Prime())
}

// 유한체의 원소를 더하는 함수 (더한 결과는 같은 유한체에 속함)
func (f fieldElement) Add(other FieldElement) (FieldElement, error) {
	num, err := f.add(other)
	if err != nil {
		return nil, err
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 더하는 내부 함수 (더한 결과는 같은 유한체에 속함)
func (f fieldElement) add(other FieldElement) (*big.Int, error) {
	if !sameBN(f.prime, other.Prime()) {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.Prime())
	}

	return addBN(f.num, other.Num(), f.prime), nil
}

// 유한체의 원소를 빼는 함수 (뺀 결과는 같은 유한체에 속함)
func (f fieldElement) Sub(other FieldElement) (FieldElement, error) {
	num, err := f.sub(other)
	if err != nil {
		return nil, err
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 빼는 내부 함수 (뺀 결과는 같은 유한체에 속함)
func (f fieldElement) sub(other FieldElement) (*big.Int, error) {
	if !sameBN(f.prime, other.Prime()) {
		return nil, fmt.Errorf("cannot subtract two numbers in different Fields %d, %d", f.prime, other.Prime())
	}

	return subBN(f.num, other.Num(), f.prime), nil
}

// 유한체의 원소를 곱하는 함수 (곱한 결과는 같은 유한체에 속함)
func (f fieldElement) Mul(other FieldElement) (FieldElement, error) {
	num, err := f.mul(other)
	if err != nil {
		return nil, err
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 곱하는 내부 함수 (곱한 결과는 같은 유한체에 속함)
func (f fieldElement) mul(other FieldElement) (*big.Int, error) {
	if !sameBN(f.prime, other.Prime()) {
		return nil, fmt.Errorf("cannot multiply two numbers in different Fields %d, %d", f.prime, other.Prime())
	}

	return mulBN(f.num, other.Num(), f.prime), nil
}

// 유한체의 원소를 거듭제곱하는 함수 (거듭제곱한 결과는 같은 유한체에 속함)
func (f fieldElement) Pow(exp *big.Int) (FieldElement, error) {
	return NewFieldElement(f.pow(exp), f.prime)
}

func (f fieldElement) pow(exp *big.Int) *big.Int {
	// 페르마의 소정리를 이용하여 지수를 줄이고 거듭제곱을 계산
	return powBN(f.num, big.NewInt(0).Mod(exp, big.NewInt(0).Sub(f.prime, big.NewInt(1))), f.prime)
}

// 유한체의 원소를 나누는 함수 (나눈 결과는 같은 유한체에 속함)
func (f fieldElement) Div(other FieldElement) (FieldElement, error) {
	num, err := f.div(other)
	if err != nil {
		return nil, err
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 나누는 내부 함수 (나눈 결과는 같은 유한체에 속함)
func (f fieldElement) div(other FieldElement) (*big.Int, error) {
	if !sameBN(f.prime, other.Prime()) {
		return nil, fmt.Errorf("cannot divide two numbers in different Fields %d, %d", f.prime, other.Prime())
	}

	// 페르마의 소정리를 이용하여 나눗셈을 곱셈으로 변환
	// a / b = a * (b ** (p - 2)) mod p
	return mulBN(f.num, invBN(other.Num(), f.prime), f.prime), nil
}

// secp256k1 유한체의 원소 구조체
type s256FieldElement struct {
	fieldElement
}

// secp256k1 유한체의 원소를 생성하는 함수
func NewS256FieldElement(num *big.Int) (FieldElement, error) {
	if num == nil {
		return nil, fmt.Errorf("num cannot be nil")
	}

	// 유한체의 원소가 유한체의 크기보다 크거나 같은 경우 에러
	if num.Cmp(P) != -1 {
		return nil, fmt.Errorf("num %s not in field range 0 to %s", num, P.Sub(P, big.NewInt(1)))
	}

	f := fieldElement{num, P}

	return &s256FieldElement{f}, nil
}

// secp256k1 유한체의 원소를 더하는 함수 (더한 결과는 같은 유한체에 속함)
func (f s256FieldElement) Add(other FieldElement) (FieldElement, error) {
	num, err := f.add(other)
	if err != nil {
		return nil, err
	}
	return NewS256FieldElement(num)
}

// secp256k1 유한체의 원소를 빼는 함수 (뺀 결과는 같은 유한체에 속함)
func (f s256FieldElement) Sub(other FieldElement) (FieldElement, error) {
	num, err := f.sub(other)
	if err != nil {
		return nil, err
	}
	return NewS256FieldElement(num)
}

// secp256k1 유한체의 원소를 곱하는 함수 (곱한 결과는 같은 유한체에 속함)
func (f s256FieldElement) Mul(other FieldElement) (FieldElement, error) {
	num, err := f.mul(other)
	if err != nil {
		return nil, err
	}
	return NewS256FieldElement(num)
}

// secp256k1 유한체의 원소를 거듭제곱하는 함수 (거듭제곱한 결과는 같은 유한체에 속함)
func (f s256FieldElement) Pow(exp *big.Int) (FieldElement, error) {
	return NewS256FieldElement(f.pow(exp))
}

// secp256k1 유한체의 원소를 나누는 함수 (나눈 결과는 같은 유한체에 속함)
func (f s256FieldElement) Div(other FieldElement) (FieldElement, error) {
	num, err := f.div(other)
	if err != nil {
		return nil, err
	}
	return NewS256FieldElement(num)
}
