package bloomfilter

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestAdd(t *testing.T) {
	bf := New(10, 5, 99)
	item := []byte("Hello World")
	bf.Add(item)

	expected, _ := hex.DecodeString("0000000a080000000140")

	if !bytes.EqualFold(bf.FilterBytes(), expected) {
		t.Errorf("expected %x, got %x", expected, bf.FilterBytes())
	}

	item = []byte("Goodbye!")
	bf.Add(item)

	expected, _ = hex.DecodeString("4000600a080000010940")

	if !bytes.EqualFold(bf.FilterBytes(), expected) {
		t.Errorf("expected %x, got %x", expected, bf.FilterBytes())
	}
}

func TestFilterload(t *testing.T) {
	bf := New(10, 5, 99)
	item := []byte("Hello World")
	bf.Add(item)
	item = []byte("Goodbye!")
	bf.Add(item)

	expected, _ := hex.DecodeString("0a4000600a080000010940050000006300000001")

	actual, _ := bf.Filterload().Serialize()

	if !bytes.EqualFold(actual, expected) {
		t.Errorf("expected %x, got %x", expected, actual)
	}
}
