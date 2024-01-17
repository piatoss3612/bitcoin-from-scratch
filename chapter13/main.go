package main

import (
	"chapter13/block"
	"chapter13/bloomfilter"
	"chapter13/merkleblock"
	"chapter13/network"
	"chapter13/tx"
	"chapter13/utils"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

func main() {
	// VerifyBlockHeaders()
	CheckIfBloomfilterWorks()
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

func CheckIfBloomfilterWorks() {
	startBlock, _ := hex.DecodeString("00000000000538d5c2246336644f9a4956551afb44ba47278759ec55ea912e19")

	address := "mwJn1YPMq7y5F8J3LkC5Hxg9PHyZ5K4cFv"
	h160, _ := utils.DecodeBase58(address)

	node, err := network.NewSimpleNode("84.250.85.135", 18333, network.TestNet, false) // testnet
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	log.Println("Connected to", node.Host, "on port", node.Port)

	bf := bloomfilter.New(30, 5, 90210)
	bf.Add(h160)

	resp, err := node.HandShake()
	if err != nil {
		log.Fatal(err)
	}

	if ok := <-resp; !ok {
		log.Fatal("Handshake failed")
	}

	time.Sleep(1 * time.Second)

	log.Println("Sending filterload message")

	if err := node.Send(bf.Filterload()); err != nil {
		log.Fatal(err)
	}

	log.Println("Sending getheaders message")

	getheaders := network.DefaultGetHeadersMessage()
	getheaders.StartBlock = startBlock

	if err := node.Send(getheaders); err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	envelopes, errs := node.WaitFor([]network.Command{network.HeadersCommand}, done)
	getdata := network.NewGetDataMessage()

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
			case headersEnvelope := <-envelopes:
				if headersEnvelope == nil {
					continue
				}

				headers, err := network.ParseHeadersMessage(headersEnvelope.Payload)
				if err != nil {
					log.Fatal(err)
				}

				for _, header := range headers.Headers {
					ok, err := header.CheckProofOfWork()
					if err != nil {
						log.Fatal(err)
					}

					if !ok {
						log.Fatal("Block does not satisfy proof of work")
					}

					hash, err := header.Hash()
					if err != nil {
						log.Fatal(err)
					}

					getdata.AddData(network.FiltedBlockDataItem, hash)
				}

				return
			}
		}
	}()

	<-done

	time.Sleep(1 * time.Second)

	log.Println("Sending getdata message")

	if err := node.Send(getdata); err != nil {
		log.Fatal(err)
	}

	done = make(chan struct{})

	envelopes, errs = node.WaitFor([]network.Command{network.MerkleBlockCommand, network.TxCommand}, done)

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
			case envelope := <-envelopes:
				switch envelope.Command.String() {
				case network.MerkleBlockCommand.String():
					mb := merkleblock.MerkleBlock{}
					err := mb.Parse(envelope.Payload)
					if err != nil {
						log.Fatalf("Error parsing merkle block: %s", err)
					}

					ok, err := mb.IsValid()
					if err != nil {
						log.Fatalf("Error validating merkle block: %s", err)
					}

					if !ok {
						log.Fatal("Merkle block is not valid")
					}
				case network.TxCommand.String():
					tx, err := tx.ParseTx(envelope.Payload)
					if err != nil {
						log.Fatalf("Error parsing tx: %s", err)
					}

					for i, out := range tx.Outputs {
						if strings.EqualFold(out.ScriptPubKey.Address(true), address) {
							id, _ := tx.ID()
							fmt.Printf("Found matching tx %s at output index %d\n", id, i)

							return
						}
					}
				}
			}
		}
	}()

	<-done

	log.Println("Done")
}
