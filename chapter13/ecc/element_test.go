package ecc

import (
	"math/big"
	"testing"
)

func TestElementNE(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(2), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(2), big.NewInt(31))
	c, _ := NewFieldElement(big.NewInt(15), big.NewInt(31))

	if !a.Equal(b) {
		t.Errorf("2 == 2 should be true")
	}

	if a.Equal(c) {
		t.Errorf("2 == 15 should be false")
	}

	if a.NotEqual(b) {
		t.Errorf("2 != 2 should be false")
	}
}

func TestElementAdd(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(2), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(15), big.NewInt(31))
	c, _ := NewFieldElement(big.NewInt(17), big.NewInt(31))

	added, _ := a.Add(b)

	if !added.Equal(c) {
		t.Errorf("2 + 15 should be 17")
	}

	a, _ = NewFieldElement(big.NewInt(17), big.NewInt(31))
	b, _ = NewFieldElement(big.NewInt(21), big.NewInt(31))
	c, _ = NewFieldElement(big.NewInt(7), big.NewInt(31))

	added, _ = a.Add(b)

	if !added.Equal(c) {
		t.Errorf("17 + 21 should be 7")
	}
}

func TestElementSub(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(29), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(4), big.NewInt(31))
	c, _ := NewFieldElement(big.NewInt(25), big.NewInt(31))

	subbed, _ := a.Sub(b)

	if !subbed.Equal(c) {
		t.Errorf("29 - 4 should be 25")
	}

	a, _ = NewFieldElement(big.NewInt(15), big.NewInt(31))
	b, _ = NewFieldElement(big.NewInt(30), big.NewInt(31))
	c, _ = NewFieldElement(big.NewInt(16), big.NewInt(31))

	subbed, _ = a.Sub(b)

	if !subbed.Equal(c) {
		t.Errorf("15 - 30 should be 16")
	}
}

func TestElementMul(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(24), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(19), big.NewInt(31))
	c, _ := NewFieldElement(big.NewInt(22), big.NewInt(31))

	mul, _ := a.Mul(b)

	if !mul.Equal(c) {
		t.Errorf("24 * 19 should be 22")
	}
}

func TestElementPow(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(17), big.NewInt(31))
	pow, _ := a.Pow(big.NewInt(3))

	expected, _ := NewFieldElement(big.NewInt(15), big.NewInt(31))

	if !pow.Equal(expected) {
		t.Errorf("17 ** 3 should be 15")
	}

	a, _ = NewFieldElement(big.NewInt(5), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(18), big.NewInt(31))

	pow, _ = a.Pow(big.NewInt(5))
	pow, _ = pow.Mul(b)

	expected, _ = NewFieldElement(big.NewInt(16), big.NewInt(31))

	if !pow.Equal(expected) {
		t.Errorf("5 ** 5 * 18 should be 16")
	}
}

func TestElementDiv(t *testing.T) {
	a, _ := NewFieldElement(big.NewInt(3), big.NewInt(31))
	b, _ := NewFieldElement(big.NewInt(24), big.NewInt(31))

	div, _ := a.Div(b)

	expected, _ := NewFieldElement(big.NewInt(4), big.NewInt(31))

	if !div.Equal(expected) {
		t.Errorf("3 / 24 should be 4")
	}

	a, _ = NewFieldElement(big.NewInt(17), big.NewInt(31))

	pow, _ := a.Pow(big.NewInt(-3))

	expected, _ = NewFieldElement(big.NewInt(29), big.NewInt(31))

	if !pow.Equal(expected) {
		t.Errorf("17 ** -3 should be 29")
	}

	a, _ = NewFieldElement(big.NewInt(4), big.NewInt(31))
	b, _ = NewFieldElement(big.NewInt(11), big.NewInt(31))

	pow, _ = a.Pow(big.NewInt(-4))
	pow, _ = pow.Mul(b)

	expected, _ = NewFieldElement(big.NewInt(13), big.NewInt(31))

	if !pow.Equal(expected) {
		t.Errorf("4 ** -4 should be 13")
	}
}
