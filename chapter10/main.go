package main

import (
	"bytes"
	"chapter10/network"
	"encoding/hex"
	"fmt"
	"io"
	"net"
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
	conn, err := net.Dial("tcp", "localhost:18555")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr())

	msg := network.DefaultVersionMessage()

	msgBytes, _ := msg.Serialize()

	envelope, _ := network.New(msg.Command, msgBytes, network.SimNet)

	envelopeBytes, _ := envelope.Serialize()

	fmt.Println("Send:", hex.EncodeToString(envelopeBytes))

	n, err := conn.Write(envelopeBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println("Sent", n, "bytes")

	for {
		buf := make([]byte, 1024)

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
		}

		fmt.Println(hex.EncodeToString(buf[:n]))

		decodedEnvelope, err := network.ParseNetworkEnvelope(buf[:n])
		if err != nil {
			panic(err)
		}

		fmt.Println(decodedEnvelope)
	}
}
