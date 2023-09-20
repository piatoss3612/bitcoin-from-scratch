package main

import (
	"chapter08/utils"
	"encoding/hex"
	"fmt"
)

func main() {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	fmt.Println(utils.H160ToP2shAddress(h160, false))
}
