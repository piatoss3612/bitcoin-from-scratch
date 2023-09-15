package main

import (
	"chapter07/ecc"
	"chapter07/script"
	"chapter07/tx"
	"chapter07/utils"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	// checkFee()
	// checkSig()
	// checkModifiedTx()
	checkConstructTx()
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

func checkConstructTx() {
	prevTx := "0d6fe5213c0b3291f208cba8bfb59b7476dffacc4e5cb66f6eb20a080843a299"
	prevIndex := 13
	txIn := tx.NewTxIn(prevTx, prevIndex, nil)

	changeAmount := int(0.33 * 1e8)
	changeH160, _ := utils.DecodeBase58("mzx5YhAH9kNHtcN481u6WkjeHjYtVeKVh2")
	changeScript := script.NewP2PKHScript(changeH160)
	changeOutput := tx.NewTxOut(changeAmount, changeScript)

	targetAmount := int(0.1 * 1e8)
	targetH160, _ := utils.DecodeBase58("mnrVtF8DWjMu839VW3rBfgYaAfKk8983Xf")
	targetScript := script.NewP2PKHScript(targetH160)
	targetOutput := tx.NewTxOut(targetAmount, targetScript)

	txObj := tx.NewTx(1, []*tx.TxIn{txIn}, []*tx.TxOut{changeOutput, targetOutput}, 0, true)
	fmt.Println(txObj)

}
