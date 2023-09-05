package tx

import "chapter05/utils"

func ParseTx(b []byte, testnet bool) (*Tx, error) {
	version := utils.LittleEndianToInt(b[:4])
	b = b[4:]

	numInputs, read := utils.ReadVarint(b)
	b = b[read:]

	inputs := []*TxIn{}

	for i := 0; i < numInputs; i++ {
		txIn, read := ParseTxIn(b)
		inputs = append(inputs, txIn)
		b = b[read:]
	}

	numOutputs, read := utils.ReadVarint(b)
	b = b[read:]

	outputs := []*TxOut{}

	for i := 0; i < numOutputs; i++ {
		txOut, read := ParseTxOut(b)
		outputs = append(outputs, txOut)
		b = b[read:]
	}

	lockTime := utils.LittleEndianToInt(b[:4])

	return NewTx(version, inputs, outputs, lockTime, testnet), nil
}

func ParseTxIn(b []byte) (*TxIn, int) {
	prevTx := utils.ReverseBytes(b[:32])
	prevIndex := utils.LittleEndianToInt(b[32:36])
	scriptSig := b[36:] // TODO: parse scriptSig
	seqNo := utils.LittleEndianToInt(b[len(b)-4:])

	return NewTxIn(prevIndex, string(prevTx), string(scriptSig), seqNo), len(b)
}

func ParseTxOut(b []byte) (*TxOut, int) {
	value := utils.LittleEndianToInt(b[:8])
	scriptPubKey := b[8:] // TODO: parse scriptPubKey

	return NewTxOut(value, string(scriptPubKey)), len(b)
}
