package ecc

import "fmt"

type Signature interface {
	fmt.Stringer
	R() FieldElement
	S() FieldElement
}

type s256Signature struct {
	r FieldElement
	s FieldElement
}

func NewS256Signature(r, s FieldElement) Signature {
	return &s256Signature{r, s}
}

func (sig s256Signature) R() FieldElement {
	return sig.r
}

func (sig s256Signature) S() FieldElement {
	return sig.s
}

func (sig s256Signature) String() string {
	return fmt.Sprintf("Signature(%s, %s)", sig.r, sig.s)
}
