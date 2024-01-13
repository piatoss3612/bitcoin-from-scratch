package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

/*
	p2wpkh 이후의 스크립트를 가진 트랜잭션들은 왜 검증이 안될까?

	후보 1. 스크립트 Evaluate 메서드에 문제가 있다.
	후보 2. DER 서명 파싱에 문제가 있다. (문제 없음 확인)
	후보 3. 서명을 검증하는 과정에 문제가 있다. (확인 필요)

	이 세 개 말고는 딱히 없어보이는데...
	Evaluate 메서드에서 오류는 발생하지 않고 서명 검증에서도 오류는 발생하지 않지만
	서명 검증 결과가 false로 나온다.
	내가 모르는 매콤한 무언가가 있는 건가?
*/

func main() {
	test1()
	// test2()
}

func test1() {
	rawPubkey, _ := hex.DecodeString("038262a6c6cec93c2d3ecd6c6072efea86d02ff8e3328bbd0242b20af3425990ac")

	pubkey, err := btcec.ParsePubKey(rawPubkey)
	if err != nil {
		panic(err)
	}

	rawSig, _ := hex.DecodeString("3045022100df7b7e5cda14ddf91290e02ea10786e03eb11ee36ec02dd862fe9a326bbcb7fd02203f5b4496b667e6e281cc654a2da9e4f08660c620a1051337fa8965f727eb191901")

	sig, err := ecdsa.ParseDERSignature(rawSig)
	if err != nil {
		panic(err)
	}

	// z, _ := hex.DecodeString("34cd3ab4ce6e9f993a419042c75365f567019385750260253b4b1b703dcfe449") // 내가 구한 z
	z, _ := hex.DecodeString("645a6a03b756ef2adca77e00905810f407876b44ed2b9ea5f6383a1da0b339f0") // 실제 z -> true 나옴;;

	verified := sig.Verify(z, pubkey)

	println(verified) // false -> z가 잘못된 것 같다. 그렇다면 SigHashBip143를 잘못 구현한 것인가? -> 맞다. 잘못 구현했다.
}

func test2() {
	// Decode hex-encoded serialized public key.
	pubKeyBytes, err := hex.DecodeString("02a673638cb9587cb68ea08dbef685c" +
		"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5")
	if err != nil {
		fmt.Println(err)
		return
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode hex-encoded serialized signature.
	sigBytes, err := hex.DecodeString("30450220090ebfb3690a0ff115bb1b38b" +
		"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
		"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479")

	if err != nil {
		fmt.Println(err)
		return
	}
	signature, err := ecdsa.ParseSignature(sigBytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Verify the signature for the message using the public key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)
	fmt.Println("Signature Verified?", verified)

	// Output:
	// Signature Verified? true
}
