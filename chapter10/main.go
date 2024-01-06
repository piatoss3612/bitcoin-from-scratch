package main

import (
	"bytes"
	"chapter10/network"
	"encoding/hex"
	"fmt"
)

func main() {
	practice4()
}

func practice2() {
	msg := "f9beb4d976657261636b000000000000000000005df6e0e2"
	rawMsg, _ := hex.DecodeString(msg)

	fmt.Println(bytes.Equal(rawMsg[:4], network.NetworkMagic))
	fmt.Println(bytes.Equal(rawMsg[:4], network.TestNetworkMagic))
}

func practice3() {
	msg := "f9beb4d976657261636b000000000000000000005df6e0e2"
	rawMsg, _ := hex.DecodeString(msg)

	envelope, _ := network.ParseNetworkEnvelope(rawMsg)

	rawMsg2, _ := envelope.Serialize()

	fmt.Println(bytes.Equal(rawMsg, rawMsg2))
}

func practice4() {
	node, err := network.NewSimpleNode("localhost", 18555, network.SimNet, true)
	if err != nil {
		panic(err)
	}

	defer node.Close()

	fmt.Println("Connected to", node.Host, "on port", node.Port)

	resp, err := node.HandShake()
	if err != nil {
		panic(err)
	}

	if ok := <-resp; !ok {
		panic("Handshake failed")
	}

	fmt.Println("Handshake successful")
}
