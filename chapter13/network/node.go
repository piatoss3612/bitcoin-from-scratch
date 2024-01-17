package network

import (
	"bytes"
	"chapter13/utils"
	"fmt"
	"log"
	"net"
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
	for {
		select {
		case <-sn.serverCloseChan:
			return
		default:
			data, err := sn.ReadAll()
			if err != nil {
				if sn.Logging {
					log.Printf("Error: %s\n", err)
				}
				continue
			}

			buf := bytes.NewBuffer(data)

			for buf.Len() > 0 {
				envelope, read, err := ParseNetworkEnvelope(buf.Bytes()) // 버퍼에서 네트워크 메시지를 읽어옴
				if err != nil {
					if sn.Logging {
						log.Printf("Error: %s\n", err)
					}
					break
				}

				buf.Next(read)           // 버퍼에서 읽어온 데이터를 제거
				sn.envelopes <- envelope // envelopes 채널에 읽어온 데이터를 전송
			}
		}
	}
}

func (sn *SimpleNode) ReadAll() ([]byte, error) {
	buf := make([]byte, 1024)

	n, err := sn.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	if !IsNetworkEnvelope(buf[:n]) { // 네트워크 메시지인지 확인
		return nil, ErrInvalidNetworkMessage
	}

	payloadLength := utils.LittleEndianToInt(buf[16:20]) // 페이로드 길이를 읽어옴

	fmt.Println("Payload Length:", payloadLength)

	totalLength := 4 + 12 + 4 + payloadLength + 4 // magic + command + payload length + payload + checksum (메시지의 전체 길이)
	readCnt := n

	data := make([]byte, 0, totalLength) // data를 메시지의 전체 길이로 초기화
	data = append(data, buf[:n]...)

	fmt.Println("Total Length:", totalLength, "Read Count:", readCnt)

	for readCnt < totalLength {
		select {
		case <-sn.serverCloseChan:
			return nil, net.ErrClosed
		default:
			n, err := sn.conn.Read(buf) // n은 읽어온 데이터의 길이
			if err != nil {
				return nil, err
			}

			readCnt += n
			data = append(data, buf[:n]...)

			fmt.Println("Total Length:", totalLength, "Read Count:", readCnt)
		}
	}

	fmt.Printf("Remaining Data: %x\n", data[totalLength:])

	return data, nil
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
					if sn.Logging {
						log.Printf("Recv: %s\n", envelope.Command)
					}

					continue
				}

				if envelope.Command.Compare(VersionCommand) {
					if sn.Logging {
						log.Printf("Recv: %s\n", envelope.Command)
					}

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
