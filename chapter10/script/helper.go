package script

import (
	"bytes"
	"chapter10/utils"
	"errors"
)

func Parse(b []byte) (*Script, int, error) {
	length, read := utils.ReadVarint(b) // 가변 정수로된 스크립트의 전체 길이를 읽음

	buf := bytes.NewBuffer(b[read:]) // 가변 정수를 제외한 나머지 스크립트를 버퍼에 저장

	var cmds []any // 스크립트 명령어를 저장할 슬라이스
	var count int  // 읽은 바이트 수

	// 읽어들인 바이트 수가 전체 길이보다 작은 동안 반복
	for count < length {
		current := buf.Next(1) // 버퍼에서 1바이트를 읽음
		count += 1             // 읽은 바이트 수를 1 증가

		currentByte := current[0]

		if currentByte >= 1 && currentByte <= 75 { // 바이트 값이 1에서 75 사이인 경우: 해당 길이만큼 데이터를 읽어 원소로 추가
			dataLength := int(currentByte)
			data := buf.Next(dataLength)
			count += dataLength
			cmds = append(cmds, data)
		} else if currentByte == 76 { // 바이트 값이 76인 경우: OP_PUSHDATA1에 해당하므로 다음 한 바이트를 더 읽어 해당 길이만큼 데이터를 읽어 원소로 추가
			dataLength := utils.LittleEndianToInt(buf.Next(1))
			data := buf.Next(dataLength)
			cmds = append(cmds, data)
			count += dataLength + 1
		} else if currentByte == 77 { // 바이트 값이 77인 경우: OP_PUSHDATA2에 해당하므로 다음 두 바이트를 더 읽어 해당 길이만큼 데이터를 읽어 원소로 추가
			dataLength := utils.LittleEndianToInt(buf.Next(2))
			data := buf.Next(dataLength)
			cmds = append(cmds, data)
			count += dataLength + 2
		} else { // 그 외의 경우: 해당 바이트 값을 연산자로 간주하여 추가
			opCode := int(currentByte)
			cmds = append(cmds, opCode)
		}
	}

	// 읽은 바이트 수와 전체 길이가 일치하지 않는 경우 에러 반환
	if count != length {
		return nil, 0, errors.New("parse error: length mismatch")
	}

	// 스크립트와 읽은 바이트 수, 에러를 반환
	return New(cmds...), read + length, nil
}

func NewP2PKHScript(h160 []byte) *Script {
	return New(
		0x76, // OP_DUP
		0xa9, // OP_HASH160
		h160, // 20바이트의 데이터
		0x88, // OP_EQUALVERIFY
		0xac, // OP_CHECKSIG
	)
}

func IsP2pkhScriptPubkey(cmds []any) bool {
	if len(cmds) != 5 {
		return false
	}

	opCodeDup, ok1 := cmds[0].(int)
	opCodeHash160, ok2 := cmds[1].(int)
	h160, ok3 := cmds[2].([]byte)
	opCodeEqualVerify, ok4 := cmds[3].(int)
	opCodeCheckSig, ok5 := cmds[4].(int)

	return ok1 && ok2 && ok3 && ok4 && ok5 &&
		opCodeDup == 0x76 && opCodeHash160 == 0xa9 && len(h160) == 20 && opCodeEqualVerify == 0x88 && opCodeCheckSig == 0xac
}

func IsP2shScriptPubkey(cmds []any) bool {
	if len(cmds) != 3 {
		return false
	}

	opCodeHash160, ok1 := cmds[0].(int)
	h160, ok2 := cmds[1].([]byte)
	opCodeEqual, ok3 := cmds[2].(int)

	return ok1 && ok2 && ok3 &&
		opCodeHash160 == 0xa9 && len(h160) == 20 && opCodeEqual == 0x87
}
