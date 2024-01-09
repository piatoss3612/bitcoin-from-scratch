package network

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSerializeVersionMessage(t *testing.T) {
	v := DefaultVersionMessage()
	v.Timestamp = 0
	v.Nonce = bytes.Repeat([]byte{0x00}, 8)

	expected, _ := hex.DecodeString("7f11010000000000000000000000000000000000000000000000000000000000000000000000ffff00000000208d000000000000000000000000000000000000ffff00000000208d0000000000000000182f70726f6772616d6d696e67626974636f696e3a302e312f0000000000")

	actual, _ := v.Serialize()

	if !bytes.EqualFold(actual, expected) {
		t.Errorf("expected %x, got %x", expected, actual)
	}
}

func TestSerializeGetHeadersMessage(t *testing.T) {
	rawBlock, _ := hex.DecodeString("0000000000000000001237f46acddf58578a37e213d2a6edc4884a2fcad05ba3")

	getheaders := DefaultGetHeadersMessage()
	getheaders.StartBlock = rawBlock

	expected, _ := hex.DecodeString("7f11010001a35bd0ca2f4a88c4eda6d213e2378a5758dfcd6af437120000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	actual, _ := getheaders.Serialize()

	if !bytes.EqualFold(actual, expected) {
		t.Errorf("expected %x, got %x", expected, actual)
	}
}

func TestSerializeGetDataMessage(t *testing.T) {
	expected, _ := hex.DecodeString("020300000030eb2540c41025690160a1014c577061596e32e426b712c7ca00000000000000030000001049847939585b0652fba793661c361223446b6fc41089b8be00000000000000")

	getdata := NewGetDataMessage()

	block1, _ := hex.DecodeString("00000000000000cac712b726e4326e596170574c01a16001692510c44025eb30")
	block2, _ := hex.DecodeString("00000000000000beb88910c46f6b442312361c6693a7fb52065b583979844910")

	getdata.AddData(3, block1)
	getdata.AddData(3, block2)

	actual, _ := getdata.Serialize()

	if !bytes.EqualFold(actual, expected) {
		t.Errorf("expected %x, got %x", expected, actual)
	}
}