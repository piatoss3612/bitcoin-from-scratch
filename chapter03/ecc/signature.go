package ecc

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

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
	return fmt.Sprintf("Signature(%s, %s)", sig.r.Num().Text(16), sig.s.Num().Text(16))
}

type PrivateKey interface {
	fmt.Stringer
	Sign(z FieldElement) (Signature, error)
}

type s256PrivateKey struct {
	secret FieldElement
	point  Point
}

func NewS256PrivateKey(secret FieldElement) (PrivateKey, error) {
	point, err := G.Mul(secret.Num())
	if err != nil {
		return nil, err
	}
	return &s256PrivateKey{secret, point}, nil
}

func (pvk s256PrivateKey) Sign(z FieldElement) (Signature, error) {
	bigK, err := rand.Int(rand.Reader, N)
	if err != nil {
		return nil, err
	}

	kG, err := G.Mul(bigK)
	if err != nil {
		return nil, err
	}

	rx := kG.X()

	r, err := NewFieldElement(rx.Num(), N)
	if err != nil {
		return nil, err
	}

	k, err := NewFieldElement(bigK, N)
	if err != nil {
		return nil, err
	}

	kInv, err := k.Pow(big.NewInt(0).Sub(N, big.NewInt(2)))
	if err != nil {
		return nil, err
	}

	secretR, err := r.Mul(pvk.secret)
	if err != nil {
		return nil, err
	}

	zPlusSecretR, err := z.Add(secretR)
	if err != nil {
		return nil, err
	}

	s, err := zPlusSecretR.Mul(kInv)
	if err != nil {
		return nil, err
	}

	if s.Num().Cmp(big.NewInt(0).Div(N, big.NewInt(2))) == 1 {
		nsNum := big.NewInt(0).Sub(N, s.Num())
		ns, err := NewFieldElement(nsNum, N)
		if err != nil {
			return nil, err
		}

		return NewS256Signature(r, ns), nil
	}

	return NewS256Signature(r, s), nil
}

func (pvk s256PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey(%s)", pvk.secret.Num().Text(16))
}
