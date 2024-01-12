package block

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	rawBlock, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(rawBlock)
	if err != nil {
		t.Fatal(err)
	}

	if block.Version != 0x20000002 {
		t.Errorf("Version = %d, expected = %d", block.Version, 0x20000002)
	}

	prevBlock := "000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e"

	if !strings.EqualFold(block.PrevBlock, prevBlock) {
		t.Errorf("PrevBlock = %x, expected = %x", block.PrevBlock, prevBlock)
	}

	merkleRoot := "be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b"

	if !strings.EqualFold(block.MerkleRoot, merkleRoot) {
		t.Errorf("MerkleRoot = %x, expected = %x", block.MerkleRoot, merkleRoot)
	}

	if block.Timestamp != 0x59a7771e {
		t.Errorf("Timestamp = %d, expected = %d", block.Timestamp, 0x59a7771e)
	}

	bits := big.NewInt(0).SetBytes([]byte{0xe9, 0x3c, 0x01, 0x18})

	if block.Bits != int(bits.Int64()) {
		t.Errorf("Bits = %x, expected = %x", block.Bits, bits)
	}

	nonce := big.NewInt(0).SetBytes([]byte{0xa4, 0xff, 0xd7, 0x1d})

	if block.Nonce != int(nonce.Int64()) {
		t.Errorf("Nonce = %x, expected = %x", block.Nonce, nonce)
	}
}
