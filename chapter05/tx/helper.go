package tx

import "chapter05/utils"

func ParseTx(b []byte) (*Tx, error) {
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

	return NewTx(version, inputs), nil
}

func ParseTxIn(b []byte) (*TxIn, int) {
	prevTx := utils.ReverseBytes(b[:32])
	prevIndex := utils.LittleEndianToInt(b[32:36])
	scriptSig := b[36:] // TODO: parse scriptSig
	seqNo := utils.LittleEndianToInt(b[len(b)-4:])

	return NewTxIn(prevIndex, string(prevTx), string(scriptSig), seqNo), len(b)
}
