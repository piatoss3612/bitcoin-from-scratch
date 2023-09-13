package main

import (
	"chapter06/tx"
	"fmt"
)

func main() {
	// // 서명해시
	// z, err := hex.DecodeString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	// if err != nil {
	// 	panic(err)
	// }

	// // 공개키 sec 직렬화 값
	// sec, err := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	// if err != nil {
	// 	panic(err)
	// }

	// // 서명 der 직렬화 값
	// sig, err := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6")
	// if err != nil {
	// 	panic(err)
	// }

	// scriptPubKey := script.New(sec, 0xac) // 공개키와 연산자 OP_CHECKSIG로 잠금 스크립트 생성

	// scriptSig := script.New(sig) // 서명으로 해제 스크립트 생성

	// combined := scriptSig.Add(scriptPubKey) // 잠금 스크립트와 해제 스크립트 결합

	// fmt.Println("valid?", combined.Evaluate(z)) // 결합한 스크립트를 평가한 결과 출력

	// scriptPubKey2 := script.New(0x6e, 0x87, 0x91, 0x69, 0xa7, 0x7c, 0xa7, 0x87)

	// c1 := "255044462d312e330a25e2e3cfd30a0a0a312030206f626a0a3c3c2f57696474682032203020522f4865696768742033203020522f547970652034203020522f537562747970652035203020522f46696c7465722036203020522f436f6c6f7253706163652037203020522f4c656e6774682038203020522f42697473506572436f6d706f6e656e7420383e3e0a73747265616d0affd8fffe00245348412d3120697320646561642121212121852fec092339759c39b1a1c63c4c97e1fffe017f46dc93a6b67e013b029aaa1db2560b45ca67d688c7f84b8c4c791fe02b3df614f86db1690901c56b45c1530afedfb76038e972722fe7ad728f0e4904e046c230570fe9d41398abe12ef5bc942be33542a4802d98b5d70f2a332ec37fac3514e74ddc0f2cc1a874cd0c78305a21566461309789606bd0bf3f98cda8044629a1"
	// c2 := "255044462d312e330a25e2e3cfd30a0a0a312030206f626a0a3c3c2f57696474682032203020522f4865696768742033203020522f547970652034203020522f537562747970652035203020522f46696c7465722036203020522f436f6c6f7253706163652037203020522f4c656e6774682038203020522f42697473506572436f6d706f6e656e7420383e3e0a73747265616d0affd8fffe00245348412d3120697320646561642121212121852fec092339759c39b1a1c63c4c97e1fffe017346dc9166b67e118f029ab621b2560ff9ca67cca8c7f85ba84c79030c2b3de218f86db3a90901d5df45c14f26fedfb3dc38e96ac22fe7bd728f0e45bce046d23c570feb141398bb552ef5a0a82be331fea48037b8b5d71f0e332edf93ac3500eb4ddc0decc1a864790c782c76215660dd309791d06bd0af3f98cda4bc4629b1"

	// collision1, err := hex.DecodeString(c1)
	// if err != nil {
	// 	panic(err)
	// }

	// collision2, err := hex.DecodeString(c2)
	// if err != nil {
	// 	panic(err)
	// }

	// scriptSig2 := script.New(collision1, collision2)

	// combined2 := scriptSig2.Add(scriptPubKey2)

	// fmt.Println("valid?", combined2.Evaluate([]byte{0x00}))

	// rawTx := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

	// hexTx, err := hex.DecodeString(rawTx)
	// if err != nil {
	// 	panic(err)
	// }

	// tx1, err := tx.ParseTx(hexTx, false)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(tx1)

	// fmt.Println(tx1.Inputs[0].ScriptSig)

	fetcher := tx.NewTxFetcher()

	tx1, err := fetcher.Fetch("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81", false, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(tx1)
}
