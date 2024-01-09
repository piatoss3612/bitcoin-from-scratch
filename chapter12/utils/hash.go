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

func Murmur3(b []byte, seeds ...uint32) uint32 {
	var seed uint32 = 0
	if len(seeds) > 0 {
		seed = seeds[0]
	}

	c1, c2 := uint32(0xcc9e2d51), uint32(0x1b873593)
	length := len(b)

	h1 := seed
	roundedEnd := length & 0xfffffffc

	for i := 0; i < roundedEnd; i += 4 {
		k1 := uint32(b[i]) | uint32(b[i+1])<<8 | uint32(b[i+2])<<16 | uint32(b[i+3])<<24

		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1

		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}

	var k1 uint32

	switch length & 0x03 {
	case 3:
		k1 ^= uint32(b[roundedEnd+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(b[roundedEnd+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(b[roundedEnd])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1
	}

	h1 ^= uint32(length)
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
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
