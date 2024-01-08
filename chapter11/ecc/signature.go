package ecc

import (
	"bytes"
	"chapter11/utils"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"hash"
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

// secp256k1 개인키 구조체
type s256PrivateKey struct {
	secret []byte
	point  Point
}

// secp256k1 개인키를 생성하는 함수
func NewS256PrivateKey(secret []byte) (PrivateKey, error) {
	e := utils.BytesToBigInt(secret)

	point, err := G.Mul(e)
	if err != nil {
		return nil, err
	}

	return &s256PrivateKey{secret, point}, nil
}

// 기존 secp256k1 개인키를 문자열로 반환하는 함수 (Stringer 인터페이스 구현)
func (pvk s256PrivateKey) String() string {
	return fmt.Sprintf("R = %s", pvk.point)
}

// secp256k1 개인키로 서명을 생성하는 함수
func (pvk s256PrivateKey) Sign(z []byte) (Signature, error) {
	bigZ := utils.BytesToBigInt(z) // 서명할 메시지를 big.Int로 변환

	e := utils.BytesToBigInt(pvk.secret) // 개인키를 big.Int로 변환

	k, err := pvk.deterministicK(bigZ) // RFC6979 표준에 따라 k값을 생성
	if err != nil {
		return nil, err
	}

	kG, err := G.Mul(k) // kG를 계산
	if err != nil {
		return nil, err
	}

	r := kG.X().Num() // kG의 x좌표를 r값으로 사용

	kInv := invBN(k, N) // k의 역원

	s := mulBN(addBN(bigZ, mulBN(r, e, N), N), kInv, N) // s = (z + r * pvk.secret) * k^-1 mod N

	// s가 N/2보다 큰 경우, s = N - s로 사용
	if s.Cmp(big.NewInt(0).Div(N, big.NewInt(2))) == 1 {
		ns := big.NewInt(0).Sub(N, s)
		return NewS256Signature(r, ns), nil
	}

	// 서명 생성
	return NewS256Signature(r, s), nil
}

// RFC6979 표준에 따라 k값을 생성하는 함수
// reference: https://github.com/codahale/rfc6979/blob/master/rfc6979.go
func (pvk s256PrivateKey) deterministicK(z *big.Int) (*big.Int, error) {
	k := bytes.Repeat([]byte{0x00}, 32)
	v := bytes.Repeat([]byte{0x01}, 32)

	if z.Cmp(N) == 1 {
		z.Sub(z, N)
	}

	zBytes := z.FillBytes(make([]byte, 32))
	secreteBytes := big.NewInt(0).SetBytes(pvk.secret).FillBytes(make([]byte, 32)) // 개인키를 big.Int로 변환한 뒤, 32바이트로 채움

	alg := sha256.New

	k = pvk.mac(alg, k, append(append(v, 0x00), append(secreteBytes, zBytes...)...))
	v = pvk.mac(alg, k, v)

	k = pvk.mac(alg, k, append(append(v, 0x01), append(secreteBytes, zBytes...)...))
	v = pvk.mac(alg, k, v)

	for {
		v = pvk.mac(alg, k, v)
		candidate := big.NewInt(0).SetBytes(v)

		if candidate.Cmp(big.NewInt(1)) >= 0 && candidate.Cmp(N) == -1 {
			return candidate, nil
		}

		k = pvk.mac(alg, k, append(v, 0x00))
		v = pvk.mac(alg, k, v)
	}
}

func (pvk s256PrivateKey) mac(alg func() hash.Hash, k, m []byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// secp256k1 개인키의 점을 반환하는 함수
func (pvk s256PrivateKey) Point() Point {
	return pvk.point
}

// secp256k1 개인키의 WIF 형식을 반환하는 함수
func (pvk s256PrivateKey) WIF(compressed bool, testnet bool) string {
	secret := pvk.secret

	// secret의 길이가 32보다 작으면 비어있는 길이만큼 0x00을 추가
	if len(secret) < 32 {
		secret = append(make([]byte, 32-len(secret)), secret...)
	}

	// 압축된 공개키를 사용하는 경우, secret에 0x01을 추가
	if compressed {
		secret = append(secret, 0x01)
	}

	// testnet을 사용하는 경우, secret에 0xef를 추가
	// mainnet을 사용하는 경우, secret에 0x80을 추가
	if testnet {
		secret = append([]byte{0xef}, secret...)
	} else {
		secret = append([]byte{0x80}, secret...)
	}

	return utils.EncodeBase58Checksum(secret) // Base58Checksum 인코딩
}
