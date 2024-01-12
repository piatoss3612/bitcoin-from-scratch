package network

import "testing"

func TestHandshake(t *testing.T) {
	node, err := NewSimpleNode("71.13.92.62", 18333, TestNet, false)
	if err != nil {
		t.Fatalf("error creating node: %v", err)
	}

	res, err := node.HandShake()
	if err != nil {
		t.Fatalf("error handshaking: %v", err)
	}

	if ok := <-res; !ok {
		t.Fatalf("expected true, got false")
	}
}
