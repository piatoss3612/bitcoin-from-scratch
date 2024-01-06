package network

import (
	"bytes"
	"chapter10/block"
	"chapter10/utils"
	"crypto/rand"
	"math/big"
	"time"
)

type Message interface {
	Command() Command
	Serialize() ([]byte, error)
}

type VersionMessage struct {
	Version          int32  // 4 bytes
	Services         int64  // 8 bytes
	Timestamp        int64  // 8 bytes
	ReceiverServices int64  // 8 bytes
	ReceiverIP       []byte // 16 bytes (IPv4)
	ReceiverPort     int16  // 2 bytes
	SenderServices   int64  // 8 bytes
	SenderIP         []byte // 16 bytes (IPv4)
	SenderPort       int16  // 2 bytes
	Nonce            []byte // 8 bytes
	UserAgent        []byte // variable
	LastestBlock     int32  // 4 bytes
	Relay            bool   // 1 byte
}

func DefaultVersionMessage(network ...NetworkType) *VersionMessage {
	msg := &VersionMessage{
		Version:          70015,
		Services:         0,
		Timestamp:        time.Now().Unix(),
		ReceiverServices: 0,
		ReceiverIP:       []byte{0x00, 0x00, 0x00, 0x00},
		ReceiverPort:     8333,
		SenderServices:   0,
		SenderIP:         []byte{0x00, 0x00, 0x00, 0x00},
		SenderPort:       8333,
		Nonce:            []byte{0, 0, 0, 0, 0, 0, 0, 0},
		UserAgent:        []byte("/Satoshi:22.0.0/"),
		LastestBlock:     0,
		Relay:            false,
	}

	if len(network) > 0 {
		switch network[0] {
		case TestNet:
			msg.SenderIP = []byte{0x7F, 0x00, 0x00, 0x01}
			msg.ReceiverPort = 18333
			msg.SenderPort = 18333
		case SimNet:
			msg.ReceiverIP = []byte{0x7F, 0x00, 0x00, 0x01}
			msg.SenderIP = []byte{0x7F, 0x00, 0x00, 0x01}
			msg.ReceiverPort = 18555
			msg.SenderPort = 18555
		}
	}

	return msg
}

func NewVersionMessage(version int32, services int64, timestamp time.Time,
	receiverServices int64, receiverIP []byte, receiverPort int16,
	senderServices int64, senderIP []byte, senderPort int16,
	nonce []byte, userAgent []byte, lastBlock int32, relay bool) (*VersionMessage, error) {
	msg := DefaultVersionMessage()

	if version != 0 {
		msg.Version = version
	}

	if services != 0 {
		msg.Services = services
	}

	if !timestamp.IsZero() {
		msg.Timestamp = timestamp.Unix()
	}

	if receiverServices != 0 {
		msg.ReceiverServices = receiverServices
	}

	if receiverIP != nil {
		msg.ReceiverIP = receiverIP
	}

	if receiverPort != 0 {
		msg.ReceiverPort = receiverPort
	}

	if senderServices != 0 {
		msg.SenderServices = senderServices
	}

	if senderIP != nil {
		msg.SenderIP = senderIP
	}

	if senderPort != 0 {
		msg.SenderPort = senderPort
	}

	if nonce == nil {
		temp, err := rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(64), nil))
		if err != nil {
			return nil, err
		}

		msg.Nonce = utils.IntToLittleEndian(int(temp.Int64()), 8)
	}

	if userAgent != nil {
		msg.UserAgent = userAgent
	}

	if lastBlock != 0 {
		msg.LastestBlock = lastBlock
	}

	msg.Relay = relay

	return msg, nil
}

func (vm VersionMessage) Command() Command {
	return VersionCommand
}

