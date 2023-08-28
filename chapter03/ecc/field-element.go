package ecc

import "fmt"

// 유한체의 원소를 나타내는 구조체
type FieldElement struct {
	num   int
	prime int
}

// 유한체의 원소를 생성하는 함수
func NewFieldElement(num, prime int) (*FieldElement, error) {
	if num >= prime || num < 0 {
		return nil, fmt.Errorf("num %d not in field range 0 to %d", num, prime-1)
	}

	return &FieldElement{num, prime}, nil
}

// 유한체의 원소를 문자열로 표현하는 함수 (Stringer 인터페이스 구현)
func (f FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.num)
}

// 유한체의 원소가 같은 유한체에 속하고 같은 값을 가지는지 확인하는 함수
func (f FieldElement) Equal(other FieldElement) bool {
	return f.num == other.num && f.prime == other.prime
}

func (f FieldElement) NotEqual(other FieldElement) bool {
	return f.num != other.num || f.prime != other.prime
}

// 유한체의 원소를 더하는 함수 (더한 결과는 같은 유한체에 속함)
func (f FieldElement) Add(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num + other.num) % f.prime
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 빼는 함수 (뺀 결과는 같은 유한체에 속함)
func (f FieldElement) Sub(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num - other.num) % f.prime
	// 음수일 경우 prime을 더해줌
	if num < 0 {
		num += f.prime
	}
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 곱하는 함수 (곱한 결과는 같은 유한체에 속함)
func (f FieldElement) Mul(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num * other.num) % f.prime
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 거듭제곱하는 함수 (거듭제곱한 결과는 같은 유한체에 속함)
func (f FieldElement) Pow(exp int) (*FieldElement, error) {
	exp %= (f.prime - 1)
	// 지수가 음수일 경우 양수로 변환
	if exp < 0 {
		exp += (f.prime - 1)
	}

	num := pow(f.num, exp, f.prime) % f.prime
	return NewFieldElement(num, f.prime)
}

// 유한체의 원소를 나누는 함수 (나눈 결과는 같은 유한체에 속함)
func (f FieldElement) Div(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	// 페르마의 소정리를 이용하여 나눗셈을 곱셈으로 변환
	num := (f.num * pow(other.num, f.prime-2, f.prime)) % f.prime
	return NewFieldElement(num, f.prime)
}

// 거듭제곱을 구하는 함수 (분할정복을 이용)
func pow(num, exp, mod int) int {
	if exp == 0 {
		return 1
	}

	if exp == 1 {
		return num
	}

	if exp%2 == 0 {
		temp := pow(num, exp/2, mod)
		return (temp * temp) % mod
	}

	temp := pow(num, (exp-1)/2, mod)
	return (temp * temp * num) % mod
}
