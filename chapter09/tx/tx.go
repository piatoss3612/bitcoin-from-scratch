package tx

import (
	"chapter09/ecc"
	"chapter09/script"
	"chapter09/utils"
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
	id, err := t.ID()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("tx: %s\nversion: %d\ninputs: %s\noutputs: %s\nlocktime: %d",
		id, t.Version, t.Inputs, t.Outputs, t.Locktime)
}

// 16진수 문자열로 표현된 트랜잭션 ID를 반환하는 함수
func (t Tx) ID() (string, error) {
	h, err := t.Hash()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h), nil
}

// 트랜잭션의 해시를 반환하는 함수
func (t Tx) Hash() ([]byte, error) {
	s, err := t.Serialize()
	if err != nil {
		return nil, err
	}

	return utils.ReverseBytes(utils.Hash256(s)), nil
}

// 트랜잭션을 직렬화한 결과를 반환하는 함수
func (t Tx) Serialize() ([]byte, error) {
	result := utils.IntToLittleEndian(t.Version, 4) // 버전

	in, err := t.serializeInputs() // 입력 목록
	if err != nil {
		return nil, err
	}

	out, err := t.serializeOutputs() // 출력 목록
	if err != nil {
		return nil, err
	}

	result = append(result, in...)
	result = append(result, out...)
	result = append(result, utils.IntToLittleEndian(t.Locktime, 4)...) // 유효 시점

	return result, nil
}

// 트랜잭션 입력 목록을 직렬화한 결과를 반환하는 함수
func (t Tx) serializeInputs() ([]byte, error) {
	inputs := t.Inputs

	result := utils.EncodeVarint(len(inputs)) // 입력 개수

	// 입력 개수만큼 반복하면서 각 입력을 직렬화한 결과를 result에 추가
	for _, input := range inputs {
		s, err := input.Serialize()
		if err != nil {
			return nil, err
		}

		result = append(result, s...)
	}

	return result, nil // 직렬화한 결과를 반환
}

