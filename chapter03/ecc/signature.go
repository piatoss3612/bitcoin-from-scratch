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
	R() *big.Int
	S() *big.Int
}

type s256Signature struct {
	r *big.Int
	s *big.Int
}

func NewS256Signature(r, s *big.Int) Signature {
	return &s256Signature{r, s}
}

func (sig s256Signature) R() *big.Int {
	return sig.r
}

func (sig s256Signature) S() *big.Int {
	return sig.s
}

func (sig s256Signature) String() string {
	return fmt.Sprintf("Signature(%s, %s)", sig.r.Text(16), sig.s.Text(16))
}

type PrivateKey interface {
	fmt.Stringer
	Sign(z *big.Int) (Signature, error)
}

type s256PrivateKey struct {
	secret *big.Int
	point  Point
}

func NewS256PrivateKey(secret *big.Int) (PrivateKey, error) {
	point, err := G.Mul(secret)
	if err != nil {
		return nil, err
	}
	return &s256PrivateKey{secret, point}, nil
}

func (pvk s256PrivateKey) Sign(z *big.Int) (Signature, error) {
	k, err := pvk.deterministicK(z)
	if err != nil {
		return nil, err
	}

	kG, err := G.Mul(k)
	if err != nil {
		return nil, err
	}

	r := kG.X().Num()

	kInv := big.NewInt(0).ModInverse(k, N)

	secretR := big.NewInt(0).Mod(big.NewInt(0).Mul(pvk.secret, r), N)

	zPlusSecretR := big.NewInt(0).Mod(big.NewInt(0).Add(z, secretR), N)

	s := big.NewInt(0).Mod(big.NewInt(0).Mul(kInv, zPlusSecretR), N)

	if s.Cmp(big.NewInt(0).Div(N, big.NewInt(2))) == 1 {
		ns := big.NewInt(0).Sub(N, s)
		return NewS256Signature(r, ns), nil
	}

	return NewS256Signature(r, s), nil
}

func (pvk s256PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey(%s)", pvk.secret.Text(16))
}

// reference: https://github.com/codahale/rfc6979/blob/master/rfc6979.go
func (pvk s256PrivateKey) deterministicK(z *big.Int) (*big.Int, error) {
	k := bytes.Repeat([]byte{0x00}, 32)
	v := bytes.Repeat([]byte{0x01}, 32)

	if z.Cmp(N) == 1 {
		z.Sub(z, N)
	}

	zBytes := z.Bytes()
	secreteBytes := pvk.secret.Bytes()

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
