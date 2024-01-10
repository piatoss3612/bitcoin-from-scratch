package block

import (
	"bytes"
	"chapter13/utils"
	"encoding/hex"
	"errors"
	"math/big"
)

const TWO_WEEK = 60 * 60 * 24 * 14 // 2주

var MaxTarget = big.NewInt(0).Mul(big.NewInt(0xffff), big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0x1d-3), nil)) // ffff0000000000000000000000000000000000000000000000000000

// 블록을 파싱하는 함수
func Parse(b []byte) (*Block, error) {
	if len(b) < 80 {
		return nil, errors.New("Block is too short")
	}

	buf := bytes.NewBuffer(b)

	version := utils.LittleEndianToInt(buf.Next(4))                    // 4바이트 리틀엔디언 정수
	prevBlock := hex.EncodeToString(utils.ReverseBytes(buf.Next(32)))  // 32바이트 리틀엔디언 해시
	merkleRoot := hex.EncodeToString(utils.ReverseBytes(buf.Next(32))) // 32바이트 리틀엔디언 해시
	timestamp := utils.LittleEndianToInt(buf.Next(4))                  // 4바이트 리틀엔디언 정수
	bits := utils.BytesToInt(buf.Next(4))                              // 4바이트 리틀엔디언 정수
	nonce := utils.BytesToInt(buf.Next(4))                             // 4바이트 리틀엔디언 정수

	return New(version, prevBlock, merkleRoot, timestamp, bits, nonce, nil), nil
}

// 목푯값을 계산하는 함수
func BitsToTarget(b []byte) *big.Int {
	exp := utils.BytesToBigInt(b[len(b)-1:])         // 지수
	coef := utils.LittleEndianToBigInt(b[:len(b)-1]) // 계수

	return big.NewInt(0).Mul(coef, big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0).Sub(exp, big.NewInt(3)), nil)) // 계수 * 256^(지수-3) = 목푯값
}

// 목푯값을 비트로 변환하는 함수
func TargetToBits(target *big.Int) []byte {
	rawBytes := target.Bytes() // 목푯값을 []byte로 변환, 앞에 0은 제외됨

	var exp int
	var coef []byte

	// 만약 rawBytes가 1로 시작하면 음수가 되므로 변환해줌
	if rawBytes[0] > 0x7f {
		exp = len(rawBytes) + 1                      // 0x00을 추가했으므로 지수는 1 증가
		coef = append([]byte{0x00}, rawBytes[:2]...) // 0x00을 추가해줌
	} else {
		exp = len(rawBytes) // 지수
		coef = rawBytes[:3] // 계수
	}

	return append(utils.ReverseBytes(coef), byte(exp)) // 계수를 리틀엔디언으로 변환하고 지수를 뒤에 붙임
}

// 새로운 목푯값을 계산하는 함수
func CalculateNewBits(prevBits []byte, timeDiff int64) []byte {
	if timeDiff > TWO_WEEK*4 {
		timeDiff = TWO_WEEK * 4
	} else if timeDiff < TWO_WEEK/4 {
		timeDiff = TWO_WEEK / 4
	}

	newTarget := big.NewInt(0).Div(big.NewInt(0).Mul(BitsToTarget(prevBits), big.NewInt(timeDiff)), big.NewInt(TWO_WEEK))

	if newTarget.Cmp(MaxTarget) == 1 {
		newTarget = MaxTarget
	}

	return TargetToBits(newTarget)
}
