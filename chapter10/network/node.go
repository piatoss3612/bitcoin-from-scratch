package network

import (
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
		serverCloseChan: make(chan struct{}),
	}

	return node, nil
}

func (sn *SimpleNode) Close() error {
	close(sn.serverCloseChan)

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

	// s := hex.EncodeToString(buf[:n])
	fmt.Println("Read bytes:", n)

	envelope, err := ParseNetworkEnvelope(buf[:n])
	if err != nil {
		return nil, err
	}

	return envelope, nil
}

func (sn *SimpleNode) WaitFor(commands []Command) (<-chan *NetworkEnvelope, <-chan error) {
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
			default:
				fmt.Println("Waiting for messages...")
				envelope, err := sn.Read()
				if err != nil {
					// fmt.Println("Error reading message:", err)
					// continue
					errors <- err
					return
				}

				if envelope == nil {
					continue
				}

				if envelope.Command.Compare(PingCommand) {
					err = sn.Send(NewPongMessage(envelope.Payload), sn.Network)
					if err != nil {
						errors <- err
						continue
					}

					if sn.Logging {
						log.Printf("Send: %s\n", PongCommand)
					}
				} else if _, ok := commandsMap[envelope.Command.String()]; ok {
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

	envelopes, errors := sn.WaitFor([]Command{VersionCommand, VerAckCommand})

	go func() {
		defer close(respChan)
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
