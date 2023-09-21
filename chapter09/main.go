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
	calcDifficulty()
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