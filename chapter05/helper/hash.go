package helper

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// sha256(sha256(b))를 구하는 함수
func Hash256(b []byte) []byte {
	h1 := sha256.New()
	_, _ = h1.Write(b)
	intermediateHash := h1.Sum(nil)
	h2 := sha256.New()
	_, _ = h2.Write(intermediateHash)
	return h2.Sum(nil)
}

// ripemd160(sha256(b))를 구하는 함수
func Hash160(b []byte) []byte {
	// sha256 해시값을 구함
	h256 := sha256.New()
	_, _ = h256.Write(b)
	hash1 := h256.Sum(nil)

	// sha256 해시값을 사용하여 ripemd160 해시값을 구함
	ripemd160 := ripemd160.New()
	_, _ = ripemd160.Write(hash1)

	return ripemd160.Sum(nil)
}
