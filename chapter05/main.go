package main

import (
	"chapter05/tx"
	"chapter05/utils"
	"encoding/hex"
	"fmt"
)

func main() {
	b1 := utils.EncodeVarint(255)
	fmt.Println(hex.EncodeToString(b1))
	fmt.Println(utils.ReadVarint(b1))

	b2 := utils.EncodeVarint(555)
	fmt.Println(hex.EncodeToString(b2))
	fmt.Println(utils.ReadVarint(b2))

	b3 := utils.EncodeVarint(70015)
	fmt.Println(hex.EncodeToString(b3))
	fmt.Println(utils.ReadVarint(b3))

	b4 := utils.EncodeVarint(18005558675309)
	fmt.Println(hex.EncodeToString(b4))
	fmt.Println(utils.ReadVarint(b4))

	tf := tx.NewTxFetcher()

	tx1, err := tf.Fetch("c05fd9e2a85a716e2c2679052cc24839f91bb7df510169aa54ed02a6c67dae9e", false, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(tx1)
}
