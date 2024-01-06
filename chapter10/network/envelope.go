package network

import (
	"chapter10/utils"
	"encoding/hex"
	"fmt"
)

type NetworkEnvelope struct {
	magic   []byte // 4 bytes
	command []byte // 12 bytes
	payload []byte // variable
}

func New(command, payload []byte, network ...NetworkType) (*NetworkEnvelope, error) {
	ne := &NetworkEnvelope{
		magic:   NetworkMagic,
		command: command,
		payload: payload,
	}

	if len(network) > 0 {
		switch network[0] {
		case TestNet:
			ne.magic = TestNetworkMagic
		case SimNet:
			ne.magic = SimNetMagic
		}
	}

	return ne, nil
}

func (ne NetworkEnvelope) String() string {
	return fmt.Sprintf("%s %s", ne.command, hex.EncodeToString(ne.payload))
}

func (ne NetworkEnvelope) Serialize() ([]byte, error) {
	result := ne.magic[:]

	command := make([]byte, 12)
	copy(command, ne.command)
	result = append(result, command...)

	result = append(result, utils.IntToLittleEndian(len(ne.payload), 4)...)
	result = append(result, utils.Hash256(ne.payload)[:4]...)
	result = append(result, ne.payload...)

	return result, nil
}
