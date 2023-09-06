package tx

import (
	"chapter06/utils"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

// 트랜잭션을 나타내는 구조체
type Tx struct {
	Version  int      // 트랜잭션의 버전
	Inputs   []*TxIn  // 트랜잭션의 입력 목록
	Outputs  []*TxOut // 트랜잭션의 출력 목록
	Locktime int      // 트랜잭션의 유효 시점
	Testnet  bool     // 테스트넷인지 여부
}

// Tx 생성자 함수
func NewTx(version int, inputs []*TxIn, outputs []*TxOut, locktime int, testnet bool) *Tx {
	tx := &Tx{
		Version:  version,
		Inputs:   inputs,
		Outputs:  outputs,
		Locktime: locktime,
		Testnet:  testnet,
	}

	return tx
}

// 트랜잭션의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t Tx) String() string {
	return fmt.Sprintf("tx: %s\nversion: %d\ninputs: %s\noutputs: %s\nlocktime: %d",
		t.ID(), t.Version, t.Inputs, t.Outputs, t.Locktime)
}

// 16진수 문자열로 표현된 트랜잭션 ID를 반환하는 함수
func (t Tx) ID() string {
	return hex.EncodeToString(t.Hash())
}

// 트랜잭션의 해시를 반환하는 함수
func (t Tx) Hash() []byte {
	return utils.ReverseBytes(utils.Hash256(t.Serialize()))
}

// 트랜잭션을 직렬화한 결과를 반환하는 함수
func (t Tx) Serialize() []byte {
	result := utils.IntToLittleEndian(t.Version, 4)                    // 버전
	result = append(result, t.serializeInputs()...)                    // 입력 목록
	result = append(result, t.serializeOutputs()...)                   // 출력 목록
	result = append(result, utils.IntToLittleEndian(t.Locktime, 4)...) // 유효 시점
	return result
}

