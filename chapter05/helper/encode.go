package helper

import (
	"math/big"
	"strings"
	"unsafe"
)

var base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz" // base58 인코딩에 사용할 문자열

// 바이트 슬라이스를 base58로 인코딩하는 함수
func EncodeBase58(s []byte) string {
	// 앞에 몇 바이트가 0x00인지 확인
	count := 0
	for _, c := range s {
		if c == 0x00 {
			count++
		} else {
			break
		}
	}

	n := big.NewInt(0).SetBytes(s) // 바이트 슬라이스를 big.Int로 변환

	// 0x00이 몇개 있는지에 따라 0x01을 count만큼 반복하여 prefix를 만듬
	// 이 prefix는 pay-to-pubkey-hash(P2PKH)에서 필요함 (6장에서 설명)
	prefix := strings.Repeat("1", count)

	result := strings.Builder{}

	// n을 58로 나눈 나머지에 해당하는 문자를 base58Alphabet에서 찾아서 문자열을 만듦
	for n.Cmp(big.NewInt(0)) == 1 {
		mod := big.NewInt(0)
		n.DivMod(n, big.NewInt(58), mod)
		result.WriteByte(base58Alphabet[mod.Int64()])
	}

	// prefix를 붙임
	result.WriteString(prefix)

	// 문자열을 뒤집음
	resultBytes := ReverseBytes(StringToBytes(result.String()))

	return BytesToString((resultBytes))
}

// 바이트 슬라이스를 뒤집는 함수
func ReverseBytes(b []byte) []byte {
	result := make([]byte, len(b))

	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = b[j], b[i]
	}

	return result
}

// 바이트 슬라이스의 체크섬을 추가하여 base58로 인코딩하는 함수
func EncodeBase58Checksum(b []byte) string {
	return EncodeBase58(append(b, Hash256(b)[:4]...))
}

// 리틀엔디언으로 인코딩된 바이트 슬라이스를 big.Int로 변환하는 함수
func LittleEndianToBigInt(b []byte) *big.Int {
	return BytesToBigInt(ReverseBytes(b))
}

// big.Int를 리틀엔디언으로 인코딩된 바이트 슬라이스로 변환하는 함수
func BigIntToLittleEndian(n *big.Int) []byte {
	return ReverseBytes(n.Bytes())
}

// 바이트 슬라이스를 big.Int로 변환하는 함수
func BytesToBigInt(b []byte) *big.Int {
	return big.NewInt(0).SetBytes(b)
}

// 문자열을 바이트 슬라이스로 변환하는 함수
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// 바이트 슬라이스를 문자열로 변환하는 함수
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
