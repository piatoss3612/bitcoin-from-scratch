package block

import (
	"chapter09/utils"
	"encoding/hex"
	"errors"
	"math/big"
)

type Block struct {
	Version    int
	PrevBlock  string
	MerkleRoot string
	Timestamp  int
	Bits       int
	Nonce      int
}

func New(version int, prevBlock, merkleRoot string, timestamp, bits, nonce int) *Block {
	return &Block{
		Version:    version,
		PrevBlock:  prevBlock,
		MerkleRoot: merkleRoot,
		Timestamp:  timestamp,
		Bits:       bits,
		Nonce:      nonce,
	}
}

// 블록을 직렬화하는 함수
func (b Block) Serialize() ([]byte, error) {
	result := make([]byte, 0, 80)

	version := utils.IntToLittleEndian(b.Version, 4)     // version 4바이트 리틀엔디언
	prevBlockBytes, err := hex.DecodeString(b.PrevBlock) // 16진수 문자열을 []byte로 변환
	if err != nil {
		return nil, err
	}
	prevBlock := utils.ReverseBytes(prevBlockBytes)        // prevBlock 32바이트 리틀엔디언
	merkleRootBytes, err := hex.DecodeString(b.MerkleRoot) // 16진수 문자열을 []byte로 변환
	if err != nil {
		return nil, err
	}
	merkleRoot := utils.ReverseBytes(merkleRootBytes)    // merkleRoot 32바이트 리틀엔디언
	timestamp := utils.IntToLittleEndian(b.Timestamp, 4) // timestamp 4바이트 리틀엔디언
	bits := utils.IntToBytes(b.Bits, 4)                  // bits 4바이트 빅엔디언
	nonce := utils.IntToBytes(b.Nonce, 4)                // nonce 4바이트 빅엔디언

	totalLength := len(version) + len(prevBlock) + len(merkleRoot) + len(timestamp) + len(bits) + len(nonce)

	if totalLength > 80 {
		return nil, errors.New("the size of block is too big")
	}

	result = append(result, version...)
	result = append(result, prevBlock...)
	result = append(result, merkleRoot...)
	result = append(result, timestamp...)
	result = append(result, bits...)
	result = append(result, nonce...)

	return result, nil
}

// 블록의 해시를 계산하는 함수
func (b Block) Hash() ([]byte, error) {
	s, err := b.Serialize()
	if err != nil {
		return nil, err
	}
	return utils.ReverseBytes(utils.Hash256(s)), nil
}

func (b Block) Bip9() bool {
	return b.Version>>29 == 0x001
}

func (b Block) Bip91() bool {
	return b.Version>>4&1 == 1
}

func (b Block) Bip141() bool {
	return b.Version>>1&1 == 1
}

// 블록의 난이도를 계산하는 함수
func (b Block) Difficulty() *big.Float {
	target := BitsToTarget(utils.IntToBytes(b.Bits, 4))
	return big.NewFloat(0).Mul(big.NewFloat(0xffff), big.NewFloat(0).Quo(big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(256), big.NewInt(0).Sub(big.NewInt(0x1d), big.NewInt(3)), nil)), big.NewFloat(0).SetInt(target))) // 0xffff * 256^(0x1d - 3) / target
}
