package ecc

import (
	"bytes"
	"fmt"
	"math/big"
)

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

// secp256k1 서명을 DER 형식으로 반환하는 함수
func (sig s256Signature) DER() []byte {
	r := sig.r.Bytes()
	s := sig.s.Bytes()

	// r, s의 비어있는 바이트를 제거
	r = bytes.TrimLeftFunc(r, func(r rune) bool {
		return r == 0x00
	})
	s = bytes.TrimLeftFunc(s, func(r rune) bool {
		return r == 0x00
	})

	// r의 첫번째 바이트가 0x80 이상이면 0x00을 추가
	if r[0]&0x80 != 0 {
		r = append([]byte{0x00}, r...)
	}

	// s의 첫번째 바이트가 0x80 이상이면 0x00을 추가
	if s[0]&0x80 != 0 {
		s = append([]byte{0x00}, s...)
	}

	// r, s의 길이를 1바이트로 표현
	rLen := byte(len(r))
	sLen := byte(len(s))

	r = append([]byte{0x02, rLen}, r...)
	s = append([]byte{0x02, sLen}, s...)

	// r, s를 연결
	result := append(r, s...)

	// DER 형식의 서명을 반환
	return append([]byte{0x30, byte(len(result))}, result...)
}
