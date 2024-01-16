package main

import (
	"chapter13/block"
	"chapter13/network"
	"chapter13/utils"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func main() {
	VerifyBlockHeaders()

	// rawCmd, _ := hex.DecodeString("66656566696c746572")
	// cmd := network.ParseCommand(rawCmd)

	// fmt.Println(cmd)
}

func VerifyBlockHeaders() {
	rawGenesisBlock, _ := hex.DecodeString("0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4adae5494dffff001d1aa4ae18") // testnet
	// rawGenesisBlock, _ := hex.DecodeString("0100000000000000000000000000000000000000000000000000000000000000000000003ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a29ab5f49ffff001d1dac2b7c") // mainnet
	previous, _ := block.Parse(rawGenesisBlock)

	node, err := network.NewSimpleNode("71.13.92.62", 18333, network.TestNet, true) // testnet
	// node, err := network.NewSimpleNode("165.227.133.233", 8333, network.MainNet, true) // mainnet
	if err != nil {
		log.Fatal(err)
	}

	defer node.Close()

	log.Println("Connected to", node.Host, "on port", node.Port)

	resp, err := node.HandShake()
	if err != nil {
		log.Fatal(err)
	}

	if ok := <-resp; !ok {
		log.Fatal("Handshake failed")
	}

	time.Sleep(1 * time.Second)

	getheaders := network.DefaultGetHeadersMessage()
	hexPreviousHash, _ := previous.Hash()
	previousHash := hex.EncodeToString(hexPreviousHash)
	firstEpochTimestamp := previous.Timestamp
	expectedBits, _ := hex.DecodeString("ffff001d")

	getheaders.StartBlock = hexPreviousHash

	err = node.Send(getheaders, network.TestNet)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	envelopes, errs := node.WaitFor([]network.Command{network.HeadersCommand}, done)

	go func() {
		defer close(done)

		for {
			select {
			case err := <-errs:
				if err == io.EOF {
					log.Println("Connection closed")
					return
				}
				log.Fatalf("Error receiving message: %s", err)
			case headers := <-envelopes:
				msg, err := network.ParseHeadersMessage(headers.Payload)
				if err != nil {
					log.Fatalf("Error parsing headers message: %s", err)
				}

				fmt.Println("Received headers message with", len(msg.Headers), "headers")

				count := 1

				for i, header := range msg.Headers {
					// 작업 증명 검증
					ok, err := header.CheckProofOfWork()
					if err != nil {
						log.Fatalf("Error checking proof of work for block header %d: %s", i, err)
					}

					if !ok {
						log.Fatalf("Block header %d does not satisfy proof of work", i)
					}

					// 이전 블록 해시 검증
					if !strings.EqualFold(header.PrevBlock, previousHash) {
						log.Fatalf("Block header %d's previous hash is not correct", i)
					}

					// 블록 난이도 계산
					if count%2016 == 0 {
						timeDiff := previous.Timestamp - firstEpochTimestamp
						expectedBits = block.CalculateNewBits(utils.IntToBytes(previous.Bits, 4), int64(timeDiff))
						firstEpochTimestamp = previous.Timestamp
						fmt.Println("New epoch, expected bits are", hex.EncodeToString(expectedBits))
					}

					// 블록 난이도 검증
					if header.Bits != utils.BytesToInt(expectedBits) {
						log.Fatalf("Block header %d's bits are not correct: expected %d, got %d", i, utils.BytesToInt(expectedBits), header.Bits)
					}

					hexHash, err := header.Hash()
					if err != nil {
						log.Fatalf("Error hashing block header %d: %s", i, err)
					}

					hash := hex.EncodeToString(hexHash)

					previousHash = hash

					count++
				}

				return
			}
		}
	}()

	<-done

	fmt.Println("Done")
}
