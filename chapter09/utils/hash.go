package utils

import (
	"crypto/sha1"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// sha256(sha256(b))를 구하는 함수
func Hash256(b []byte) []byte {
	return Sha256(Sha256(b))
}

// ripemd160(sha256(b))를 구하는 함수
func Hash160(b []byte) []byte {
	return Ripemd160(Sha256(b))
}

func Ripemd160(b []byte) []byte {
	h := ripemd160.New()
	_, _ = h.Write(b)
	return h.Sum(nil)
}

func Sha1(b []byte) []byte {
	h := sha1.New()
	_, _ = h.Write(b)
	return h.Sum(nil)
}

func Sha256(b []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(b)
	return h.Sum(nil)
}

// hash160 값을 p2pkh 주소로 변환하는 함수
func H160ToP2pkhAddress(h160 []byte, testnet bool) string {
	var prefix byte
	if testnet {
		prefix = 0x6f
	} else {
		prefix = 0x00
	}
	return EncodeBase58Checksum(append([]byte{prefix}, h160...))
}

// hash160 값을 p2sh 주소로 변환하는 함수
func H160ToP2shAddress(h160 []byte, testnet bool) string {
	var prefix byte
	if testnet {
		prefix = 0xc4
	} else {
		prefix = 0x05
	}
	return EncodeBase58Checksum(append([]byte{prefix}, h160...))
}
