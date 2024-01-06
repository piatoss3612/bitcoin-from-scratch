package network

import (
	"bytes"
	"chapter10/utils"
	"fmt"
)

var (
	ErrInvalidNetworkMagic = fmt.Errorf("invalid network magic")
	ErrInvalidPayload      = fmt.Errorf("invalid payload")
	ErrInvalidNetwork      = fmt.Errorf("invalid network")
	ErrInvalidCommand      = fmt.Errorf("invalid command")
)

func ParseNetworkEnvelope(b []byte) (*NetworkEnvelope, error) {
	buf := bytes.NewBuffer(b)

	magic := buf.Next(4)

	if !IsNetworkMagicValid(magic) {
		return nil, ErrInvalidNetworkMagic
	}

	command, err := ParseCommand(buf.Next(12))
	if err != nil {
		return nil, err
	}

	payloadLength := utils.LittleEndianToInt(buf.Next(4))
	payloadChecksum := buf.Next(4)

	payload := buf.Next(payloadLength)

	if !bytes.Equal(payloadChecksum, utils.Hash256(payload)[:4]) {
		return nil, ErrInvalidPayload
	}

	return &NetworkEnvelope{
		Magic:   magic,
		Command: command,
		Payload: payload,
	}, nil
}

func ParseVerAckMessage(b []byte) *VerAckMessage {
	return NewVerAckMessage()
}

func ParseCommand(b []byte) (Command, error) {
	cmd := Command(bytes.Trim(b, "\x00"))
	if !cmd.IsValid() {
		return nil, ErrInvalidCommand
	}

	return cmd, nil
}
