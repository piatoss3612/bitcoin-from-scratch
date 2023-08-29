package ecc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
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
	/*
		bigK, err := rand.Int(rand.Reader, N)
		if err != nil {
			return nil, err
		}
	*/
	bigK, err := pvk.deterministicK(z)
	if err != nil {
		return nil, err
	}

	fmt.Println(bigK)

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

// reference: https://github.com/codahale/rfc6979/blob/master/rfc6979.go
func (pvk s256PrivateKey) deterministicK(z FieldElement) (*big.Int, error) {
	k := bytes.Repeat([]byte{0x00}, 32)
	v := bytes.Repeat([]byte{0x01}, 32)

	zNum := z.Num()

	if zNum.Cmp(N) == 1 {
		zNum = big.NewInt(0).Sub(zNum, N)
	}

	zBytes := zNum.Bytes()
	secreteBytes := pvk.secret.Num().Bytes()

	alg := sha256.New

	k = pvk.mac(alg, k, append(append(v, 0x00), append(secreteBytes, zBytes...)...), k)
	v = pvk.mac(alg, k, v, v)

	k = pvk.mac(alg, k, append(append(v, 0x01), append(secreteBytes, zBytes...)...), k)
	v = pvk.mac(alg, k, v, v)

	for {
		v = pvk.mac(alg, k, v, v)
		candidate := big.NewInt(0).SetBytes(v)

		if candidate.Cmp(big.NewInt(0)) == 1 && candidate.Cmp(N) == -1 {
			return candidate, nil
		}

		k = pvk.mac(alg, k, append(v, 0x00), k)
		v = pvk.mac(alg, k, v, v)
	}
}

func (pvk s256PrivateKey) mac(alg func() hash.Hash, k, m, buf []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(buf[:0])
}
