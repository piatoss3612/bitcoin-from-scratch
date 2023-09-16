package main

import (
	"chapter07/ecc"
	"chapter07/script"
	"chapter07/tx"
	"chapter07/utils"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

func main() {
	// checkFee()
	// checkSig()
	// checkModifiedTx()
	// checkGenTx()
	// checkGenScriptSig()
	// checkSignInput()
	checkGenTestnetTx()
}

func checkFee() {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	parsedTx, _ := tx.ParseTx(rawTx, false)

	fetcher := tx.NewTxFetcher()

	fee, err := parsedTx.Fee(fetcher)
	if err != nil {
		panic(err)
	}

	fmt.Println(fee >= 0)
}

func checkSig() {
	sec, _ := hex.DecodeString("0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	der, _ := hex.DecodeString("3045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed")
	z, _ := hex.DecodeString("27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6")

	point, _ := ecc.ParsePoint(sec)
	sig, _ := ecc.ParseSignature(der)

	ok, err := point.Verify(z, sig)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok)
}

func checkModifiedTx() {
	modifiedTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88acfeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac1943060001000000")

	h256 := utils.Hash256(modifiedTx)
	z := big.NewInt(0).SetBytes(h256)

	sec, _ := hex.DecodeString("0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	der, _ := hex.DecodeString("3045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed")

	point, _ := ecc.ParsePoint(sec)
	sig, _ := ecc.ParseSignature(der)

	ok, err := point.Verify(z.Bytes(), sig)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok)
}

func checkGenTx() {
	prevTx := "0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299" // 이전 트랜잭션 ID
	prevIndex := 13                                                              // 이전 트랜잭션의 출력 인덱스
	txIn := tx.NewTxIn(prevTx, prevIndex, nil)                                   // 트랜잭션 입력 생성 (해제 스크립트는 비어있음)

	changeAmount := int(0.33 * 1e8)                                           // 출력 금액
	changeH160, _ := utils.DecodeBase58("mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2") // 잠금 스크립트를 생성할 주소
	changeScript := script.NewP2PKHScript(changeH160)                         // p2pkh 잠금 스크립트 생성
	changeOutput := tx.NewTxOut(changeAmount, changeScript)                   // 트랜잭션 출력 생성

	targetAmount := int(0.1 * 1e8)                                            // 출력 금액
	targetH160, _ := utils.DecodeBase58("mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf") // 잠금 스크립트를 생성할 주소
	targetScript := script.NewP2PKHScript(targetH160)                         // p2pkh 잠금 스크립트 생성
	targetOutput := tx.NewTxOut(targetAmount, targetScript)                   // 트랜잭션 출력 생성

	txObj := tx.NewTx(1, []*tx.TxIn{txIn}, []*tx.TxOut{changeOutput, targetOutput}, 0, true) // 트랜잭션 생성
	fmt.Println(txObj)
}

func checkGenScriptSig() {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	parsedTx, _ := tx.ParseTx(rawTx, false) // 트랜잭션 파싱

	z, _ := parsedTx.SigHash(0) // 서명 해시 생성

	privateKey, _ := ecc.NewS256PrivateKey(big.NewInt(8675309).Bytes()) // 개인 키 생성

	point, _ := privateKey.Sign(z) // 서명 생성

	der := point.DER() // 서명을 DER 형식으로 변환

	sig := append(der, byte(tx.SIGHASH_ALL)) // DER 서명에 해시 타입을 추가
	sec := privateKey.Point().SEC(true)      // 공개 키를 압축된 SEC 형식으로 변환
	scriptSig := script.New(sig, sec)        // 해제 스크립트 생성

	parsedTx.Inputs[0].ScriptSig = scriptSig // 해제 스크립트를 트랜잭션 입력에 추가
	encoded, err := parsedTx.Serialize()     // 트랜잭션 직렬화
	if err != nil {
		panic(err)
	}

	expected := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006a47304402207db2402a3311a3b845b038885e3dd889c08126a8570f26a844e3e4049c482a11022010178cdca4129eacbeab7c44648bf5ac1f9cac217cd609d216ec2ebc8d242c0a012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67feffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	fmt.Println(strings.Compare(hex.EncodeToString(encoded), expected))
}

func checkSignInput() {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	parsedTx, _ := tx.ParseTx(rawTx, false) // 트랜잭션 파싱

	fmt.Println(parsedTx.Inputs[0].ScriptSig)

	privateKey, _ := ecc.NewS256PrivateKey(big.NewInt(8675309).Bytes()) // this private key might invalid for mainnet transaction

	ok, err := parsedTx.SignInput(0, privateKey, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok) // false
}

func checkGenTestnetTx() {
	secret := utils.LittleEndianToBigInt(utils.Hash256(utils.StringToBytes("piatoss rules the world")))
	privateKey, _ := ecc.NewS256PrivateKey(secret.Bytes())

	address := privateKey.Point().Address(true, true)

	prevTx := "e770e0b481166da7d0d139c855e86633a12dbd4fa9b97f33a31fc9a458f8ddd7"
	prevIndex := 0
	txIn := tx.NewTxIn(prevTx, prevIndex, nil)

	balance := 1193538

	changeAmount := balance - (balance * 6 / 10) // 40% of balance
	changeH160, _ := utils.DecodeBase58(address)
	changeScript := script.NewP2PKHScript(changeH160)
	changeOutput := tx.NewTxOut(changeAmount, changeScript)

	targetAmount := balance * 6 / 10 // 60% of balance
	targetH160, _ := utils.DecodeBase58("mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv")
	targetScript := script.NewP2PKHScript(targetH160)
	targetOutput := tx.NewTxOut(targetAmount, targetScript)

	txObj := tx.NewTx(1, []*tx.TxIn{txIn}, []*tx.TxOut{changeOutput, targetOutput}, 0, true)

	ok, err := txObj.SignInput(0, privateKey, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(ok)

	serializedTx, err := txObj.Serialize()
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(serializedTx))
}
