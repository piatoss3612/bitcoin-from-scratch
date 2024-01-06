package network

import (
	"bytes"
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
		ReceiverIP:       []byte{0x7F, 0x00, 0x00, 0x01},
		ReceiverPort:     8333,
		SenderServices:   0,
		SenderIP:         []byte{0x7F, 0x00, 0x00, 0x01},
		SenderPort:       8333,
		Nonce:            []byte{0, 0, 0, 0, 0, 0, 0, 0},
		UserAgent:        []byte("/Satoshi:0.0.1/"),
		LastestBlock:     0,
		Relay:            false,
	}

	if len(network) > 0 {
		switch network[0] {
		case TestNet:
			msg.ReceiverPort = 18333
			msg.SenderPort = 18333
		case SimNet:
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

type VerAckMessage struct {
}

func NewVerAckMessage() *VerAckMessage {
	return &VerAckMessage{}
}

func (vam VerAckMessage) Command() Command {
	return VerAckCommand
}

func (vam VerAckMessage) Serialize() []byte {
	return []byte{}
}
