package utils

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func TestLitleEndianToInt(t *testing.T) {
	tests := []struct {
		caseName     string
		littleEndian string
		expected     int
	}{
		{"1", "99c3980000000000", 10011545},
		{"2", "a135ef0100000000", 32454049},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			h, _ := hex.DecodeString(test.littleEndian)
			actual := LittleEndianToInt(h)
			if actual != test.expected {
				t.Errorf("LittleEndianToInt(%s) = %d, expected = %d", test.littleEndian, actual, test.expected)
			}
		})
	}
}

func TestIntToLittleEndian(t *testing.T) {
	tests := []struct {
		caseName string
		n        int
		length   int
		expected string
	}{
		{"1", 1, 4, "01000000"},
		{"2", 10011545, 8, "99c3980000000000"},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			actual := IntToLittleEndian(test.n, test.length)
			if hex.EncodeToString(actual) != test.expected {
				t.Errorf("IntToLittleEndian(%d) = %s, expected = %s", test.n, hex.EncodeToString(actual), test.expected)
			}
		})
	}
}

func TestBase58(t *testing.T) {
	addr := "mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf"
	h160, _ := DecodeBase58(addr)

	expected, _ := hex.DecodeString("507b27411ccf7f16f10297de6cef3f291623eddf")

	if !bytes.EqualFold(h160, expected) {
		t.Fatalf("DecodeBase58(%s) = %x, expected = %x", addr, h160, expected)
	}

	enc := EncodeBase58Checksum(append([]byte{0x6f}, h160...))

	if !strings.EqualFold(enc, addr) {
		t.Fatalf("EncodeBase58Checksum(%x) = %s, expected = %s", h160, enc, addr)
	}
}

func TestBitFieldToBytes(t *testing.T) {
	bitField := []byte{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	want := "4000600a080000010940"
	actual := hex.EncodeToString(BitFieldToBytes(bitField))
	if !strings.EqualFold(actual, want) {
		t.Fatalf("BitFieldToBytes(%x) = %s, expected = %s", bitField, actual, want)
	}
}
