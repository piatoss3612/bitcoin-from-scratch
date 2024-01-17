package network

import (
	"bytes"
	"chapter13/utils"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type SimpleNode struct {
	Host    string
	Port    int
	Network NetworkType
	Logging bool

	conn            net.Conn
	envelopes       chan *NetworkEnvelope
	serverCloseChan chan struct{}
}

func NewSimpleNode(host string, port int, network NetworkType, logging bool) (*SimpleNode, error) {
	if port == 0 {
		switch network {
		case MainNet:
			port = DefaultMainNetPort
		case TestNet:
			port = DefaultTestNetPort
		case RegTest:
			port = DefaultRegTestPort
		case SimNet:
			port = DefaultSimNetPort
		default:
			return nil, ErrInvalidNetwork
		}
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	node := &SimpleNode{
		Host:            host,
		Port:            port,
		Network:         network,
		Logging:         logging,
		conn:            conn,
		envelopes:       make(chan *NetworkEnvelope, 1000),
		serverCloseChan: make(chan struct{}),
	}

	go node.readMessages()

	return node, nil
}

func (sn *SimpleNode) Close() error {
	close(sn.serverCloseChan)
	close(sn.envelopes)

	return sn.conn.Close()
}

func (sn *SimpleNode) Send(msg Message, network ...NetworkType) error {
	msgBytes, err := msg.Serialize()
	if err != nil {
		return err
	}

	var envelope *NetworkEnvelope

	if len(network) > 0 {
		envelope, err = NewEnvelope(msg.Command(), msgBytes, network[0])
	} else {
		envelope, err = NewEnvelope(msg.Command(), msgBytes, sn.Network)
	}

	if err != nil {
		return err
	}

	envelopeBytes, err := envelope.Serialize()
	if err != nil {
		return err
	}

	if sn.Logging {
		log.Printf("Send: %s\n", envelope.Command)
	}

	_, err = sn.conn.Write(envelopeBytes)
	return err
}

func (sn *SimpleNode) Read() (*NetworkEnvelope, error) {
	// 너무 작은 버퍼를 사용해서 데이터를 전부 읽어오지 못하는 문제가 있음 (책에서 파이썬으로 구현한 코드는 그런 문제가 없음)
	// 32MB 버퍼를 사용해도 안됨
	buf := make([]byte, 32*1024*1024)

	n, err := sn.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	envelope, _, err := ParseNetworkEnvelope(buf[:n])
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

func (sn *SimpleNode) readMessages() {
	mu := sync.Mutex{} // data 버퍼에 동시에 접근하는 것을 방지하기 위해 사용

	data := new(bytes.Buffer) // 읽어온 데이터를 저장하는 버퍼

	go func() {
		for {
			select {
			case <-sn.serverCloseChan:
				data.Reset()
				return
			default:
				mu.Lock()
				if data.Len() < 24 { // 4 + 12 + 4 + 4 (magic + command + payload length + checksum) (최소한의 길이를 만족하지 못하면 메시지를 읽어올 수 없음)
					mu.Unlock()
					time.Sleep(100 * time.Millisecond)
					continue
				}

				b := data.Bytes()                                  // 버퍼의 데이터를 읽어옴
				payloadLength := utils.LittleEndianToInt(b[16:20]) // 페이로드 길이를 읽어옴

				totalLength := 4 + 12 + 4 + payloadLength + 4 // magic + command + payload length + payload + checksum (메시지의 전체 길이)

				if len(b) < totalLength {
					mu.Unlock()
					time.Sleep(100 * time.Millisecond)
					continue
				}

				rawMsg := make([]byte, totalLength)
				copy(rawMsg, b[:totalLength])

				data.Next(totalLength) // 버퍼에서 읽어온 데이터를 제거

				mu.Unlock()

				envelope, _, err := ParseNetworkEnvelope(rawMsg) // rawMsg를 NetworkEnvelope로 변환
				if err != nil {
					if sn.Logging {
						log.Printf("Error: %s\n", err)
					}
					continue
				}

				sn.envelopes <- envelope // envelopes 채널에 읽어온 데이터를 전송
			}
		}
	}()

	for {
		buf := make([]byte, 1024) // 1KB 버퍼를 사용

		select {
		case <-sn.serverCloseChan:
			return
		default:
			n, err := sn.conn.Read(buf) // n은 읽어온 데이터의 길이
			if err != nil {
				if err == io.EOF {
					if sn.Logging {
						log.Println("Connection closed by peer")
					}
					return
				}

				if sn.Logging {
					log.Printf("Error: %s\n", err)
				}
				continue
			}

			mu.Lock()

			_, err = data.Write(buf[:n])
			if err != nil {
				if sn.Logging {
					log.Printf("Error: %s\n", err)
				}
				mu.Unlock()
				continue
			}

			mu.Unlock()
		}
	}
}

func (sn *SimpleNode) WaitFor(commands []Command, done <-chan struct{}) (<-chan *NetworkEnvelope, <-chan error) {
	envelopes := make(chan *NetworkEnvelope)
	errors := make(chan error)
	commandsMap := make(map[string]bool)

	for _, command := range commands {
		commandsMap[command.String()] = true
	}

	go func() {
		defer func() {
			close(envelopes)
			close(errors)
		}()
		for {
			select {
			case <-sn.serverCloseChan:
				return
			case <-done:
				return
			case envelope := <-sn.envelopes:
				if envelope == nil {
					continue
				}

				if sn.Logging {
					log.Printf("Recv: %s\n", envelope.Command)
				}

				if envelope.Command.Compare(PingCommand) {
					err := sn.Send(NewPongMessage(envelope.Payload), sn.Network)
					if err != nil {
						errors <- err
						continue
					}

					continue
				}

				if _, ok := commandsMap[envelope.Command.String()]; ok {
					envelopes <- envelope
				}
			}
		}
	}()

	return envelopes, errors
}

func (sn *SimpleNode) HandShake() (<-chan bool, error) {
	msg := DefaultVersionMessage()

	err := sn.Send(msg, sn.Network)
	if err != nil {
		return nil, err
	}

	respChan := make(chan bool)
	done := make(chan struct{})

	envelopes, errors := sn.WaitFor([]Command{VersionCommand, VerAckCommand}, done)

	go func() {
		defer close(respChan)
		defer close(done)

		for {
			select {
			case envelope := <-envelopes:
				if envelope == nil {
					continue
				}

				if envelope.Command.Compare(VerAckCommand) {
					continue
				}

				if envelope.Command.Compare(VersionCommand) {
					ack := NewVerAckMessage()

					err = sn.Send(ack, sn.Network)
					if err != nil {
						if sn.Logging {
							log.Printf("Error: %s\n", err)
						}
						return
					}

					if sn.Logging {
						log.Println("Successfully established connection")
					}

					respChan <- true
					return
				}

			case err := <-errors:
				if err != nil {
					if sn.Logging {
						log.Printf("Error: %s\n", err)
					}
					return
				}
			}
		}
	}()

	return respChan, nil
}
