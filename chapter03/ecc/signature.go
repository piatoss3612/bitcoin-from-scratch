package ecc

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
	"math/big"
)

// 서명 인터페이스
type Signature interface {
	fmt.Stringer
	R() *big.Int
	S() *big.Int
}

// secp256k1 서명 구조체
type s256Signature struct {
	r *big.Int
	s *big.Int
}

// secp256k1 서명을 생성하는 함수
func NewS256Signature(r, s *big.Int) Signature {
	return &s256Signature{r, s}
}

// secp256k1 서명의 r값을 반환하는 함수
func (sig s256Signature) R() *big.Int {
	return sig.r
}

// secp256k1 서명의 s값을 반환하는 함수
func (sig s256Signature) S() *big.Int {
	return sig.s
}

// secp256k1 서명을 문자열로 반환하는 함수 (Stringer 인터페이스 구현)
func (sig s256Signature) String() string {
	return fmt.Sprintf("Signature(%s, %s)", sig.r.Text(16), sig.s.Text(16))
}

// 개인키 인터페이스
type PrivateKey interface {
	fmt.Stringer
	Sign(z *big.Int) (Signature, error)
}

// secp256k1 개인키 구조체
type s256PrivateKey struct {
	secret *big.Int
	point  Point
}

// secp256k1 개인키를 생성하는 함수
func NewS256PrivateKey(secret *big.Int) (PrivateKey, error) {
	point, err := G.Mul(secret)
	if err != nil {
		return nil, err
	}
	return &s256PrivateKey{secret, point}, nil
}

// secp256k1 개인키로 서명을 생성하는 함수
func (pvk s256PrivateKey) Sign(z *big.Int) (Signature, error) {
	k, err := pvk.deterministicK(z) // RFC6979 표준에 따라 k값을 생성
	if err != nil {
		return nil, err
	}

	kG, err := G.Mul(k) // kG를 계산
	if err != nil {
		return nil, err
	}

	r := kG.X().Num() // kG의 x좌표를 r값으로 사용

	kInv := big.NewInt(0).ModInverse(k, N) // k의 역원

	secretR := big.NewInt(0).Mod(big.NewInt(0).Mul(pvk.secret, r), N) // 개인키와 r값의 곱을 계산

	zPlusSecretR := big.NewInt(0).Mod(big.NewInt(0).Add(z, secretR), N) // z + secretR을 계산

	s := big.NewInt(0).Mod(big.NewInt(0).Mul(kInv, zPlusSecretR), N) // kInv * (z + secretR)을 계산

	// s가 N/2보다 큰 경우, s = N - s로 사용
	if s.Cmp(big.NewInt(0).Div(N, big.NewInt(2))) == 1 {
		ns := big.NewInt(0).Sub(N, s)
		return NewS256Signature(r, ns), nil
	}

	// 서명 생성
	return NewS256Signature(r, s), nil
}

// secp256k1 개인키를 문자열로 반환하는 함수 (Stringer 인터페이스 구현)
func (pvk s256PrivateKey) String() string {
	return fmt.Sprintf("PrivateKey(%s)", pvk.secret.Text(16))
}

// RFC6979 표준에 따라 k값을 생성하는 함수
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
