package network

import (
	"chapter10/utils"
	"encoding/hex"
	"fmt"
)

var (
	NetworkMagic     = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	TestNetworkMagic = []byte{0x0b, 0x11, 0x09, 0x07}
)

type NetworkEnvelope struct {
	magic   []byte
	command []byte
	payload []byte
}

func New(command, payload []byte, testnet bool) (*NetworkEnvelope, error) {
	ne := &NetworkEnvelope{
		magic:   NetworkMagic,
		command: command,
		payload: payload,
	}

	if testnet {
		ne.magic = TestNetworkMagic
	}

	return ne, nil
}

func (ne NetworkEnvelope) String() string {
	return fmt.Sprintf("%s %s", ne.command, hex.EncodeToString(ne.payload))
}

func (ne NetworkEnvelope) Serialize() []byte {
	result := ne.magic[:]
	result = append(result, ne.command...)
	result = append(result, utils.IntToLittleEndian(len(ne.payload), 4)...)
	result = append(result, utils.Hash256(ne.payload)[:4]...)
	result = append(result, ne.payload...)

	return result
}
