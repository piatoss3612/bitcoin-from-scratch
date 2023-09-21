package block

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

// TODO: implement Parse, Serialize, Hash
func Parse(b []byte) (*Block, error) {
	return nil, nil
}

func (b *Block) Serialize() []byte {
	return nil
}

func (b *Block) Hash() []byte {
	return nil
}
