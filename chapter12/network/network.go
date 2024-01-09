package network

import "bytes"

var (
	NetworkMagic     = []byte{0xf9, 0xbe, 0xb4, 0xd9}
	TestNetworkMagic = []byte{0x0b, 0x11, 0x09, 0x07}
	RegTestMagic     = []byte{0xfa, 0xbf, 0xb5, 0xda}
	SimNetMagic      = []byte{0x16, 0x1c, 0x14, 0x12}
)

type NetworkType uint8

const (
	MainNet NetworkType = iota
	TestNet
	RegTest
	SimNet
)

const (
	DefaultMainNetPort = 8333
	DefaultTestNetPort = 18333
	DefaultRegTestPort = 18444
	DefaultSimNetPort  = 18555
)

func (nt NetworkType) String() string {
	switch nt {
	case MainNet:
		return "MainNet"
	case TestNet:
		return "TestNet"
	case RegTest:
		return "RegTest"
	case SimNet:
		return "SimNet"
	default:
		return "Unknown"
	}
}

func (nt NetworkType) Magic() []byte {
	switch nt {
	case MainNet:
		return NetworkMagic
	case TestNet:
		return TestNetworkMagic
	case RegTest:
		return RegTestMagic
	case SimNet:
		return SimNetMagic
	default:
		return nil
	}
}

func IsNetworkMagicValid(magic []byte) bool {
	return bytes.Equal(magic, NetworkMagic) || bytes.Equal(magic, TestNetworkMagic) || bytes.Equal(magic, SimNetMagic) || bytes.Equal(magic, RegTestMagic)
}

type Command []byte

var (
	VersionCommand    Command = []byte("version")
	VerAckCommand     Command = []byte("verack")
	PingCommand       Command = []byte("ping")
	PongCommand       Command = []byte("pong")
	GetHeadersCommand Command = []byte("getheaders")
	HeadersCommand    Command = []byte("headers")
	FilterloadCommand Command = []byte("filterload")
)

func (c Command) String() string {
	return string(bytes.Trim(c, "\x00"))
}

func (c Command) Serialize() []byte {
	b := make([]byte, 12)
	copy(b, c)

	return b
}

func (c Command) Compare(other Command) bool {
	return bytes.Equal(c, other)
}
