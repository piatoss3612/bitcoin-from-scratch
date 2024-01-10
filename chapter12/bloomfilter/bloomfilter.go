package bloomfilter

import (
	"chapter12/network"
	"chapter12/utils"
)

const BIP37Constant = 0xfba4c795

type BloomFilter struct {
	Size     uint32
	Funcs    uint32
	Tweak    uint32
	BitField []byte
}

func New(size, funcs, tweak uint32) *BloomFilter {
	return &BloomFilter{
		Size:     size,
		Funcs:    funcs,
		Tweak:    tweak,
		BitField: make([]byte, size*8),
	}
}

func (b *BloomFilter) Add(data []byte) {
	for i := uint32(0); i < b.Funcs; i++ {
		seed := i*BIP37Constant + b.Tweak
		h := utils.Murmur3(data, seed)
		bit := h % (b.Size * 8)
		b.BitField[bit] = 1
	}
}

func (b *BloomFilter) Filterload(flag ...bool) *network.GenericMessage {
	payload := utils.EncodeVarint(int(b.Size))
	payload = append(payload, b.FilterBytes()...)
	payload = append(payload, utils.IntToLittleEndian(int(b.Funcs), 4)...)
	payload = append(payload, utils.IntToLittleEndian(int(b.Tweak), 4)...)

	if len(flag) > 0 {
		if flag[0] {
			payload = append(payload, 0x01)
		} else {
			payload = append(payload, 0x00)
		}
	} else {
		payload = append(payload, 0x01)
	}

	return network.NewGenericMessage(network.FilterloadCommand, payload)
}

func (b BloomFilter) FilterBytes() []byte {
	return utils.BitFieldToBytes(b.BitField)
}
