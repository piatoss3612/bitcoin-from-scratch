package block

import (
	"bytes"
	"chapter09/utils"
	"encoding/hex"
	"errors"
	"math/big"
)

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

	return New(version, prevBlock, merkleRoot, timestamp, bits, nonce), nil
}

// 목푯값을 계산하는 함수
func BitsToTarget(b []byte) *big.Int {
	exp := utils.BytesToBigInt(b[len(b)-1:])         // 지수
	coef := utils.LittleEndianToBigInt(b[:len(b)-1]) // 계수

	return big.NewInt(0).Mul(coef, big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0).Sub(exp, big.NewInt(3)), nil)) // 계수 * 256^(지수-3) = 목푯값
}
