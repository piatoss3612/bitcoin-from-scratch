package network

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestParseEnvelope(t *testing.T) {
	rawMsg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	envelope, err := ParseNetworkEnvelope(rawMsg)
	if err != nil {
		t.Fatalf("error parsing envelope: %v", err)
	}

	if envelope.Command.String() != "verack" {
		t.Fatalf("expected command verack, got %s", envelope.Command)
	}

	if len(envelope.Payload) != 0 {
		t.Fatalf("expected empty payload, got %v", envelope.Payload)
	}

	rawMsg, _ = hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	envelope, err = ParseNetworkEnvelope(rawMsg)
	if err != nil {
		t.Fatalf("error parsing envelope: %v", err)
	}

	if envelope.Command.String() != "version" {
		t.Fatalf("expected command version, got %s", envelope.Command)
	}

	if !bytes.Equal(envelope.Payload, rawMsg[24:]) {
		t.Fatalf("expected payload %v, got %v", rawMsg[24:], envelope.Payload)
	}
}

func TestSerializeEnvelope(t *testing.T) {
	rawMsg, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")
	envelope, _ := ParseNetworkEnvelope(rawMsg)
	s, err := envelope.Serialize()
	if err != nil {
		t.Fatalf("error serializing envelope: %v", err)
	}

	if !bytes.Equal(s, rawMsg) {
		t.Fatalf("expected %v, got %v", rawMsg, s)
	}

	rawMsg, _ = hex.DecodeString("f9beb4d976657273696f6e0000000000650000005f1a69d2721101000100000000000000bc8f5e5400000000010000000000000000000000000000000000ffffc61b6409208d010000000000000000000000000000000000ffffcb0071c0208d128035cbc97953f80f2f5361746f7368693a302e392e332fcf05050001")
	envelope, _ = ParseNetworkEnvelope(rawMsg)
	s, err = envelope.Serialize()
	if err != nil {
		t.Fatalf("error serializing envelope: %v", err)
	}

	if !bytes.Equal(s, rawMsg) {
		t.Fatalf("expected %v, got %v", rawMsg, s)
	}
}
