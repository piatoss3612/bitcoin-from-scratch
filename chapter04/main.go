package main

import (
	"chapter04/ecc"
	"fmt"
)

func main() {
	passphrase := "piatoss will be the blockchain core developer."
	secret := ecc.LittleEndianToBigInt(ecc.Hash256([]byte(passphrase)))

	// 개인키 생성
	privateKey, err := ecc.NewS256PrivateKey(secret.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println(privateKey.Point().Address(true, true))
}