// 트랜잭션 출력 목록을 직렬화한 결과를 반환하는 함수
func (t Tx) serializeOutputs() ([]byte, error) {
	outputs := t.Outputs

	result := utils.EncodeVarint(len(outputs)) // 출력 개수

	// 출력 개수만큼 반복하면서 각 출력을 직렬화한 결과를 result에 추가
	for _, output := range outputs {
		s, err := output.Serialize()
		if err != nil {
			return nil, err
		}
		result = append(result, s...)
	}

	return result, nil // 직렬화한 결과를 반환
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

// 트랜잭션의 서명해시를 반환하는 함수
// inputIndex는 서명해시를 만들 때 사용할 입력의 인덱스
// redeemScripts는 리딤 스크립트 목록
func (t Tx) SigHash(inputIndex int, redeemScripts ...*script.Script) ([]byte, error) {
	// 입력 인덱스가 트랜잭션의 입력 개수보다 크면 에러를 반환
	if inputIndex >= len(t.Inputs) {
		return nil, fmt.Errorf("input index %d greater than the number of inputs %d", inputIndex, len(t.Inputs))
	}

	s := utils.IntToLittleEndian(t.Version, 4) // 버전

	in, err := t.serializeInputsForSig(inputIndex, redeemScripts...) // 입력 목록, 입력의 인덱스와 리딤 스크립트 목록을 사용
	if err != nil {
		return nil, err
	}

	s = append(s, in...)

	out, err := t.serializeOutputs() // 출력 목록
	if err != nil {
		return nil, err
	}

	s = append(s, out...)

	s = append(s, utils.IntToLittleEndian(t.Locktime, 4)...) // 유효 시점

	s = append(s, utils.IntToLittleEndian(SIGHASH_ALL, 4)...) // SIGHASH_ALL (4바이트)

	h256 := utils.Hash256(s) // 해시를 생성

	return h256, nil // 해시를 반환
}

// 서명해시를 만들 때 사용할 입력 목록을 직렬화한 결과를 반환하는 함수
func (t Tx) serializeInputsForSig(inputIndex int, redeemScripts ...*script.Script) ([]byte, error) {
	inputs := t.Inputs

	result := utils.EncodeVarint(len(inputs)) // 입력 개수

	for i, input := range inputs {
		var scriptSig *script.Script // 해제 스크립트, 기본값은 nil

		if i == inputIndex { // 입력 인덱스가 inputIndex와 같으면
			if len(redeemScripts) > 0 { // 리딤 스크립트가 있으면
				scriptSig = redeemScripts[0] // 리딤 스크립트를 사용
			} else {
				scriptPubKey, err := input.ScriptPubKey(NewTxFetcher(), t.Testnet) // 이전 트랜잭션 출력의 잠금 스크립트를 가져옴
				if err != nil {
					return nil, err
				}

				scriptSig = scriptPubKey // 이전 트랜잭션 출력의 잠금 스크립트를 사용
			}
		}

		s, err := NewTxIn(input.PrevTx, input.PrevIndex, scriptSig, input.SeqNo).Serialize() // scriptSig를 사용하는 새로운 입력을 생성하고 직렬화
		if err != nil {
			return nil, err
		}

		result = append(result, s...) // 직렬화한 결과를 result에 추가
	}

	return result, nil // 직렬화한 결과를 반환
}

// 트랜잭션의 입력을 검증하는 함수
func (t Tx) VerifyInput(inputIndex int) (bool, error) {
	if inputIndex >= len(t.Inputs) {
		return false, fmt.Errorf("input index %d greater than the number of inputs %d", inputIndex, len(t.Inputs))
	}

	input := t.Inputs[inputIndex] // 입력을 가져옴

	scriptSig := input.ScriptSig // 해제 스크립트

	scriptPubKey, err := input.ScriptPubKey(NewTxFetcher(), t.Testnet) // 이전 트랜잭션 출력의 잠금 스크립트를 가져옴
	if err != nil {
		return false, err
	}

	var redeemScripts []*script.Script // 리딤 스크립트 목록

	if script.IsP2shScriptPubkey(scriptPubKey.Cmds) { // 이전 트랜잭션 출력의 잠금 스크립트가 P2SH 스크립트인 경우
		rawRedeem, ok := scriptSig.Cmds[len(scriptSig.Cmds)-1].([]byte) // 해제 스크립트의 마지막 원소가 리딤 스크립트
		if !ok {
			return false, fmt.Errorf("last element should be the redeem script")
		}

		redeemScript, _, err := script.Parse(append([]byte{byte(len(rawRedeem))}, rawRedeem...)) // 리딤 스크립트 파싱
		if err != nil {
			return false, err
		}

		redeemScripts = append(redeemScripts, redeemScript)
	}

	z, err := t.SigHash(inputIndex, redeemScripts...) // 서명해시를 가져옴
	if err != nil {
		return false, err
	}

	combined := scriptSig.Add(scriptPubKey) // 해제 스크립트와 잠금 스크립트를 결합

	return combined.Evaluate(z) // 결합한 스크립트를 평가
}

// 트랜잭션의 입력에 서명하는 함수
func (t Tx) SignInput(inputIndex int, privateKey ecc.PrivateKey, compressed bool) (bool, error) {
	if inputIndex >= len(t.Inputs) {
		return false, fmt.Errorf("input index %d greater than the number of inputs %d", inputIndex, len(t.Inputs))
	}

	z, err := t.SigHash(inputIndex) // 서명해시를 가져옴
	if err != nil {
		return false, err
	}

	point, err := privateKey.Sign(z) // 서명 생성
	if err != nil {
		return false, err
	}

	der := point.DER() // 서명을 DER 형식으로 직렬화

	sig := append(der, byte(SIGHASH_ALL))     // 직렬화한 서명에 해시 유형을 추가 (SIGHASH_ALL)
	sec := privateKey.Point().SEC(compressed) // 공개 키를 SEC 형식으로 직렬화 (압축)

	scriptSig := script.New(sig, sec) // 해제 스크립트 생성

	t.Inputs[inputIndex].ScriptSig = scriptSig // 트랜잭션의 입력에 해제 스크립트를 설정

	return t.VerifyInput(inputIndex) // 트랜잭션의 입력을 검증
}

// 트랜잭션을 검증하는 함수
func (t Tx) Verify() (bool, error) {
	fee, err := t.Fee(NewTxFetcher()) // 수수료를 가져옴
	if err != nil {
		return false, err
	}

	// 수수료가 음수이면 유효하지 않은 트랜잭션
	if fee < 0 {
		return false, nil
	}

	// 트랜잭션의 입력을 검증
	for i := 0; i < len(t.Inputs); i++ {
		ok, err := t.VerifyInput(i)
		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// 트랜잭션 입력을 나타내는 구조체
type TxIn struct {
	PrevIndex int            // 이전 트랜잭션의 출력 인덱스
	PrevTx    string         // 이전 트랜잭션의 해시
	ScriptSig *script.Script // 해제 스크립트
	SeqNo     int            // 시퀀스 번호
}

// TxIn 생성자 함수
func NewTxIn(prevTx string, prevIndex int, scriptSig *script.Script, seqNos ...int) *TxIn {
	tx := &TxIn{
		PrevTx:    prevTx,
		PrevIndex: prevIndex,
		ScriptSig: scriptSig,
		SeqNo:     0xffffffff,
	}

	if scriptSig == nil {
		tx.ScriptSig = script.New()
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
func (t TxIn) Serialize() ([]byte, error) {
	hexPrevTx, err := hex.DecodeString(t.PrevTx)
	if err != nil {
		return nil, err
	}

	result := utils.ReverseBytes(hexPrevTx)

	result = append(result, utils.IntToLittleEndian(t.PrevIndex, 4)...)

	serializedScript, err := t.ScriptSig.Serialize()
	if err != nil {
		return nil, err
	}

	result = append(result, serializedScript...)
	result = append(result, utils.IntToLittleEndian(t.SeqNo, 4)...)

	return result, nil
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
func (t TxIn) ScriptPubKey(fetcher *TxFetcher, testnet bool) (*script.Script, error) {
	tx, err := t.FetchTx(fetcher, testnet)
	if err != nil {
		return nil, err
	}

	return tx.Outputs[t.PrevIndex].ScriptPubKey, nil
}

// 트랜잭션 출력을 나타내는 구조체
type TxOut struct {
	Amount       int            // 금액
	ScriptPubKey *script.Script // 잠금 스크립트
}

// TxOut 생성자 함수
func NewTxOut(amount int, scriptPubKey *script.Script) *TxOut {
	out := &TxOut{
		Amount:       amount,
		ScriptPubKey: scriptPubKey,
	}

	if scriptPubKey == nil {
		out.ScriptPubKey = script.New()
	}

	return out
}

// TxOut의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t TxOut) String() string {
	return fmt.Sprintf("%d:%s", t.Amount, t.ScriptPubKey)
}

// TxOut을 직렬화한 결과를 반환하는 함수
func (t TxOut) Serialize() ([]byte, error) {
	result := utils.IntToLittleEndian(t.Amount, 8)

	serializedScript, err := t.ScriptPubKey.Serialize()
	if err != nil {
		return nil, err
	}

	result = append(result, serializedScript...)

	return result, nil
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
			parsedTx, err := ParseTx(rawHex, testnet)
			if err != nil {
				return nil, err
			}
			parsedTx.Locktime = utils.LittleEndianToInt(rawHex[len(rawHex)-4:])

			tx = parsedTx
		} else {
			parsedTx, err := ParseTx(rawHex, testnet)
			if err != nil {
				return nil, err
			}

			tx = parsedTx
		}

		expectedID, err := tx.ID() // 트랜잭션의 ID를 가져옴
		if err != nil {
			return nil, err
		}

		// 가져온 트랜잭션의 ID가 txID와 다르면 에러를 반환
		if expectedID != txID {
			return nil, fmt.Errorf("tx ID mismatch %s != %s", expectedID, txID)
		}

		tf.cache[txID] = tx // 가져온 트랜잭션을 tf.cache에 저장
	}

	tf.cache[txID].Testnet = testnet // 테스트넷 여부를 설정
	return tf.cache[txID], nil       // 트랜잭션을 반환
}
