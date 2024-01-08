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

func MerkleParent(left, right []byte) []byte {
	return Hash256(append(left, right...))
}

func MerkleParentLevel(hashes [][]byte) [][]byte {
	if len(hashes)%2 == 1 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	parentLevel := make([][]byte, len(hashes)/2)
	for i := 0; i < len(hashes); i += 2 {
		parentLevel[i/2] = MerkleParent(hashes[i], hashes[i+1])
	}
	return parentLevel
}

func MerkleRoot(hashes [][]byte) []byte {
	for len(hashes) > 1 {
		hashes = MerkleParentLevel(hashes)
	}
	return hashes[0]
}
