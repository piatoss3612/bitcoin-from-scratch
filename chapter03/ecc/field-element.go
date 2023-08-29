package ecc

import (
	"fmt"
	"math/big"
)

// 유한체의 원소를 나타내는 구조체
type FieldElement struct {
	num   *big.Int
	prime *big.Int
}

// 유한체의 원소를 생성하는 함수
func NewFieldElement(num, prime *big.Int) (*FieldElement, error) {
	if num.Cmp(prime) != -1 || num.Cmp(big.NewInt(0)) == -1 {
		return nil, fmt.Errorf("num %s not in field range 0 to %s", num, prime.Sub(prime, big.NewInt(1)))
	}

	return &FieldElement{num, prime}, nil
}

// 유한체의 원소를 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (f FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.num)
}

// 유한체의 원소가 같은 유한체에 속하고 같은 값을 가지는지 확인하는 함수
func (f FieldElement) Equal(other FieldElement) bool {
	return f.num.Cmp(other.num) == 0 && f.prime.Cmp(other.prime) == 0
}

func (f FieldElement) NotEqual(other FieldElement) bool {
	return f.num.Cmp(other.num) != 0 || f.prime.Cmp(other.prime) != 0
}

// 유한체의 원소를 더하는 함수 (더한 결과는 같은 유한체에 속함)
func (f FieldElement) Add(other FieldElement) (*FieldElement, error) {
	if f.prime.Cmp(other.prime) != 0 {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := big.NewInt(0).Mod(big.NewInt(0).Add(f.num, other.num), f.prime)
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 빼는 함수 (뺀 결과는 같은 유한체에 속함)
func (f FieldElement) Sub(other FieldElement) (*FieldElement, error) {
	if f.prime.Cmp(other.prime) != 0 {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	//num := (f.num - other.num) % f.prime

	num := big.NewInt(0).Mod(big.NewInt(0).Sub(f.num, other.num), f.prime)
	// 음수일 경우 prime을 더해줌
	if num.Cmp(big.NewInt(0)) == -1 {
		num.Add(num, f.prime)
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 곱하는 함수 (곱한 결과는 같은 유한체에 속함)
func (f FieldElement) Mul(other FieldElement) (*FieldElement, error) {
	if f.prime.Cmp(other.prime) != 0 {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := big.NewInt(0).Mod(big.NewInt(0).Mul(f.num, other.num), f.prime)
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 거듭제곱하는 함수 (거듭제곱한 결과는 같은 유한체에 속함)
func (f FieldElement) Pow(exp *big.Int) (*FieldElement, error) {
	exp.Mod(exp, big.NewInt(0).Sub(f.prime, big.NewInt(1)))
	// 지수가 음수일 경우 양수로 변환
	if exp.Cmp(big.NewInt(0)) == -1 {
		exp.Add(exp, big.NewInt(0).Sub(f.prime, big.NewInt(1)))
	}

	num := big.NewInt(0).Mod(big.NewInt(0).Exp(f.num, exp, f.prime), f.prime)
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 나누는 함수 (나눈 결과는 같은 유한체에 속함)
func (f FieldElement) Div(other FieldElement) (*FieldElement, error) {
	if f.prime.Cmp(other.prime) != 0 {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	// 페르마의 소정리를 이용하여 나눗셈을 곱셈으로 변환
	num := big.NewInt(0).Mod(big.NewInt(0).Mul(f.num, pow(other.num, big.NewInt(0).Sub(f.prime, big.NewInt(2)), f.prime)), f.prime)

	return NewFieldElement(num, f.prime)
}

// 거듭제곱을 구하는 함수 (분할정복을 이용)
func pow(num, exp, mod *big.Int) *big.Int {
	if exp.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1)
	}

	if exp.Cmp(big.NewInt(1)) == 0 {
		return num
	}

	if exp.Bit(0) == 0 {
		temp := pow(num, exp.Div(exp, big.NewInt(2)), mod)
		return big.NewInt(0).Mod(big.NewInt(0).Mul(temp, temp), mod)
	}

	temp := pow(num, exp.Div(exp.Sub(exp, big.NewInt(1)), big.NewInt(2)), mod)
	return big.NewInt(0).Mod(big.NewInt(0).Mul(big.NewInt(0).Mul(temp, temp), num), mod)
}
