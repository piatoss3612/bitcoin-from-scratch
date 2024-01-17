package main

import (
	"bytes"
	"chapter13/block"
	"chapter13/bloomfilter"
	"chapter13/ecc"
	"chapter13/merkleblock"
	"chapter13/network"
	"chapter13/script"
	"chapter13/tx"
	"chapter13/utils"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// VerifyBlockHeaders()
	// CheckIfBloomfilterWorks()
	// GenTestnetTx()
	Practice6()
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

func GenTestnetTx() {
	secret1 := utils.LittleEndianToBigInt(utils.Hash256(utils.StringToBytes("piatoss rules the world"))) // 개인 키 생성
	privateKey1, _ := ecc.NewS256PrivateKey(secret1.Bytes())

	secret2 := utils.LittleEndianToBigInt(utils.Hash256(utils.StringToBytes("piatoss ruins the world"))) // 개인 키 생성
	privateKey2, _ := ecc.NewS256PrivateKey(secret2.Bytes())

	address1 := privateKey1.Point().Address(true, true) // 비트코인을 보내는 주소
	address2 := privateKey2.Point().Address(true, true) // 비트코인을 받는 주소

	prevTx := "e770e0b481166da7d0d139c855e86633a12dbd4fa9b97f33a31fc9a458f8ddd7" // 이전 트랜잭션 ID (faucet -> address1)
	prevIndex := 0                                                               // 이전 트랜잭션의 출력 인덱스
	txIn := tx.NewTxIn(prevTx, prevIndex, nil)

	balance := 1193538 // 잔고

	changeAmount := balance - (balance * 7 / 10)            // 잔액
	changeH160, _ := utils.DecodeBase58(address1)           // 잔액을 받을 주소
	changeScript := script.NewP2pkhScript(changeH160)       // p2pkh 잠금 스크립트 생성
	changeOutput := tx.NewTxOut(changeAmount, changeScript) // 트랜잭션 출력 생성

	targetAmount := balance * 6 / 10                        // 사용할 금액
	targetH160, _ := utils.DecodeBase58(address2)           // 받을 주소
	targetScript := script.NewP2pkhScript(targetH160)       // p2pkh 잠금 스크립트 생성
	targetOutput := tx.NewTxOut(targetAmount, targetScript) // 트랜잭션 출력 생성

	txObj := tx.NewTx(1, []*tx.TxIn{txIn}, []*tx.TxOut{changeOutput, targetOutput}, 0, true, false) // address1 -> address2로 비트코인 전송 트랜잭션 생성

	ok, err := txObj.SignInput(0, privateKey1, true) // 서명 생성
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Signature valid?", ok)

	serializedTx, err := txObj.Serialize() // 트랜잭션 직렬화
	if err != nil {
		log.Fatal(err)
	}

	hexTx := hex.EncodeToString(serializedTx)

	log.Println("Sending transaction:", hexTx)

	body := bytes.NewBuffer([]byte(hexTx))

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodPost, "https://blockstream.info/testnet/api/tx", body) // testnet 블록체인에 트랜잭션 전송
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body) // 응답 바디 읽기
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(respBody)) // 트랜잭션 ID 출력 - 4036b355ee1274384e397b58c3d3caff755acb20a8747c340cda29f5bd03ef24
}

