package main

import (
	"bytes"
	"chapter10/network"
	"encoding/hex"
	"fmt"
)

func main() {

}

func parseNetworkEnvelop() {
	msg := "f9beb4d976657261636b000000000000000000005df6e0e2"
	rawMsg, _ := hex.DecodeString(msg)

	fmt.Println(bytes.Equal(rawMsg[:4], network.NetworkMagic))
	fmt.Println(bytes.Equal(rawMsg[:4], network.TestNetworkMagic))

	envelope, err := network.ParseNetworkEnvelope(rawMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println(envelope)

	rawMsg2 := envelope.Serialize()

	fmt.Println(bytes.Equal(rawMsg, rawMsg2))
	fmt.Println(hex.EncodeToString(rawMsg2))
}
