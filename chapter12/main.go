package main

import (
	"bytes"
	"chapter12/bloomfilter"
	"chapter12/utils"
	"fmt"
	"math/big"
)

func main() {
	practice6()
}

func practice1() {
	var bitFieldSize int64 = 10 // 그룹의 개수
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))
	h := utils.Hash256([]byte("hello world")) // 해시값

	n := utils.BytesToBigInt(h)                         // 해시값을 빅엔디언 정수로 변환
	bit := big.NewInt(0).Mod(n, big.NewInt(10)).Int64() // 해시값을 그룹의 개수로 나눈 나머지를 구함

	bitField[bit] = 1     // 해당 그룹의 비트를 1로 설정
	fmt.Println(bitField) // [0 0 0 0 0 0 0 0 0 1]
}

func practice2() {
	var bitFieldSize int64 = 10 // 그룹의 개수
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))

	for _, item := range [][]byte{[]byte("hello world"), []byte("goodbye")} {
		h := utils.Hash256(item) // 해시값
		n := utils.BytesToBigInt(h)
		bit := big.NewInt(0).Mod(n, big.NewInt(10)).Int64() // 해시값을 그룹의 개수로 나눈 나머지를 구함

		bitField[bit] = 1 // 해당 그룹의 비트를 1로 설정
	}

	fmt.Println(bitField) // [0 0 1 0 0 0 0 0 0 1]
}

func practice3() {
	var bitFieldSize int64 = 10 // 그룹의 개수
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))

	for _, item := range [][]byte{[]byte("hello world"), []byte("goodbye")} {
		h := utils.Hash160(item) // 해시값
		n := utils.BytesToBigInt(h)
		bit := big.NewInt(0).Mod(n, big.NewInt(10)).Int64() // 해시값을 그룹의 개수로 나눈 나머지를 구함

		bitField[bit] = 1 // 해당 그룹의 비트를 1로 설정
	}

	fmt.Println(bitField) // [1 1 0 0 0 0 0 0 0 0]
}

func practice4() {
	var bitFieldSize int64 = 10 // 그룹의 개수
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))
	hasher := []func(b []byte) []byte{utils.Hash256, utils.Hash160}

	for _, item := range [][]byte{[]byte("hello world"), []byte("goodbye")} {
		for _, h := range hasher {
			n := utils.BytesToBigInt(h(item))
			bit := big.NewInt(0).Mod(n, big.NewInt(10)).Int64() // 해시값을 그룹의 개수로 나눈 나머지를 구함

			bitField[bit] = 1 // 해당 그룹의 비트를 1로 설정
		}
	}

	fmt.Println(bitField) // [1 1 1 0 0 0 0 0 0 1]
}

func practice5() {
	fieldSize := 2
	numOfFuncs := 2
	tweak := 42

	bitFieldSize := uint32(fieldSize * 8)
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))

	for _, item := range [][]byte{[]byte("hello world"), []byte("goodbye")} {
		for i := 0; i < numOfFuncs; i++ {
			seed := i*bloomfilter.BIP37Constant + tweak
			h := utils.Murmur3(item, uint32(seed))
			bit := h % bitFieldSize
			bitField[bit] = 1
		}
	}

	fmt.Println(bitField) // [0 0 0 0 0 1 1 0 0 1 1 0 0 0 0 0]
}

func practice6() {
	fieldSize := 10
	numOfFuncs := 5
	tweak := 99

	bitFieldSize := uint32(fieldSize * 8)
	bitField := bytes.Repeat([]byte{0}, int(bitFieldSize))

	for _, item := range [][]byte{[]byte("Hello World"), []byte("Goodbye!")} {
		for i := 0; i < numOfFuncs; i++ {
			seed := i*bloomfilter.BIP37Constant + tweak
			h := utils.Murmur3(item, uint32(seed))
			bit := h % bitFieldSize
			bitField[bit] = 1
		}
	}

	b := utils.BitFieldToBytes(bitField)
	fmt.Printf("%x\n", b) // 4000600a080000010940
}

func practice7() {}
