package tx

import (
	"bytes"
	"chapter13/script"
	"chapter13/utils"
	"encoding/hex"
	"errors"
)

const (
	SIGHASH_ALL    = 1
	SIGHASH_NONE   = 2
	SIGHASH_SINGLE = 3
)

var (
	ErrInvalidSegwitTx = errors.New("invalid segwit tx")
)

// 트랜잭션을 파싱하는 함수
func ParseTx(b []byte, testnet ...bool) (*Tx, error) {
	buf := bytes.NewBuffer(b)
	buf.Next(4) // skip 4 bytes

	var parseMethod func([]byte, bool) (*Tx, error)

	if buf.Next(1)[0] == 0x00 {
		parseMethod = parseSegwitTx
	} else {
		parseMethod = parseLegacyTx
	}

	if len(testnet) > 0 {
		return parseMethod(b, testnet[0])
	}

	return parseMethod(b, false)
}

// 세그윗 이전의 트랜잭션을 파싱하는 함수
func parseLegacyTx(b []byte, testnet bool) (*Tx, error) {
	buf := bytes.NewBuffer(b)

	version := utils.LittleEndianToInt(buf.Next(4))

	numInputs, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	inputs := []*TxIn{}

	for i := 0; i < numInputs; i++ {
		txIn, read, err := ParseTxIn(buf.Bytes())
		if err != nil {
			return nil, err
		}

		inputs = append(inputs, txIn)
		buf.Next(read)
	}

	numOutputs, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	outputs := []*TxOut{}

	for i := 0; i < numOutputs; i++ {
		txOut, read, err := ParseTxOut(buf.Bytes())
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, txOut)
		buf.Next(read)
	}

	lockTime := utils.LittleEndianToInt(buf.Next(4))

	return NewTx(version, inputs, outputs, lockTime, testnet, false), nil
}

// 세그윗 트랜잭션을 파싱하는 함수
func parseSegwitTx(b []byte, testnet bool) (*Tx, error) {
	buf := bytes.NewBuffer(b)

	version := utils.LittleEndianToInt(buf.Next(4))
	marker := buf.Next(2)

	if marker[0] != 0x00 || marker[1] != 0x01 {
		return nil, ErrInvalidSegwitTx
	}

	numInputs, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	inputs := []*TxIn{}

	for i := 0; i < numInputs; i++ {
		txIn, read, err := ParseTxIn(buf.Bytes())
		if err != nil {
			return nil, err
		}

		inputs = append(inputs, txIn)
		buf.Next(read)
	}

	numOutputs, read := utils.ReadVarint(buf.Bytes())
	buf.Next(read)

	outputs := []*TxOut{}

	for i := 0; i < numOutputs; i++ {
		txOut, read, err := ParseTxOut(buf.Bytes())
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, txOut)
		buf.Next(read)
	}

	for _, txin := range inputs {
		numItems, read := utils.ReadVarint(buf.Bytes())
		buf.Next(read)

		items := [][]byte{}

		for i := 0; i < numItems; i++ {
			itemLen, read := utils.ReadVarint(buf.Bytes())
			buf.Next(read)

			if itemLen == 0 {
				items = append(items, []byte{})
			} else {
				items = append(items, buf.Next(int(itemLen)))
			}
		}

		txin.Witness = items
	}

	lockTime := utils.LittleEndianToInt(buf.Next(4))

	return NewTx(version, inputs, outputs, lockTime, testnet, true), nil
}

// TxIn을 파싱하는 함수
func ParseTxIn(b []byte) (*TxIn, int, error) {
	buf := bytes.NewBuffer(b)

	prevTx := utils.ReverseBytes(buf.Next(32)) // 이전 트랜잭션의 해시값 (32바이트, 리틀엔디언)

	prevIndex := utils.LittleEndianToInt(buf.Next(4)) // 이전 트랜잭션의 인덱스 (4바이트, 리틀엔디언)

	scriptSig, read, err := script.Parse(buf.Bytes()) // 스크립트 (가변)
	if err != nil {
		return nil, 0, err
	}

	buf.Next(read)

	seqNo := utils.LittleEndianToInt(buf.Next(4)) // 시퀀스 번호 (4바이트, 리틀엔디언)

	return NewTxIn(hex.EncodeToString(prevTx), prevIndex, scriptSig, seqNo), 40 + read, nil
}

// TxOut을 파싱하는 함수
func ParseTxOut(b []byte) (*TxOut, int, error) {
	buf := bytes.NewBuffer(b)

	value := utils.LittleEndianToInt(buf.Next(8))

	scriptPubKey, read, err := script.Parse(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}

	return NewTxOut(value, scriptPubKey), 8 + read, nil
}
