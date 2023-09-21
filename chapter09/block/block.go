package block

import "chapter09/utils"

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
func (b *Block) Serialize() []byte {
	result := make([]byte, 80)
	result = append(result, utils.IntToLittleEndian(b.Version, 4)...)                 // version 4바이트 리틀엔디언
	result = append(result, utils.ReverseBytes(utils.StringToBytes(b.PrevBlock))...)  // prevBlock 32바이트 리틀엔디언
	result = append(result, utils.ReverseBytes(utils.StringToBytes(b.MerkleRoot))...) // merkleRoot 32바이트 리틀엔디언
	result = append(result, utils.IntToLittleEndian(b.Timestamp, 4)...)               // timestamp 4바이트 리틀엔디언
	result = append(result, utils.IntToBytes(b.Bits, 4)...)                           // bits 4바이트 빅엔디언
	result = append(result, utils.IntToBytes(b.Nonce, 4)...)                          // nonce 4바이트 빅엔디언
	return result
}

// TODO: implement Hash
func (b *Block) Hash() []byte {
	return nil
}
