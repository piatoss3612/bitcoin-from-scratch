package network

import (
	"bytes"
	"chapter13/block"
	"chapter13/utils"
	"fmt"
)

var (
	ErrInvalidNetworkMagic   = fmt.Errorf("invalid network magic")
	ErrInvalidPayload        = fmt.Errorf("invalid payload")
	ErrInvalidNetwork        = fmt.Errorf("invalid network")
	ErrInvalidNetworkMessage = fmt.Errorf("invalid network message")
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

	command := ParseCommand(buf.Next(12))

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

func ParseCommand(b []byte) Command {
	return Command(bytes.Trim(b, "\x00"))
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

		buf.Next(80)

		numOfTxns, n := utils.ReadVarint(buf.Bytes())

		if numOfTxns != 0 {
			return nil, fmt.Errorf("block %d has %d transactions", i, numOfTxns)
		}

		buf.Next(n)
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
