package main

import (
	"chapter05/helper"
	"encoding/hex"
	"fmt"
)

func main() {
	b1 := helper.EncodeVarint(255)
	fmt.Println(hex.EncodeToString(b1))
	fmt.Println(helper.ReadVarint(b1))

	b2 := helper.EncodeVarint(555)
	fmt.Println(hex.EncodeToString(b2))
	fmt.Println(helper.ReadVarint(b2))

	b3 := helper.EncodeVarint(70015)
	fmt.Println(hex.EncodeToString(b3))
	fmt.Println(helper.ReadVarint(b3))

	b4 := helper.EncodeVarint(18005558675309)
	fmt.Println(hex.EncodeToString(b4))
	fmt.Println(helper.ReadVarint(b4))
}
