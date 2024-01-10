package merkleblock

import (
	"bytes"
	"chapter13/utils"
)

type MerkleBlock struct {
	Version       int
	PrevBlockHash []byte
	MerkleRoot    []byte
	Timestamp     int
	Bits          int
	Nonce         int
	TotalTx       int
	Hashes        [][]byte
	Flags         []byte
}

func (m *MerkleBlock) Parse(b []byte) error {
	buf := bytes.NewBuffer(b)
	m.Version = utils.LittleEndianToInt(buf.Next(4))
	m.PrevBlockHash = utils.ReverseBytes(buf.Next(32))
	m.MerkleRoot = utils.ReverseBytes(buf.Next(32))
	m.Timestamp = utils.LittleEndianToInt(buf.Next(4))
	m.Bits = utils.BytesToInt(buf.Next(4))
	m.Nonce = utils.BytesToInt(buf.Next(4))
	m.TotalTx = utils.LittleEndianToInt(buf.Next(4))
	hashCount, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	m.Hashes = make([][]byte, hashCount)
	for i := 0; i < int(hashCount); i++ {
		m.Hashes[i] = utils.ReverseBytes(buf.Next(32))
	}

	flagCount, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	m.Flags = buf.Next(int(flagCount))

	return nil
}

func (m *MerkleBlock) IsValid() (bool, error) {
	flagBits := utils.BytesToBitField(m.Flags)
	hashes := m.Hashes
	for i, hash := range hashes {
		hashes[i] = utils.ReverseBytes(hash)
	}

	mt := NewMerkleTree(m.TotalTx)
	err := mt.PopulateTree(flagBits, hashes)
	if err != nil {
		return false, err
	}

	return bytes.Equal(utils.ReverseBytes(mt.Root()), m.MerkleRoot), nil
}
