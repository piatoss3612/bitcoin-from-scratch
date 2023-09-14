package tx

import (
	"bytes"
	"chapter07/script"
	"chapter07/utils"
	"encoding/hex"
)

const (
	SIGHASH_ALL    = 1
	SIGHASH_NONE   = 2
	SIGHASH_SINGLE = 3
)

// 트랜잭션을 파싱하는 함수
func ParseTx(b []byte, testnet bool) (*Tx, error) {
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

	return NewTx(version, inputs, outputs, lockTime, testnet), nil
}

// TxIn을 파싱하는 함수
func ParseTxIn(b []byte) (*TxIn, int, error) {
	buf := bytes.NewBuffer(b)

	prevTx := utils.ReverseBytes(buf.Next(32))

	prevIndex := utils.LittleEndianToInt(buf.Next(4))

	scriptSig, read, err := script.Parse(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}

	buf.Next(read)

	seqNo := utils.LittleEndianToInt(buf.Next(4))

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
