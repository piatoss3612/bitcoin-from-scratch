package network

import (
	"bytes"
	"chapter10/utils"
	"fmt"
)

var (
	ErrInvalidNetworkMagic = fmt.Errorf("invalid network magic")
	ErrInvalidPayload      = fmt.Errorf("invalid payload")
)

func ParseNetworkEnvelope(b []byte) (*NetworkEnvelope, error) {
	buf := bytes.NewBuffer(b)

	magic := buf.Next(4)

	if !IsNetworkMagicValid(magic) {
		return nil, ErrInvalidNetworkMagic
	}

	command := buf.Next(12)

	payloadLength := utils.LittleEndianToInt(buf.Next(4))
	payloadChecksum := buf.Next(4)

	payload := buf.Next(payloadLength)

	if !bytes.Equal(payloadChecksum, utils.Hash256(payload)[:4]) {
		return nil, ErrInvalidPayload
	}

	return &NetworkEnvelope{
		magic:   magic,
		command: command,
		payload: payload,
	}, nil
}

func ParseVerAckMessage(b []byte) *VerAckMessage {
	return NewVerAckMessage()
}
