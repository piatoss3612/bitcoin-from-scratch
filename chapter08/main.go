package main

import (
	"chapter08/ecc"
	"chapter08/script"
	"chapter08/tx"
	"chapter08/utils"
	"encoding/hex"
	"fmt"
)

func main() {
	// testH160ToP2shAddress()
	// testGenSigHashWithModifiedTx()
	// testVerifyFirstSigWithModifiedTx()
	// testVerifySecondSigWithModifiedTx()
	testVerifyInput()
}

func testH160ToP2shAddress() {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	fmt.Println(utils.H160ToP2shAddress(h160, false))
}

func testGenSigHashWithModifiedTx() {
	modifiedTx, _ := hex.DecodeString("0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152aeffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2cc15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c00000000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e6b3c192ecfb52cc8984ee7b6c56870000000001000000")
	s256 := utils.Hash256(modifiedTx)
	z := utils.BytesToBigInt(s256)
	fmt.Println(z.Text(16))
}

func testVerifyFirstSigWithModifiedTx() {
	modifiedTx, _ := hex.DecodeString("0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152aeffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2cc15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c00000000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e6b3c192ecfb52cc8984ee7b6c56870000000001000000")
	s256 := utils.Hash256(modifiedTx)
	z := utils.BytesToBigInt(s256)
	sec, _ := hex.DecodeString("022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb70")
	der, _ := hex.DecodeString("3045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a89937")

	point, _ := ecc.ParsePoint(sec)
	sig, _ := ecc.ParseSignature(der)

	ok, err := point.Verify(z.Bytes(), sig)
	fmt.Println(ok, err)
}

func testVerifySecondSigWithModifiedTx() {
	// 원래 트랜잭션
	txBytes, _ := hex.DecodeString("0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000db00483045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a8993701483045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152aeffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2cc15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c00000000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e6b3c192ecfb52cc8984ee7b6c568700000000")
	// 공개키
	sec, _ := hex.DecodeString("03b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb71")
	// 원래 트랜잭션의 첫 번째 입력의 두 번째 서명
	der, _ := hex.DecodeString("3045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e754022")
	hexRedeemScript, _ := hex.DecodeString("475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152ae")

	redeemScript, _, _ := script.Parse(hexRedeemScript) // 리딤 스크립트 파싱

	txObj, _ := tx.ParseTx(txBytes, false) // 원래 트랜잭션 파싱
	fmt.Println(txObj)

	s := utils.IntToLittleEndian(txObj.Version, 4)          // 트랜잭션의 버전 직렬화
	s = append(s, utils.EncodeVarint(len(txObj.Inputs))...) // 트랜잭션의 입력 개수 직렬화

	i := txObj.Inputs[0]
	newInput := tx.NewTxIn(i.PrevTx, i.PrevIndex, redeemScript, i.SeqNo) // 리딤 스크립트를 해제 스크립트 자리에 넣은 새로운 입력 생성
	serialized, _ := newInput.Serialize()                                // 새로운 입력 직렬화
	s = append(s, serialized...)

	s = append(s, utils.EncodeVarint(len(txObj.Outputs))...) // 트랜잭션의 출력 개수 직렬화
	for _, o := range txObj.Outputs {                        // 트랜잭션의 출력 직렬화
		serialized, _ := o.Serialize()
		s = append(s, serialized...)
	}

	s = append(s, utils.IntToLittleEndian(txObj.Locktime, 4)...) // 트랜잭션의 락타임 직렬화
	s = append(s, utils.IntToLittleEndian(tx.SIGHASH_ALL, 4)...) // 트랜잭션의 서명 해시 직렬화

	z := utils.BytesToBigInt(utils.Hash256(s)) // 서명 해시 생성

	point, _ := ecc.ParsePoint(sec)   // 공개키 파싱
	sig, _ := ecc.ParseSignature(der) // 서명 파싱

	ok, err := point.Verify(z.Bytes(), sig) // 서명 검증
	fmt.Println(ok, err)
}

func testVerifyInput() {
	txBytes, _ := hex.DecodeString("0100000001868278ed6ddfb6c1ed3ad5f8181eb0c7a385aa0836f01d5e4789e6bd304d87221a000000db00483045022100dc92655fe37036f47756db8102e0d7d5e28b3beb83a8fef4f5dc0559bddfb94e02205a36d4e4e6c7fcd16658c50783e00c341609977aed3ad00937bf4ee942a8993701483045022100da6bee3c93766232079a01639d07fa869598749729ae323eab8eef53577d611b02207bef15429dcadce2121ea07f233115c6f09034c0be68db99980b9a6c5e75402201475221022626e955ea6ea6d98850c994f9107b036b1334f18ca8830bfff1295d21cfdb702103b287eaf122eea69030a0e9feed096bed8045c8b98bec453e1ffac7fbdbd4bb7152aeffffffff04d3b11400000000001976a914904a49878c0adfc3aa05de7afad2cc15f483a56a88ac7f400900000000001976a914418327e3f3dda4cf5b9089325a4b95abdfa0334088ac722c0c00000000001976a914ba35042cfe9fc66fd35ac2224eebdafd1028ad2788acdc4ace020000000017a91474d691da1574e6b3c192ecfb52cc8984ee7b6c568700000000")
	txObj, _ := tx.ParseTx(txBytes, false)

	ok, err := txObj.VerifyInput(0)
	fmt.Println(ok, err)
}
