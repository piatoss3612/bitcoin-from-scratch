package tx

import "fmt"

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
	// TODO: implement
	panic("not implemented")
}

func (t Tx) ID() []byte {
	// TODO: implement
	panic("not implemented")
}

func (t Tx) Hash() []byte {
	// TODO: implement
	panic("not implemented")
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
