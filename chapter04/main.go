package main

import (
	"chapter04/ecc"
	"fmt"
	"math/big"
)

func main() {
	a, _ := big.NewInt(0).SetString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 16)
	b, _ := big.NewInt(0).SetString("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 16)
	c, _ := big.NewInt(0).SetString("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 16)

	fmt.Println(ecc.EncodeBase58(a.Bytes()))
	fmt.Println(ecc.EncodeBase58(b.Bytes()))
	fmt.Println(ecc.EncodeBase58(c.Bytes()))
}