// 트랜잭션 입력 목록을 직렬화한 결과를 반환하는 함수
func (t Tx) serializeInputs() []byte {
	result := utils.EncodeVarint(len(t.Inputs)) // 입력 개수

	// 입력 개수만큼 반복하면서 각 입력을 직렬화한 결과를 result에 추가
	for _, input := range t.Inputs {
		result = append(result, input.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

// 트랜잭션 출력 목록을 직렬화한 결과를 반환하는 함수
func (t Tx) serializeOutputs() []byte {
	result := utils.EncodeVarint(len(t.Outputs)) // 출력 개수

	// 출력 개수만큼 반복하면서 각 출력을 직렬화한 결과를 result에 추가
	for _, output := range t.Outputs {
		result = append(result, output.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

// 트랜잭션의 수수료를 반환하는 함수
func (t Tx) Fee(fetcher *TxFetcher) (int, error) {
	totalIn, err := t.totalInput(fetcher)
	if err != nil {
		return 0, err
	}

	totalOut := t.totalOutput()

	return totalIn - totalOut, nil
}

// 트랜잭션 입력의 비트코인 총량을 반환하는 함수
func (t Tx) totalInput(fetcher *TxFetcher) (int, error) {
	total := 0

	for _, input := range t.Inputs {
		value, err := input.Value(fetcher, t.Testnet)
		if err != nil {
			return 0, err
		}

		total += value
	}

	return total, nil
}

// 트랜잭션 출력의 비트코인 총량을 반환하는 함수
func (t Tx) totalOutput() int {
	total := 0

	for _, output := range t.Outputs {
		total += output.Amount
	}

	return total
}

// 트랜잭션 입력을 나타내는 구조체
type TxIn struct {
	PrevIndex int    // 이전 트랜잭션의 출력 인덱스
	PrevTx    string // 이전 트랜잭션의 해시
	ScriptSig string // 해제 스크립트
	SeqNo     int    // 시퀀스 번호
}

// TxIn 생성자 함수
func NewTxIn(prevIndex int, prevTx, scriptSig string, seqNos ...int) *TxIn {
	tx := &TxIn{
		PrevTx:    prevTx,
		PrevIndex: prevIndex,
		ScriptSig: scriptSig,
		SeqNo:     0xffffffff,
	}

	if len(seqNos) > 0 {
		tx.SeqNo = seqNos[0]
	}

	return tx
}

// TxIn의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t TxIn) String() string {
	return fmt.Sprintf("%s:%d", t.PrevTx, t.PrevIndex)
}

// TxIn을 직렬화한 결과를 반환하는 함수
func (t TxIn) Serialize() []byte {
	result := utils.ReverseBytes(utils.StringToBytes(t.PrevTx))
	result = append(result, utils.IntToLittleEndian(t.PrevIndex, 4)...)
	// TODO: scriptSig를 직렬화한 결과를 result에 추가
	result = append(result, utils.IntToLittleEndian(t.SeqNo, 4)...)
	return result
}

// TxIn의 이전 트랜잭션을 가져오는 함수
func (t TxIn) FetchTx(fetcher *TxFetcher, testnet bool) (*Tx, error) {
	return fetcher.Fetch(t.PrevTx, testnet, false)
}

// TxIn의 이전 트랜잭션 출력의 금액을 반환하는 함수
func (t TxIn) Value(fetcher *TxFetcher, testnet bool) (int, error) {
	tx, err := t.FetchTx(fetcher, testnet)
	if err != nil {
		return 0, err
	}

	return tx.Outputs[t.PrevIndex].Amount, nil
}

// TxIn의 이전 트랜잭션 출력의 잠금 스크립트를 반환하는 함수
func (t TxIn) ScriptPubKey(fetcher *TxFetcher, testnet bool) (string, error) {
	tx, err := t.FetchTx(fetcher, testnet)
	if err != nil {
		return "", err
	}

	return tx.Outputs[t.PrevIndex].ScriptPubKey, nil
}

// 트랜잭션 출력을 나타내는 구조체
type TxOut struct {
	Amount       int    // 금액
	ScriptPubKey string // 잠금 스크립트
}

// TxOut 생성자 함수
func NewTxOut(amount int, scriptPubKey string) *TxOut {
	return &TxOut{
		Amount:       amount,
		ScriptPubKey: scriptPubKey,
	}
}

// TxOut의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t TxOut) String() string {
	return fmt.Sprintf("%d:%s", t.Amount, t.ScriptPubKey)
}

// TxOut을 직렬화한 결과를 반환하는 함수
func (t TxOut) Serialize() []byte {
	result := utils.IntToLittleEndian(t.Amount, 8)
	// TODO: scriptPubKey를 직렬화한 결과를 result에 추가
	return result
}

// 트랜잭션을 가져오는 구조체
type TxFetcher struct {
	client *http.Client
	cache  map[string]*Tx
}

// TxFetcher 생성자 함수
func NewTxFetcher(clients ...*http.Client) *TxFetcher {
	tf := &TxFetcher{
		client: &http.Client{},
		cache:  make(map[string]*Tx),
	}

	if len(clients) > 0 && clients[0] != nil {
		tf.client = clients[0]
	}

	return tf
}

// 트랜잭션을 가져올 URL을 반환하는 함수
func (tf TxFetcher) GetURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/"
	}
	return "https://blockstream.info/api/"
}

// 트랜잭션을 가져오는 함수
func (tf *TxFetcher) Fetch(txID string, testnet, fresh bool) (*Tx, error) {
	// fresh가 true이거나 tf.cache에 txID가 없으면 트랜잭션을 가져옴
	if fresh || tf.cache[txID] == nil {
		url := fmt.Sprintf("%s/tx/%s/hex", tf.GetURL(testnet), txID)

		resp, err := tf.client.Get(url) // GET 요청을 보내 트랜잭션의 16진수 직렬화 결과를 가져옴
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error fetching %s: %s", txID, resp.Status)
		}

		raw, err := io.ReadAll(resp.Body) // 응답 바디를 읽어서 raw에 저장
		if err != nil {
			return nil, err
		}

		rawHex := make([]byte, hex.DecodedLen(len(raw)))

		_, err = hex.Decode(rawHex, raw) // raw를 16진수로 디코딩한 결과를 rawHex에 저장
		if err != nil {
			return nil, err
		}

		var tx *Tx

		if rawHex[4] == 0x00 {
			rawHex = append(rawHex[:4], rawHex[6:]...)
			tx, err = ParseTx(rawHex, testnet)
			if err != nil {
				return nil, err
			}
			tx.Locktime = utils.LittleEndianToInt(rawHex[len(rawHex)-4:])
		} else {
			tx, err = ParseTx(rawHex, testnet)
			if err != nil {
				return nil, err
			}
		}

		// 가져온 트랜잭션의 ID가 txID와 다르면 에러를 반환
		if tx.ID() != txID {
			return nil, fmt.Errorf("tx ID mismatch %s != %s", tx.ID(), txID)
		}

		tf.cache[txID] = tx // 가져온 트랜잭션을 tf.cache에 저장
	}
	tf.cache[txID].Testnet = testnet // 테스트넷 여부를 설정
	return tf.cache[txID], nil       // 트랜잭션을 반환
}
