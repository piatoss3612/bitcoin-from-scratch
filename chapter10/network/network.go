package network

import "bytes"

var (
	NetworkMagic     = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	TestNetworkMagic = []byte{0x0b, 0x11, 0x09, 0x07}
	SimNetMagic      = []byte{0x0b, 0x11, 0x09, 0x05}
)

type NetworkType uint8

const (
	MainNet NetworkType = iota
	TestNet
	SimNet
)

func (nt NetworkType) String() string {
	switch nt {
	case MainNet:
		return "MainNet"
	case TestNet:
		return "TestNet"
	case SimNet:
		return "SimNet"
	default:
		return "Unknown"
	}
}

func IsNetworkMagicValid(magic []byte) bool {
	return bytes.Equal(magic, NetworkMagic) || bytes.Equal(magic, TestNetworkMagic) || bytes.Equal(magic, SimNetMagic)
}
