package script

import (
	"bytes"
	"chapter13/utils"
	"errors"
)

func Parse(b []byte) (*Script, int, error) {
	length, read := utils.ReadVarint(b) // 가변 정수로된 스크립트의 전체 길이를 읽음

	buf := bytes.NewBuffer(b[read:]) // 가변 정수를 제외한 나머지 스크립트를 버퍼에 저장

	var cmds []Command // 스크립트 명령어를 저장할 슬라이스
	var count int      // 읽은 바이트 수

	// 읽어들인 바이트 수가 전체 길이보다 작은 동안 반복
	for count < length {
		current := buf.Next(1) // 버퍼에서 1바이트를 읽음
		count += 1             // 읽은 바이트 수를 1 증가

		currentByte := current[0]

		if currentByte >= 1 && currentByte <= 75 { // 바이트 값이 1에서 75 사이인 경우: 해당 길이만큼 데이터를 읽어 원소로 추가
			n := int(currentByte)                     // 읽어들인 바이트 값을 정수로 변환
			cmds = append(cmds, NewElem(buf.Next(n))) // 해당 길이만큼 데이터를 읽어 원소로 추가
			count += n                                // 읽은 바이트 수를 해당 길이만큼 증가
		} else if currentByte == 76 { // 바이트 값이 76인 경우: OP_PUSHDATA1에 해당하므로 다음 한 바이트를 더 읽어 해당 길이만큼 데이터를 읽어 원소로 추가
			n := utils.LittleEndianToInt(buf.Next(1))
			cmds = append(cmds, NewElem(buf.Next(n)))
			count += n + 1
		} else if currentByte == 77 { // 바이트 값이 77인 경우: OP_PUSHDATA2에 해당하므로 다음 두 바이트를 더 읽어 해당 길이만큼 데이터를 읽어 원소로 추가
			n := utils.LittleEndianToInt(buf.Next(2))
			cmds = append(cmds, NewElem(buf.Next(n)))
			count += n + 2
		} else { // 그 외의 경우: 해당 바이트 값을 연산자로 간주하여 추가
			opCode := int(currentByte)
			cmds = append(cmds, NewOpCode(OpCode(opCode)))
		}
	}

	// 읽은 바이트 수와 전체 길이가 일치하지 않는 경우 에러 반환
	if count != length {
		return nil, 0, errors.New("parse error: length mismatch")
	}

	// 스크립트와 읽은 바이트 수, 에러를 반환
	return New(cmds...), read + length, nil
}

func NewP2pkhScript(h160 []byte) *Script {
	return New(
		NewOpCode(OpCodeDup),         // OP_DUP
		NewOpCode(OpCodeHash160),     // OP_HASH160
		NewElem(h160),                // 20바이트의 데이터
		NewOpCode(OpCodeEqualVerify), // OP_EQUALVERIFY
		NewOpCode(OpCodeCheckSig),    // OP_CHECKSIG
	)
}

func NewP2shScript(h160 []byte) *Script {
	return New(
		NewOpCode(OpCodeHash160), // OP_HASH160
		NewElem(h160),            // 20바이트의 데이터
		NewOpCode(OpCodeEqual),   // OP_EQUAL
	)
}

func NewP2wpkhScript(h160 []byte) *Script {
	return New(
		NewOpCode(OpCode0), // OP_0
		NewElem(h160),      // 20바이트의 데이터
	)
}

func NewP2wshScript(h256 []byte) *Script {
	return New(
		NewOpCode(OpCode0), // OP_0
		NewElem(h256),      // 32바이트의 데이터
	)
}

func IsP2pkhScriptPubkey(cmds []Command) bool {
	return len(cmds) == 5 && cmds[0].Code == OpCodeDup && cmds[1].Code == OpCodeHash160 &&
		len(cmds[2].Elem) == 20 && cmds[3].Code == OpCodeEqualVerify && cmds[4].Code == OpCodeCheckSig
}

func IsP2shScriptPubkey(cmds []Command) bool {
	return len(cmds) == 3 && cmds[0].Code == OpCodeHash160 &&
		len(cmds[1].Elem) == 20 && cmds[2].Code == OpCodeEqual
}

func IsP2wpkhScriptPubkey(cmds []Command) bool {
	return len(cmds) == 2 && cmds[0].Code == OpCode0 && len(cmds[1].Elem) == 20
}

func IsP2wshScriptPubkey(cmds []Command) bool {
	return len(cmds) == 2 && cmds[0].Code == OpCode0 && len(cmds[1].Elem) == 32
}
