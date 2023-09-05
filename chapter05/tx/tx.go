package tx

import (
	"chapter05/utils"
	"encoding/hex"
	"fmt"
)

type Tx struct {
	version  int      // 트랜잭션의 버전
	inputs   []*TxIn  // 트랜잭션의 입력 목록
	outputs  []*TxOut // 트랜잭션의 출력 목록
	lockTime int      // 트랜잭션의 유효 시점
	testnet  bool     // 테스트넷인지 여부
}

func NewTx(version int, inputs []*TxIn, outputs []*TxOut, lockTime int, testnet bool) *Tx {
	tx := &Tx{
		version:  version,
		inputs:   inputs,
		outputs:  outputs,
		lockTime: lockTime,
		testnet:  testnet,
	}

	return tx
}

func (t Tx) String() string {
	return fmt.Sprintf("tx: %s\nversion: %d\ninputs: %s\noutputs: %s\nlocktime: %d",
		t.ID(), t.version, t.inputs, t.outputs, t.lockTime)
}

func (t Tx) ID() string {
	return hex.EncodeToString(t.Hash())
}

func (t Tx) Hash() []byte {
	return utils.ReverseBytes(utils.Hash256(t.Serialize()))
}

func (t Tx) Serialize() []byte {
	result := utils.IntToLittleEndian(t.version, 4)                    // 버전
	result = append(result, t.serializeInputs()...)                    // 입력 목록
	result = append(result, t.serializeOutputs()...)                   // 출력 목록
	result = append(result, utils.IntToLittleEndian(t.lockTime, 4)...) // 유효 시점
	return result
}

func (t Tx) serializeInputs() []byte {
	result := utils.EncodeVarint(len(t.inputs)) // 입력 개수

	// 입력 개수만큼 반복하면서 각 입력을 직렬화한 결과를 result에 추가
	for _, input := range t.inputs {
		result = append(result, input.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

func (t Tx) serializeOutputs() []byte {
	result := utils.EncodeVarint(len(t.outputs)) // 출력 개수

	// 출력 개수만큼 반복하면서 각 출력을 직렬화한 결과를 result에 추가
	for _, output := range t.outputs {
		result = append(result, output.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

// Transaction 입력을 나타내는 구조체
type TxIn struct {
	prevIndex int    // 이전 트랜잭션의 출력 인덱스
	prevTx    string // 이전 트랜잭션의 해시
	scriptSig string // 해제 스크립트
	seqNo     int    // 시퀀스 번호
}

// TxIn 생성자 함수
func NewTxIn(prevIndex int, prevTx, scriptSig string, seqNos ...int) *TxIn {
	tx := &TxIn{
		prevTx:    prevTx,
		prevIndex: prevIndex,
		scriptSig: scriptSig,
		seqNo:     0xffffffff,
	}

	if len(seqNos) > 0 {
		tx.seqNo = seqNos[0]
	}

	return tx
}

// TxIn의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t TxIn) String() string {
	return fmt.Sprintf("%s:%d", t.prevTx, t.prevIndex)
}

// TxIn을 직렬화한 결과를 반환하는 함수
func (t TxIn) Serialize() []byte {
	result := utils.ReverseBytes(utils.StringToBytes(t.prevTx))
	result = append(result, utils.IntToLittleEndian(t.prevIndex, 4)...)
	// TODO: scriptSig를 직렬화한 결과를 result에 추가
	result = append(result, utils.IntToLittleEndian(t.seqNo, 4)...)
	return result
}

type TxOut struct {
	amount       int    // 금액
	scriptPubKey string // 잠금 스크립트
}

func NewTxOut(amount int, scriptPubKey string) *TxOut {
	return &TxOut{
		amount:       amount,
		scriptPubKey: scriptPubKey,
	}
}

func (t TxOut) String() string {
	return fmt.Sprintf("%d:%s", t.amount, t.scriptPubKey)
}

func (t TxOut) Serialize() []byte {
	result := utils.IntToLittleEndian(t.amount, 8)
	// TODO: scriptPubKey를 직렬화한 결과를 result에 추가
	return result
}
