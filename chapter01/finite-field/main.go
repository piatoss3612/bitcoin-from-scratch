package main

import (
	"fmt"
)

type FieldElement struct {
	num   int
	prime int
}

func New(num, prime int) (*FieldElement, error) {
	if num >= prime || num < 0 {
		return nil, fmt.Errorf("num %d not in field range 0 to %d", num, prime-1)
	}

	return &FieldElement{num, prime}, nil
}

func (f FieldElement) String() string {
	return fmt.Sprintf("FieldElement_%d(%d)", f.prime, f.num)
}

func (f FieldElement) Equals(other FieldElement) bool {
	return f.num == other.num && f.prime == other.prime
}

func (f FieldElement) Add(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num + other.num) % f.prime
	return New(num, f.prime)
}

func (f FieldElement) Sub(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num - other.num) % f.prime
	if num < 0 {
		num += f.prime
	}
	return New(num, f.prime)
}

func (f FieldElement) Mul(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num * other.num) % f.prime
	return New(num, f.prime)
}

func (f FieldElement) Pow(exp int) (*FieldElement, error) {
	if exp < 0 {
		exp = exp%(f.prime-1) + (f.prime - 1)
	}

	num := pow(f.num, exp, f.prime) % f.prime
	return New(num, f.prime)
}

func (f FieldElement) Div(other FieldElement) (*FieldElement, error) {
	if f.prime != other.prime {
		return nil, fmt.Errorf("cannot add two numbers in different Fields %d, %d", f.prime, other.prime)
	}

	num := (f.num * pow(other.num, f.prime-2, f.prime)) % f.prime
	return New(num, f.prime)
}

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

func main() {
	a, _ := New(7, 13)
	b, _ := New(6, 13)

	fmt.Println(a.Equals(*b))
	fmt.Println(a.Equals(*a))

	c, _ := a.Add(*b)
	fmt.Println(c)

	d, _ := a.Sub(*b)
	fmt.Println(d)

	e, _ := a.Mul(*b)
	fmt.Println(e)

	f, _ := a.Pow(-3)
	fmt.Println(f)

	g, _ := a.Div(*b)
	fmt.Println(g)
}