func Practice6() {
	secret1 := utils.LittleEndianToBigInt(utils.Hash256(utils.StringToBytes("piatoss rules the world"))) // 개인 키 생성
	privateKey1, _ := ecc.NewS256PrivateKey(secret1.Bytes())

	secret2 := utils.LittleEndianToBigInt(utils.Hash256(utils.StringToBytes("piatoss ruins the world"))) // 개인 키 생성
	privateKey2, _ := ecc.NewS256PrivateKey(secret2.Bytes())

	address1 := privateKey1.Point().Address(true, true) // 비트코인을 보내는 주소
	address2 := privateKey2.Point().Address(true, true) // 비트코인을 받는 주소

	h160, _ := utils.DecodeBase58(address1) // address1의 공개 키 해시

	node, err := network.NewSimpleNode("84.250.85.135", 18333, network.TestNet, false) // testnet 노드에 연결
	if err != nil {
		log.Fatal(err)
	}
	defer node.Close()

	log.Println("Connected to", node.Host, "on port", node.Port)

	bf := bloomfilter.New(30, 5, 90210) // 블룸 필터 생성
	bf.Add(h160)                        // 블룸 필터에 공개 키 해시 추가

	resp, err := node.HandShake() // 핸드셰이크
	if err != nil {
		log.Fatal(err)
	}

	if ok := <-resp; !ok { // 핸드셰이크 실패 시 종료
		log.Fatal("Handshake failed")
	}

	time.Sleep(1 * time.Second)

	log.Println("Sending filterload message")

	if err := node.Send(bf.Filterload()); err != nil { // 블룸 필터로 필터로드 메시지 전송
		log.Fatal(err)
	}

	log.Println("Sending getheaders message")

	getheaders := network.DefaultGetHeadersMessage()                                                      // getheaders 메시지 생성
	startBlock, _ := hex.DecodeString("0000000000000014fbf5791d8333ef9cea0c87aff27f98e102dbcf40963aea8b") // 2573311번 블록 (address1 -> address2로 비트코인 전송한 블록은 2573313번 블록)
	getheaders.StartBlock = startBlock                                                                    // getheaders 메시지에 start_block 필드 설정

	if err := node.Send(getheaders); err != nil { // getheaders 메시지 전송
		log.Fatal(err)
	}

	done := make(chan struct{})

	envelopes, errs := node.WaitFor([]network.Command{network.HeadersCommand}, done) // headers 메시지 대기

	getdata := network.NewGetDataMessage() // getdata 메시지 생성

	go func() {
		defer close(done)

		for {
			select {
			case err := <-errs:
				if err == io.EOF { // 에러가 EOF일 경우 종료
					log.Println("Connection closed")
					return
				}
				log.Fatalf("Error receiving message: %s", err)
			case headersEnvelope := <-envelopes:
				if headersEnvelope == nil {
					continue
				}

				headers, err := network.ParseHeadersMessage(headersEnvelope.Payload) // headers 메시지 파싱
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Received headers message with %d headers\n", len(headers.Headers))

				// 블록 헤더 검증
				for _, header := range headers.Headers {
					ok, err := header.CheckProofOfWork() // 작업 증명 검증
					if err != nil {
						log.Fatal(err)
					}

					if !ok {
						log.Fatal("Block does not satisfy proof of work")
					}

					hash, err := header.Hash() // 블록 해시
					if err != nil {
						log.Fatal(err)
					}

					getdata.AddData(network.FiltedBlockDataItem, hash) // getdata 메시지에 필터링된 블록 데이터 추가
				}

				return
			}
		}
	}()

	<-done

	time.Sleep(1 * time.Second)

	log.Println("Sending getdata message")

	if err := node.Send(getdata); err != nil { // getdata 메시지 전송
		log.Fatal(err)
	}

	done = make(chan struct{})

	envelopes, errs = node.WaitFor([]network.Command{network.MerkleBlockCommand, network.TxCommand}, done) // merkleblock, tx 메시지 대기

	go func() {
		defer close(done)

		for {
			select {
			case err := <-errs:
				if err == io.EOF { // 에러가 EOF일 경우 종료
					log.Println("Connection closed")
					return
				}
				log.Fatalf("Error receiving message: %s", err)
			case envelope := <-envelopes:
				switch envelope.Command.String() {
				case network.MerkleBlockCommand.String():
					mb := merkleblock.MerkleBlock{}
					err := mb.Parse(envelope.Payload) // merkleblock 메시지 파싱
					if err != nil {
						log.Fatalf("Error parsing merkle block: %s", err)
					}

					ok, err := mb.IsValid() // merkleblock 메시지 검증
					if err != nil {
						log.Fatalf("Error validating merkle block: %s", err)
					}

					if !ok {
						log.Fatal("Merkle block is not valid")
					}
				case network.TxCommand.String():
					transaction, err := tx.ParseTx(envelope.Payload) // tx 메시지 파싱
					if err != nil {
						log.Fatalf("Error parsing tx: %s", err)
					}

					for i, out := range transaction.Outputs { // tx 메시지의 출력 검색
						if strings.EqualFold(out.ScriptPubKey.Address(true), address1) {
							prevTx, _ := transaction.ID()

							log.Printf("Found matching tx %s at output index %d\n", prevTx, i)

							prevIndex := 0                             // 이전 트랜잭션의 출력 인덱스
							txIn := tx.NewTxIn(prevTx, prevIndex, nil) // 이전 트랜잭션의 출력을 참조하는 트랜잭션 입력 생성

							balance := 358062 // 잔고 (하드코딩)

							changeAmount := balance - (balance * 7 / 10)            // 잔액
							changeH160, _ := utils.DecodeBase58(address1)           // 잔액을 받을 주소
							changeScript := script.NewP2pkhScript(changeH160)       // p2pkh 잠금 스크립트 생성
							changeOutput := tx.NewTxOut(changeAmount, changeScript) // 트랜잭션 출력 생성

							targetAmount := balance * 6 / 10                        // 사용할 금액
							targetH160, _ := utils.DecodeBase58(address2)           // 받을 주소
							targetScript := script.NewP2pkhScript(targetH160)       // p2pkh 잠금 스크립트 생성
							targetOutput := tx.NewTxOut(targetAmount, targetScript) // 트랜잭션 출력 생성

							txObj := tx.NewTx(1, []*tx.TxIn{txIn}, []*tx.TxOut{changeOutput, targetOutput}, 0, true, false) // address1 -> address2로 비트코인 전송 트랜잭션 생성

							ok, err := txObj.SignInput(0, privateKey1, true) // 서명 생성
							if err != nil {
								log.Fatal(err)
							}

							log.Println("Signature valid?", ok)

							serializedTx, err := txObj.Serialize() // 트랜잭션 직렬화
							if err != nil {
								log.Fatal(err)
							}

							hexTx := hex.EncodeToString(serializedTx)

							log.Println("Sending transaction:", hexTx)

							body := bytes.NewBuffer([]byte(hexTx))

							client := http.DefaultClient

							req, err := http.NewRequest(http.MethodPost, "https://blockstream.info/testnet/api/tx", body) // testnet 블록체인에 트랜잭션 전송
							if err != nil {
								log.Fatal(err)
							}

							resp, err := client.Do(req)
							if err != nil {
								log.Fatal(err)
							}
							defer resp.Body.Close()

							respBody, err := io.ReadAll(resp.Body) // 응답 바디 읽기
							if err != nil {
								log.Fatal(err)
							}

							fmt.Println("Transaction ID:", string(respBody)) // 트랜잭션 ID 출력

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
