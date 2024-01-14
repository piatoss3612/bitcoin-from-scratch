package script

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestOpHash160(t *testing.T) {
	stack := []any{[]byte("hello world")}
	if !OpHash160(&stack) {
		t.Error("opHash160 failed")
	}

	elem, ok := stack[0].([]byte)
	if !ok {
		t.Error("opHash160 failed")
	}

	hexElem := hex.EncodeToString(elem)

	if !strings.EqualFold(hexElem, "d7d5ee7824ff93f94c3055af9382c86c68b5ca92") {
		t.Error("opHash160 failed")
	}
}

func TestOpCheckSig(t *testing.T) {
	z, _ := hex.DecodeString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	sec, _ := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	sig, _ := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab601")

	stack := []any{sig, sec}

	if !OpCheckSig(&stack, z) {
		t.Error("opCheckSig failed")
	}

	elem, ok := stack[0].([]byte)
	if !ok {
		t.Error("opCheckSig failed")
	}

	n := DecodeNum(elem)

	if n != 1 {
		t.Error("opCheckSig failed")
	}
}

func TestOpCheckMultiSig(t *testing.T) {
	z, _ := hex.DecodeString("e71bfa115715d6fd33796948126f40a8cdd39f187e4afb03896795189fe1423c")
	sig1, _ := hex.DecodeString("3045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a8993701")
	sig2, _ := hex.DecodeString("3045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201")
	sec1, _ := hex.DecodeString("022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb70")
	sec2, _ := hex.DecodeString("03b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb71")

	stack := []any{[]byte(""), sig1, sig2, []byte{0x02}, sec1, sec2, []byte{0x02}}

	if !OpCheckMultiSig(&stack, z) {
		t.Error("opCheckMultiSig failed")
	}

	elem, ok := stack[0].([]byte)
	if !ok {
		t.Error("opCheckMultiSig failed")
	}

	n := DecodeNum(elem)

	if n != 1 {
		t.Error("opCheckMultiSig failed")
	}
}