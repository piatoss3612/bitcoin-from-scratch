package script

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestParse(t *testing.T) {
	rawScriptPubkey, _ := hex.DecodeString("6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")

	scriptPubkey, _, err := Parse(rawScriptPubkey)
	if err != nil {
		t.Fatal(err)
	}

	expected1, _ := hex.DecodeString("304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601")

	if !bytes.Equal(scriptPubkey.Cmds[0].Elem, expected1) {
		t.Errorf("expected %x, got %x", expected1, scriptPubkey.Cmds[0].Elem)
	}

	expected2, _ := hex.DecodeString("035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")

	if !bytes.Equal(scriptPubkey.Cmds[1].Elem, expected2) {
		t.Errorf("expected %x, got %x", expected2, scriptPubkey.Cmds[1].Elem)
	}
}

func TestSerialize(t *testing.T) {
	rawScriptPubkey, _ := hex.DecodeString("6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")

	scriptPubkey, _, err := Parse(rawScriptPubkey)
	if err != nil {
		t.Fatal(err)
	}

	b, err := scriptPubkey.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, rawScriptPubkey) {
		t.Errorf("expected %x, got %x", rawScriptPubkey, b)
	}
}
