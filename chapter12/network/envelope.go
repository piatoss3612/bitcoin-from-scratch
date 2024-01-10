package network

import (
	"chapter12/utils"
	"encoding/hex"
	"fmt"
)

type NetworkEnvelope struct {
	Magic   []byte  // 4 bytes
	Command Command // 12 bytes
	Payload []byte  // variable
}

func NewEnvelope(command, payload []byte, network ...NetworkType) (*NetworkEnvelope, error) {
	ne := &NetworkEnvelope{
		Magic:   NetworkMagic,
		Command: command,
		Payload: payload,
	}

	if len(network) > 0 {
		switch network[0] {
		case TestNet:
			ne.Magic = TestNetworkMagic
		case SimNet:
			ne.Magic = SimNetMagic
		}
	}

	return ne, nil
}

func (ne NetworkEnvelope) String() string {
	return fmt.Sprintf("%s %s", ne.Command, hex.EncodeToString(ne.Payload))
}

func (ne NetworkEnvelope) Serialize() ([]byte, error) {
	result := ne.Magic[:]

	result = append(result, ne.Command.Serialize()...)
	result = append(result, utils.IntToLittleEndian(len(ne.Payload), 4)...)
	result = append(result, utils.Hash256(ne.Payload)[:4]...)
	result = append(result, ne.Payload...)

	return result, nil
}
