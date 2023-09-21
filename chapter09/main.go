package main

import (
	"chapter09/block"
	"chapter09/script"
	"chapter09/utils"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	// readCoinbaseTxScriptSig()
	// parseHeightFromCoinbaseTxScriptSig()
	// readBlockID()
	// readBlockVersionBIP9()
	// calcTargetFromBits()
	// calcDifficulty()
	// calcNewTarget()
	calcNewTargetAndConvertToBits()
}

func readCoinbaseTxScriptSig() {
	rawScript, _ := hex.DecodeString("4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73")
	scriptSig, _, _ := script.Parse(rawScript)

	content, ok := scriptSig.Cmds[2].([]byte)
	if !ok {
		fmt.Println("Fail to parse scriptSig")
		return
	}

	fmt.Println(string(content)) // The Times 03/Jan/2009 Chancellor on brink of second bailout for banks
}

func parseHeightFromCoinbaseTxScriptSig() {
	rawScript, _ := hex.DecodeString("5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00")
	scriptSig, _, _ := script.Parse(rawScript)

	fmt.Println(scriptSig)

	height, ok := scriptSig.Cmds[0].([]byte)
	if !ok {
		fmt.Println("Fail to parse scriptSig")
		return
	}
	fmt.Println(utils.LittleEndianToInt(height)) // 465879
}

func readBlockID() {
	rawBlockHeader, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	//blockHash := utils.Hash256(rawBlockHeader)

	//blockID := hex.EncodeToString(utils.ReverseBytes(blockHash))

	parsed, _ := block.Parse(rawBlockHeader)

	blockHash, _ := parsed.Hash()
	blockID := hex.EncodeToString(utils.ReverseBytes(blockHash))

	fmt.Println(blockID) // 0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523
}

func readBlockVersionBIP9() {
	rawBlockHeader, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")

	b, _ := block.Parse(rawBlockHeader)

	fmt.Println("BIP9:", b.Version>>29 == 0x001) // 처음 3비트가 001이면 BIP9 활성화
	fmt.Println("BIP91:", b.Version>>4&1 == 1)   // 4번째 비트가 1이면 BIP91 활성화
	fmt.Println("BIP141:", b.Version>>1&1 == 1)  // 2번째 비트가 1이면 BIP141 활성화

	fmt.Println("BIP9:", b.Bip9())     // 처음 3비트가 001이면 BIP9 활성화
	fmt.Println("BIP91:", b.Bip91())   // 4번째 비트가 1이면 BIP91 활성화
	fmt.Println("BIP141:", b.Bip141()) // 2번째 비트가 1이면 BIP141 활성화
}

func calcTargetFromBits() {
	bits, _ := hex.DecodeString("e93c0118")
	// exp := big.NewInt(0).SetBytes([]byte{bits[len(bits)-1]}) // 지수
	// coef := utils.LittleEndianToBigInt(bits[:len(bits)-1])   // 계수

	// target := big.NewInt(0).Mul(coef, big.NewInt(0).Exp(big.NewInt(256), big.NewInt(0).Sub(exp, big.NewInt(3)), nil)) // 계수 * 256^(지수-3) = 목푯값

	// fmt.Println(hex.EncodeToString(target.FillBytes(make([]byte, 32)))) // 0000000000000000013ce9000000000000000000000000000000000000000000

	target := block.BitsToTarget(bits)

	rawBlockHeader, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	proof := utils.LittleEndianToBigInt(utils.Hash256(rawBlockHeader))

	fmt.Println(proof.Cmp(target) < 0) // proof가 target보다 작으면 Cmp는 -1을 반환
}

func calcDifficulty() {
	bits, _ := hex.DecodeString("e93c0118")
	target := block.BitsToTarget(bits)

	// difficulty = 0xffff * 256^(0x1d - 3) / target
	difficulty := big.NewFloat(0).Mul(big.NewFloat(0xffff), big.NewFloat(0).Quo(big.NewFloat(0).SetInt(new(big.Int).Exp(big.NewInt(256), big.NewInt(0).Sub(big.NewInt(0x1d), big.NewInt(3)), nil)), big.NewFloat(0).SetInt(target))) // 0xffff * 256^(0x1d - 3) / target

	fmt.Println(difficulty.Text('f', -1)) // 888171856257.3206
}

func calcNewTarget() {
	rawLastBlock, _ := hex.DecodeString("00000020fdf740b0e49cf75bb3d5168fb3586f7613dcc5cd89675b0100000000000000002e37b144c0baced07eb7e7b64da916cd3121f2427005551aeb0ec6a6402ac7d7f0e4235954d801187f5da9f5")
	rawFirstBlock, _ := hex.DecodeString("000000201ecd89664fd205a37566e694269ed76e425803003628ab010000000000000000bfcade29d080d9aae8fd461254b041805ae442749f2a40100440fc0e3d5868e55019345954d80118a1721b2e")

	lastBlock, _ := block.Parse(rawLastBlock)
	firstBlock, _ := block.Parse(rawFirstBlock)

	timeDiff := lastBlock.Timestamp - firstBlock.Timestamp

	twoWeek := 60 * 60 * 24 * 14

	if timeDiff > twoWeek*4 {
		timeDiff = twoWeek * 4
	} else if timeDiff < twoWeek/4 {
		timeDiff = twoWeek / 4
	}

	newTarget := big.NewInt(0).Div(big.NewInt(0).Mul(lastBlock.Target(), big.NewInt(int64(timeDiff))), big.NewInt(int64(twoWeek))).FillBytes(make([]byte, 32))

	fmt.Println(hex.EncodeToString(newTarget))
}

func calcNewTargetAndConvertToBits() {
	rawFirstBlock, _ := hex.DecodeString("02000020f1472d9db4b563c35f97c428ac903f23b7fc055d1cfc26000000000000000000b3f449fcbe1bc4cfbcb8283a0d2c037f961a3fdf2b8bedc144973735eea707e1264258597e8b0118e5f00474")
	rawLastBlock, _ := hex.DecodeString("000000203471101bbda3fe307664b3283a9ef0e97d9a38a7eacd8800000000000000000010c8aba8479bbaa5e0848152fd3c2289ca50e1c3e58c9a4faaafbdf5803c5448ddb845597e8b0118e43a81d3")

	firstBlock, _ := block.Parse(rawFirstBlock)
	lastBlock, _ := block.Parse(rawLastBlock)

	timeDiff := lastBlock.Timestamp - firstBlock.Timestamp

	/*
		twoWeek := 60 * 60 * 24 * 14

		if timeDiff > twoWeek*4 {
			timeDiff = twoWeek * 4
		} else if timeDiff < twoWeek/4 {
			timeDiff = twoWeek / 4
		}

		newTarget := big.NewInt(0).Div(big.NewInt(0).Mul(lastBlock.Target(), big.NewInt(int64(timeDiff))), big.NewInt(int64(twoWeek)))

		newBits := block.TargetToBits(newTarget)
	*/

	newBits := block.CalculateNewBits(block.TargetToBits(lastBlock.Target()), int64(timeDiff))

	fmt.Println(hex.EncodeToString(newBits)) // 80df6217
}