func (vm VersionMessage) Serialize() ([]byte, error) {
	result := utils.IntToLittleEndian(int(vm.Version), 4)
	result = append(result, utils.IntToLittleEndian(int(vm.Services), 8)...)
	result = append(result, utils.IntToLittleEndian(int(vm.Timestamp), 8)...)
	result = append(result, utils.IntToLittleEndian(int(vm.ReceiverServices), 8)...)
	result = append(result, append(append(bytes.Repeat([]byte{0x00}, 10), []byte{0xff, 0xff}...), vm.ReceiverIP...)...)
	result = append(result, utils.IntToBytes(int(vm.ReceiverPort), 2)...)
	result = append(result, utils.IntToLittleEndian(int(vm.SenderServices), 8)...)
	result = append(result, append(append(bytes.Repeat([]byte{0x00}, 10), []byte{0xff, 0xff}...), vm.SenderIP...)...)
	result = append(result, utils.IntToBytes(int(vm.SenderPort), 2)...)
	result = append(result, vm.Nonce...)
	result = append(result, utils.EncodeVarint(len(vm.UserAgent))...)
	result = append(result, vm.UserAgent...)
	result = append(result, utils.IntToLittleEndian(int(vm.LastestBlock), 4)...)

	if vm.Relay {
		result = append(result, 0x01)
	} else {
		result = append(result, 0x00)
	}

	return result, nil
}

type VerAckMessage struct{}

func NewVerAckMessage() *VerAckMessage {
	return &VerAckMessage{}
}

func (vam VerAckMessage) Command() Command {
	return VerAckCommand
}

func (vam VerAckMessage) Serialize() ([]byte, error) {
	return []byte{}, nil
}

type GetHeadersMessage struct {
	Version        int32 // 4 bytes
	NumberOfHashes int64 // variable
	StartBlock     []byte
	EndBlock       []byte
}

func DefaultGetHeadersMessage() *GetHeadersMessage {
	msg := &GetHeadersMessage{
		Version:        70015,
		NumberOfHashes: 1,
		EndBlock:       bytes.Repeat([]byte{0x00}, 32),
		StartBlock:     bytes.Repeat([]byte{0x00}, 32),
	}

	return msg
}

func NewGetHeadersMessage(version int32, numberOfHashes int64, startBlock []byte, endBlock []byte) (*GetHeadersMessage, error) {
	msg := DefaultGetHeadersMessage()

	if version != 0 {
		msg.Version = version
	}

	if numberOfHashes > 0 {
		msg.NumberOfHashes = numberOfHashes
	}

	if startBlock == nil {
		return nil, ErrInvalidStartBlockHash
	}

	if len(startBlock) != 32 {
		return nil, ErrInvalidStartBlockHash
	}

	msg.StartBlock = startBlock

	if endBlock != nil {
		if len(endBlock) != 32 {
			return nil, ErrInvalidEndBlockHash
		}

		msg.EndBlock = endBlock
	}

	return msg, nil
}

func (ghm GetHeadersMessage) Command() Command {
	return GetHeadersCommand
}

func (ghm GetHeadersMessage) Serialize() ([]byte, error) {
	result := utils.IntToLittleEndian(int(ghm.Version), 4)
	result = append(result, utils.EncodeVarint(int(ghm.NumberOfHashes))...)
	result = append(result, ghm.StartBlock...)
	result = append(result, ghm.EndBlock...)

	return result, nil
}

type HeadersMessage struct {
	NumberOfHeaders int64
	Headers         []*block.Block
}

func (hm HeadersMessage) Command() Command {
	return HeadersCommand
}

func (hm HeadersMessage) Serialize() ([]byte, error) {
	return nil, nil
}

type PingMessage struct {
	Nonce []byte
}

func NewPingMessage(nonce []byte) *PingMessage {
	return &PingMessage{
		Nonce: nonce,
	}
}

func (pm PingMessage) Command() Command {
	return PingCommand
}

func (pm PingMessage) Serialize() ([]byte, error) {
	return pm.Nonce, nil
}

type PongMessage struct {
	Nonce []byte
}

func NewPongMessage(nonce []byte) *PongMessage {
	return &PongMessage{
		Nonce: nonce,
	}
}

func (pm PongMessage) Command() Command {
	return PongCommand
}

func (pm PongMessage) Serialize() ([]byte, error) {
	return pm.Nonce, nil
}
