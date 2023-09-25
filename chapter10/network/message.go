package network

import (
	"chapter10/utils"
	"crypto/rand"
	"math/big"
	"time"
)

type VersionMessage struct {
	Version          int64
	Services         int64
	Timestamp        int64
	ReceiverServices int64
	ReceiverIP       []byte
	ReceiverPort     int64
	SenderServices   int64
	SenderIP         []byte
	SenderPort       int64
	Nonce            []byte
	UserAgent        []byte
	LastBlock        int64
	Relay            bool
}

func NewVersionMessage(version int64, services int64, timestamp time.Time, receiverServices int64, receiverIP []byte, receiverPort int64, senderServices int64, senderIP []byte, senderPort int64, nonce []byte, userAgent []byte, lastBlock int64, relay bool) (*VersionMessage, error) {
	if version == 0 {
		version = 70015
	}

	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	if receiverIP == nil {
		receiverIP = []byte{0, 0, 0, 0}
	}

	if receiverPort == 0 {
		receiverPort = 8333
	}

	if senderIP == nil {
		senderIP = []byte{0, 0, 0, 0}
	}

	if senderPort == 0 {
		senderPort = 8333
	}

	if nonce == nil {
		temp, err := rand.Int(rand.Reader, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(64), nil))
		if err != nil {
			return nil, err
		}

		nonce = utils.IntToLittleEndian(int(temp.Int64()), 8)
	}

	return &VersionMessage{
		Version:          version,
		Services:         services,
		Timestamp:        timestamp.Unix(),
		ReceiverServices: receiverServices,
		ReceiverIP:       receiverIP,
		ReceiverPort:     receiverPort,
		SenderServices:   senderServices,
		SenderIP:         senderIP,
		SenderPort:       senderPort,
		Nonce:            nonce,
		UserAgent:        userAgent,
		LastBlock:        lastBlock,
		Relay:            relay,
	}, nil
}
