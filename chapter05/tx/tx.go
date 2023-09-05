package tx

import (
	"chapter05/utils"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type Tx struct {
	Version  int      // 트랜잭션의 버전
	Inputs   []*TxIn  // 트랜잭션의 입력 목록
	Outputs  []*TxOut // 트랜잭션의 출력 목록
	Locktime int      // 트랜잭션의 유효 시점
	Testnet  bool     // 테스트넷인지 여부
}

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

func (t Tx) String() string {
	return fmt.Sprintf("tx: %s\nversion: %d\ninputs: %s\noutputs: %s\nlocktime: %d",
		t.ID(), t.Version, t.Inputs, t.Outputs, t.Locktime)
}

func (t Tx) ID() string {
	return hex.EncodeToString(t.Hash())
}

func (t Tx) Hash() []byte {
	return utils.ReverseBytes(utils.Hash256(t.Serialize()))
}

func (t Tx) Serialize() []byte {
	result := utils.IntToLittleEndian(t.Version, 4)                    // 버전
	result = append(result, t.serializeInputs()...)                    // 입력 목록
	result = append(result, t.serializeOutputs()...)                   // 출력 목록
	result = append(result, utils.IntToLittleEndian(t.Locktime, 4)...) // 유효 시점
	return result
}

func (t Tx) serializeInputs() []byte {
	result := utils.EncodeVarint(len(t.Inputs)) // 입력 개수

	// 입력 개수만큼 반복하면서 각 입력을 직렬화한 결과를 result에 추가
	for _, input := range t.Inputs {
		result = append(result, input.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

func (t Tx) serializeOutputs() []byte {
	result := utils.EncodeVarint(len(t.Outputs)) // 출력 개수

	// 출력 개수만큼 반복하면서 각 출력을 직렬화한 결과를 result에 추가
	for _, output := range t.Outputs {
		result = append(result, output.Serialize()...)
	}

	return result // 직렬화한 결과를 반환
}

func (t Tx) Fee(fetcher *TxFetcher) (int, error) {
	totalIn, err := t.totalInput(fetcher)
	if err != nil {
		return 0, err
	}

	totalOut := t.totalOutput()

	return totalIn - totalOut, nil
}

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

func (t Tx) totalOutput() int {
	total := 0

	for _, output := range t.Outputs {
		total += output.Amount
	}

	return total
}

// Transaction 입력을 나타내는 구조체
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

func (t TxIn) FetchTx(fetcher *TxFetcher, testnet bool) (*Tx, error) {
	return fetcher.Fetch(t.PrevTx, testnet, false)
}

func (t TxIn) Value(fetcher *TxFetcher, testnet bool) (int, error) {
	tx, err := t.FetchTx(fetcher, testnet)
	if err != nil {
		return 0, err
	}

	return tx.Outputs[t.PrevIndex].Amount, nil
}

func (t TxIn) ScriptPubKey(fetcher *TxFetcher, testnet bool) (string, error) {
	tx, err := t.FetchTx(fetcher, testnet)
	if err != nil {
		return "", err
	}

	return tx.Outputs[t.PrevIndex].ScriptPubKey, nil
}

type TxOut struct {
	Amount       int    // 금액
	ScriptPubKey string // 잠금 스크립트
}

func NewTxOut(amount int, scriptPubKey string) *TxOut {
	return &TxOut{
		Amount:       amount,
		ScriptPubKey: scriptPubKey,
	}
}

func (t TxOut) String() string {
	return fmt.Sprintf("%d:%s", t.Amount, t.ScriptPubKey)
}

func (t TxOut) Serialize() []byte {
	result := utils.IntToLittleEndian(t.Amount, 8)
	// TODO: scriptPubKey를 직렬화한 결과를 result에 추가
	return result
}

type TxFetcher struct {
	client *http.Client
	cache  map[string]*Tx
}

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

func (tf TxFetcher) GetURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/"
	}
	return "https://blockstream.info/api/"
}

func (tf *TxFetcher) Fetch(txID string, testnet, fresh bool) (*Tx, error) {
	if fresh || tf.cache[txID] == nil {
		url := fmt.Sprintf("%s/tx/%s/hex", tf.GetURL(testnet), txID)

		resp, err := tf.client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error fetching %s: %s", txID, resp.Status)
		}

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		rawHex := make([]byte, hex.DecodedLen(len(raw)))

		_, err = hex.Decode(rawHex, raw)
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
			fmt.Println("here")
			tx, err = ParseTx(rawHex, testnet)
			if err != nil {
				return nil, err
			}
		}

		tf.cache[txID] = tx
	}
	tf.cache[txID].Testnet = testnet
	return tf.cache[txID], nil
}
