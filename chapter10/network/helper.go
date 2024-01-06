package network

import (
	"bytes"
	"chapter10/block"
	"chapter10/utils"
	"fmt"
)

var (
	ErrInvalidNetworkMagic   = fmt.Errorf("invalid network magic")
	ErrInvalidPayload        = fmt.Errorf("invalid payload")
	ErrInvalidNetwork        = fmt.Errorf("invalid network")
	ErrInvalidCommand        = fmt.Errorf("invalid command")
	ErrInvalidStartBlockHash = fmt.Errorf("invalid start block hash")
	ErrInvalidEndBlockHash   = fmt.Errorf("invalid end block hash")
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

func ParseCommand(b []byte) (Command, error) {
	cmd := Command(bytes.Trim(b, "\x00"))
	if !cmd.IsValid() {
		return nil, ErrInvalidCommand
	}

	return cmd, nil
}

func ParseHeadersMessage(b []byte) (*HeadersMessage, error) {
	buf := bytes.NewBuffer(b)

	numOfHeaders, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	headers := make([]*block.Block, numOfHeaders)

	for i := 0; i < int(numOfHeaders); i++ {
		header, err := block.Parse(buf.Bytes())
		if err != nil {
			return nil, err
		}

		headers[i] = header

		numOfTxns, _ := utils.ReadVarint(buf.Bytes())
		if numOfTxns > 0 {
			return nil, fmt.Errorf("block %d has %d transactions", i, numOfTxns)
		}

		buf.Next(81)
	}

	return &HeadersMessage{
		NumberOfHeaders: int64(numOfHeaders),
		Headers:         headers,
	}, nil
}

func ParsePingMessage(b []byte) (*PingMessage, error) {
	buf := bytes.NewBuffer(b)

	nonce := buf.Next(8)

	return &PingMessage{
		Nonce: nonce,
	}, nil
}

func ParsePongMessage(b []byte) (*PongMessage, error) {
	buf := bytes.NewBuffer(b)

	nonce := buf.Next(8)

	return &PongMessage{
		Nonce: nonce,
	}, nil
}
