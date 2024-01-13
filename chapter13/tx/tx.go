package tx

import (
	"bytes"
	"chapter13/ecc"
	"chapter13/script"
	"chapter13/utils"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 트랜잭션을 나타내는 구조체
type Tx struct {
	Version      int      // 트랜잭션의 버전
	Inputs       []*TxIn  // 트랜잭션의 입력 목록
	Outputs      []*TxOut // 트랜잭션의 출력 목록
	Locktime     int      // 트랜잭션의 유효 시점
	Testnet      bool     // 테스트넷인지 여부
	Segwit       bool     // 세그윗인지 여부
	HashPrevOuts []byte   // 이전 트랜잭션 출력의 해시
	HashSequence []byte   // 시퀀스 번호의 해시
	HashOutputs  []byte   // 출력의 해시
}

// Tx 생성자 함수
func NewTx(version int, inputs []*TxIn, outputs []*TxOut, locktime int, testnet, segwit bool) *Tx {
	tx := &Tx{
		Version:  version,
		Inputs:   inputs,
		Outputs:  outputs,
		Locktime: locktime,
		Testnet:  testnet,
		Segwit:   segwit,
	}

	return tx
}

// 트랜잭션의 문자열 표현을 반환하는 함수 (fmt.Stringer 인터페이스 구현)
func (t Tx) String() string {
	id, _ := t.ID()
	return fmt.Sprintf("tx: %s\nversion: %d\ninputs: %s\noutputs: %s\nlocktime: %d\nsegwit: %t\n",
		id, t.Version, t.Inputs, t.Outputs, t.Locktime, t.Segwit)
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
	s, err := t.serializeLegacy()
	if err != nil {
		return nil, err
	}

	return utils.ReverseBytes(utils.Hash256(s)), nil
}

// 트랜잭션을 직렬화한 결과를 반환하는 함수
func (t Tx) Serialize() ([]byte, error) {
	if t.Segwit { // 세그윗 트랜잭션인 경우
		return t.serializeSegwit()
	}

	return t.serializeLegacy()
}

func (t Tx) serializeLegacy() ([]byte, error) {
	buf := new(bytes.Buffer)

	_, err := buf.Write(utils.IntToLittleEndian(t.Version, 4)) // 버전 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Inputs))) // 입력 개수 (가변 정수)
	if err != nil {
		return nil, err
	}

	for _, input := range t.Inputs {
		s, err := input.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 입력
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Outputs))) // 출력 개수
	if err != nil {
		return nil, err
	}

	for _, output := range t.Outputs {
		s, err := output.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 출력
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(utils.IntToLittleEndian(t.Locktime, 4)) // 유효 시점
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t Tx) serializeSegwit() ([]byte, error) {
	buf := new(bytes.Buffer)

	_, err := buf.Write(utils.IntToLittleEndian(t.Version, 4)) // 버전 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte{0x00, 0x01}) // 마커 (0x00), 플래그 (0x01)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Inputs))) // 입력 개수 (가변 정수)
	if err != nil {
		return nil, err
	}

	for _, input := range t.Inputs {
		s, err := input.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 입력
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Outputs))) // 출력 개수
	if err != nil {
		return nil, err
	}

	for _, output := range t.Outputs {
		s, err := output.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 출력
		if err != nil {
			return nil, err
		}
	}

	witnesses, err := t.serializeWitnesses() // 증인 데이터를 직렬화
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(witnesses) // 증인 데이터
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(t.Locktime, 4)) // 유효 시점
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 증인 데이터를 직렬화한 결과를 반환하는 함수
func (t Tx) serializeWitnesses() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, input := range t.Inputs {
		_, err := buf.Write(utils.IntToLittleEndian(len(input.Witness), 1)) // 입력의 증인 데이터 개수
		if err != nil {
			return nil, err
		}

		for _, item := range input.Witness {
			if len(item) == 0 {
				_, err := buf.Write(utils.IntToLittleEndian(0, 1)) // 증인 데이터의 길이가 0이면 0x00을 추가
				if err != nil {
					return nil, err
				}
				continue
			}

			_, err := buf.Write(utils.EncodeVarint(len(item))) // 증인 데이터의 길이 (가변 정수)
			if err != nil {
				return nil, err
			}

			_, err = buf.Write(item) // 증인 데이터
			if err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil // 직렬화한 결과를 반환
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

	buf := new(bytes.Buffer)

	_, err := buf.Write(utils.IntToLittleEndian(t.Version, 4)) // 버전 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Inputs))) // 입력 개수 (가변 정수)
	if err != nil {
		return nil, err
	}

	for i, input := range t.Inputs {
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
		} else {
			scriptSig = nil
		}

		s, err := NewTxIn(input.PrevTx, input.PrevIndex, scriptSig, input.SeqNo).Serialize() // scriptSig를 사용하는 새로운 입력을 생성하고 직렬화
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 직렬화한 결과를 buf에 추가
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(utils.EncodeVarint(len(t.Outputs))) // 출력 개수
	if err != nil {
		return nil, err
	}

	for _, output := range t.Outputs {
		s, err := output.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(s) // 출력
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(utils.IntToLittleEndian(t.Locktime, 4)) // 유효 시점 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(SIGHASH_ALL, 4)) // SIGHASH_ALL (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	h256 := utils.Hash256(buf.Bytes()) // 해시를 생성

	return h256, nil // 해시를 반환
}

// BIP143에 따라 트랜잭션의 서명해시를 반환하는 함수
func (t *Tx) SigHashBIP143(inputIndex int, redeemScript *script.Script, witnessScript *script.Script) ([]byte, error) {
	if inputIndex >= len(t.Inputs) {
		return nil, fmt.Errorf("input index %d greater than the number of inputs %d", inputIndex, len(t.Inputs))
	}

	txin := t.Inputs[inputIndex] // 입력을 가져옴

	buf := new(bytes.Buffer)

	_, err := buf.Write(utils.IntToLittleEndian(t.Version, 4)) // 버전 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	hashPrevOuts, err := t.hashPrevouts() // -> 잘못 구현한 부분
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(hashPrevOuts) // 이전 트랜잭션 출력의 해시 (32바이트)
	if err != nil {
		return nil, err
	}

	hashSequence, err := t.hashSequence()
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(hashSequence) // 시퀀스 번호의 해시 (32바이트)
	if err != nil {
		return nil, err
	}

	hexPrevTx, err := hex.DecodeString(txin.PrevTx) // 문자열을 16진수로 디코딩 (이 부분이 빠져있었음)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.ReverseBytes(hexPrevTx)) // 이전 트랜잭션의 해시 (32바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(txin.PrevIndex, 4)) // 이전 트랜잭션의 출력 인덱스 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	var scriptCode []byte

	if witnessScript != nil {
		scriptCode, err = witnessScript.Serialize()
		if err != nil {
			return nil, err
		}
	} else if redeemScript != nil {
		scriptCode, err = script.NewP2pkhScript(redeemScript.Cmds[1].Elem).Serialize()
		if err != nil {
			return nil, err
		}
	} else {
		scriptPubkey, err := txin.ScriptPubKey(NewTxFetcher(), t.Testnet)
		if err != nil {
			return nil, err
		}

		scriptCode, err = script.NewP2pkhScript(scriptPubkey.Cmds[1].Elem).Serialize()
		if err != nil {
			return nil, err
		}
	}

	_, err = buf.Write(scriptCode) // 스크립트 코드
	if err != nil {
		return nil, err
	}

	val, err := txin.Value(NewTxFetcher(), t.Testnet)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(val, 8)) // 이전 트랜잭션 출력의 금액 (8바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(txin.SeqNo, 4)) // 시퀀스 번호 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	hashOutputs, err := t.hashOutputs()
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(hashOutputs) // 출력의 해시 (32바이트)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(t.Locktime, 4)) // 유효 시점 (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(utils.IntToLittleEndian(SIGHASH_ALL, 4)) // SIGHASH_ALL (4바이트, 리틀엔디언)
	if err != nil {
		return nil, err
	}

	h256 := utils.Hash256(buf.Bytes()) // 해시를 생성

	return h256, nil // 해시를 반환
}

func (t *Tx) hashPrevouts() ([]byte, error) {
	if t.HashPrevOuts != nil {
		return t.HashPrevOuts, nil
	}

	hashBuf := new(bytes.Buffer)
	seqBuf := new(bytes.Buffer)

	for _, input := range t.Inputs {
		hexPrevTx, err := hex.DecodeString(input.PrevTx) // 문자열을 16진수로 디코딩 (이 부분이 빠져있었음)
		if err != nil {
			return nil, err
		}

		_, err = hashBuf.Write(utils.ReverseBytes(hexPrevTx))
		if err != nil {
			return nil, err
		}

		_, err = hashBuf.Write(utils.IntToLittleEndian(input.PrevIndex, 4))
		if err != nil {
			return nil, err
		}

		_, err = seqBuf.Write(utils.IntToLittleEndian(input.SeqNo, 4))
		if err != nil {
			return nil, err
		}
	}

	t.HashPrevOuts = utils.Hash256(hashBuf.Bytes())
	t.HashSequence = utils.Hash256(seqBuf.Bytes())

	return t.HashPrevOuts, nil
}

func (t *Tx) hashSequence() ([]byte, error) {
	if t.HashSequence != nil {
		return t.HashSequence, nil
	}

	_, err := t.hashPrevouts()
	if err != nil {
		return nil, err
	}

	return t.HashSequence, nil
}

func (t *Tx) hashOutputs() ([]byte, error) {
	if t.HashOutputs != nil {
		return t.HashOutputs, nil
	}

	buf := new(bytes.Buffer)

	for _, output := range t.Outputs {
		b, err := output.Serialize()
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(b)
		if err != nil {
			return nil, err
		}
	}

	t.HashOutputs = utils.Hash256(buf.Bytes())

	return t.HashOutputs, nil
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

	var z []byte
	var witness [][]byte

	if script.IsP2shScriptPubkey(scriptPubKey.Cmds) { // 이전 트랜잭션 출력의 잠금 스크립트가 p2sh 스크립트인 경우
		command := scriptSig.Cmds[len(scriptSig.Cmds)-1] // 해제 스크립트의 마지막 명령어
		if command.IsOpCode {                            // 마지막 명령어가 op 코드인 경우 에러를 반환
			return false, fmt.Errorf("last command should be redeem script")
		}

		rawRedeem := append(utils.IntToLittleEndian(len(command.Elem), 1), command.Elem...) // 리딤 스크립트를 가져옴
		redeemScript, _, err := script.Parse(rawRedeem)                                     // 리딤 스크립트를 파싱
		if err != nil {
			return false, err
		}

		if script.IsP2wpkhScriptPubkey(redeemScript.Cmds) { // 리딤 스크립트가 p2wpkh 스크립트인 경우
			z, err = t.SigHashBIP143(inputIndex, redeemScript, nil) // BIP143에 따라 서명해시를 생성
			if err != nil {
				return false, err
			}

			witness = input.Witness
		} else if script.IsP2wshScriptPubkey(redeemScript.Cmds) { // 리딤 스크립트가 p2wsh 스크립트인 경우
			cmd := input.Witness[len(input.Witness)-1]                         // 증인 데이터의 마지막 명령어
			rawWitness := append(utils.IntToLittleEndian(len(cmd), 1), cmd...) // 증인 데이터를 가져옴
			witnessScript, _, err := script.Parse(rawWitness)                  // 증인 데이터를 파싱
			if err != nil {
				return false, err
			}

			z, err = t.SigHashBIP143(inputIndex, nil, witnessScript) // BIP143에 따라 서명해시를 생성
			if err != nil {
				return false, err
			}

			witness = input.Witness
		} else {
			z, err = t.SigHash(inputIndex, redeemScript) // 서명해시를 생성
			if err != nil {
				return false, err
			}

			witness = nil
		}
	} else {
		if script.IsP2wpkhScriptPubkey(scriptPubKey.Cmds) {
			z, err = t.SigHashBIP143(inputIndex, nil, nil)
			if err != nil {
				return false, err
			}

			witness = input.Witness
		} else if script.IsP2wshScriptPubkey(scriptPubKey.Cmds) {
			cmd := input.Witness[len(input.Witness)-1]
			rawWitness := append(utils.IntToLittleEndian(len(cmd), 1), cmd...)
			witnessScript, _, err := script.Parse(rawWitness)
			if err != nil {
				return false, err
			}

			z, err = t.SigHashBIP143(inputIndex, nil, witnessScript)
			if err != nil {
				return false, err
			}

			witness = input.Witness
		} else {
			z, err = t.SigHash(inputIndex)
			if err != nil {
				return false, err
			}

			witness = nil
		}
	}

	combined := scriptSig.Add(scriptPubKey) // 해제 스크립트와 잠금 스크립트를 결합

	return combined.Evaluate(z, witness) // 결합한 스크립트를 평가
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

	scriptSig := script.New(script.NewElem(sig), script.NewElem(sec)) // 해제 스크립트 생성

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

// 트랜잭션이 코인베이스 트랜잭션인지 여부를 반환하는 함수
func (t Tx) IsCoinbase() bool {
	return len(t.Inputs) == 1 && // 입력 개수가 1이고
		strings.EqualFold(t.Inputs[0].PrevTx, hex.EncodeToString(bytes.Repeat([]byte{0x00}, 32))) && // 이전 트랜잭션이 0x00으로 채워진 32바이트이고
		t.Inputs[0].PrevIndex == 0xffffffff // 이전 트랜잭션의 출력 인덱스가 0xffffffff인 경우 코인베이스 트랜잭션
}

// 코인베이스 트랜잭션의 해제 스크립트에 포함된 높이를 반환하는 함수
func (t Tx) CoinbaseHeight() (bool, int) {
	if !t.IsCoinbase() {
		return false, 0
	}

	// 코인베이스 트랜잭션의 해제 스크립트에서 높이를 가져옴
	scriptSig := t.Inputs[0].ScriptSig

	if len(scriptSig.Cmds) < 2 || len(scriptSig.Cmds[0].Elem) == 0 {
		return false, 0
	}

	heightBytes := scriptSig.Cmds[0].Elem // 높이를 나타내는 바이트

	return true, utils.LittleEndianToInt(heightBytes) // 리틀엔디언으로 인코딩된 높이를 반환
}

// 트랜잭션 입력을 나타내는 구조체
type TxIn struct {
	PrevIndex int            // 이전 트랜잭션의 출력 인덱스
	PrevTx    string         // 이전 트랜잭션의 해시
	ScriptSig *script.Script // 해제 스크립트
	SeqNo     int            // 시퀀스 번호
	Witness   [][]byte       // 증인 데이터 목록 (세그윗 트랜잭션인 경우)
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

	result := utils.ReverseBytes(hexPrevTx) // 이전 트랜잭션의 해시 (리틀엔디언)

	result = append(result, utils.IntToLittleEndian(t.PrevIndex, 4)...) // 이전 트랜잭션의 출력 인덱스 (리틀엔디언)

	serializedScript, err := t.ScriptSig.Serialize()
	if err != nil {
		return nil, err
	}

	result = append(result, serializedScript...)                    // 해제 스크립트
	result = append(result, utils.IntToLittleEndian(t.SeqNo, 4)...) // 시퀀스 번호 (리틀엔디언)

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
	result := utils.IntToLittleEndian(t.Amount, 8) // 금액 (8바이트, 리틀엔디언)

	serializedScript, err := t.ScriptPubKey.Serialize() // 잠금 스크립트를 직렬화
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

		tx, err := ParseTx(rawHex, testnet) // rawHex를 파싱
		if err != nil {
			return nil, err
		}

		// 세그윗을 적용했으므로 필요 없는 코드. 책에는 나와있지 않음
		// if rawHex[4] == 0x00 {
		// 	rawHex = append(rawHex[:4], rawHex[6:]...)
		// 	parsedTx, err := ParseTx(rawHex, testnet)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	parsedTx.Locktime = utils.LittleEndianToInt(rawHex[len(rawHex)-4:])

		// 	tx = parsedTx
		// } else {
		// 	parsedTx, err := ParseTx(rawHex, testnet)
		// 	if err != nil {
		// 		return nil, err
		// 	}

		// 	tx = parsedTx
		// }

		var computed string

		if tx.Segwit {
			computed, err = tx.ID()
			if err != nil {
				return nil, err
			}
		} else {
			computed = hex.EncodeToString(utils.ReverseBytes(utils.Hash256(rawHex)))
		}

		// 가져온 트랜잭션의 ID가 txID와 다르면 에러를 반환
		if computed != txID {
			return nil, fmt.Errorf("tx ID mismatch %s != %s", computed, txID)
		}

		tf.cache[txID] = tx // 가져온 트랜잭션을 tf.cache에 저장
	}

	tf.cache[txID].Testnet = testnet // 테스트넷 여부를 설정
	return tf.cache[txID], nil       // 트랜잭션을 반환
}
