package main

import (
	"chapter06/script"
	"encoding/hex"
	"fmt"
)

func main() {
	z, err := hex.DecodeString("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
	if err != nil {
		panic(err)
	}

	sec, err := hex.DecodeString("04887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	if err != nil {
		panic(err)
	}

	sig, err := hex.DecodeString("3045022000eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c022100c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6")
	if err != nil {
		panic(err)
	}

	scriptPubKey := script.New(sec, 0xac)

	scriptSig := script.New(sig)

	combined := scriptSig.Add(scriptPubKey)

	fmt.Println("valid?", combined.Evaluate(z))
}
