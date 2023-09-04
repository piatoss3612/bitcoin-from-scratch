package tx

type Tx struct {
	// version: 트랜잭션의 버전
	// txIns: 트랜잭션의 입력 목록
	// txOuts: 트랜잭션의 출력 목록
	// lockTime: 트랜잭션의 유효 시점
	// testnet: 테스트넷인지 여부

}

func New() *Tx {
	return &Tx{}
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

type TxIn struct {
}
