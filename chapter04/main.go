package main

import (
	"chapter04/ecc"
	"fmt"
	"math/big"
)

func main() {
	r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)

	sig := ecc.NewS256Signature(r, s)

	der := sig.DER()

	fmt.Printf("DER: %x\n", der)
}
