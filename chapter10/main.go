package main

import (
	"bytes"
	"chapter10/block"
	"chapter10/network"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

func main() {
	practice5()
	// practice4()
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
	node, err := network.NewSimpleNode("71.13.92.62", 18333, network.TestNet, true)
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

func practice5() {
	rawGenesisBlock, _ := hex.DecodeString("0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff001d1aa4ae18")
	genesisBlock, _ := block.Parse(rawGenesisBlock)

	node, err := network.NewSimpleNode("71.13.92.62", 18333, network.TestNet, true)
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

	time.Sleep(1 * time.Second)

	getheaders := network.DefaultGetHeadersMessage()
	genesisHash, _ := genesisBlock.Hash()
	getheaders.StartBlock = genesisHash

	err = node.Send(getheaders, network.TestNet)
	if err != nil {
		panic(err)
	}

	envelopes, errs := node.WaitFor([]network.Command{network.HeadersCommand})

	for {
		select {
		case err := <-errs:
			if err == io.EOF {
				fmt.Println("Connection closed")
				return
			}
			panic(err)
		case headers := <-envelopes:
			msg, err := network.ParseHeadersMessage(headers.Payload)
			if err != nil {
				panic(err)
			}

			fmt.Println("Received headers message with", len(msg.Headers), "headers")

			for _, header := range msg.Headers {
				hash, err := header.Hash()
				if err != nil {
					panic(err)
				}

				fmt.Println("Header hash:", hex.EncodeToString(hash[:]))
			}
		}
	}
}
