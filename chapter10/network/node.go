package network

type SimpleNode struct {
	Host    string
	Port    int
	Testnet bool
	Logging bool
}

func NewSimpleNode(host string, port int, testnet, logging bool) (*SimpleNode, error) {
	if port == 0 {
		if testnet {
			port = 18333
		} else {
			port = 8333
		}
	}

	return &SimpleNode{
		Host:    host,
		Port:    port,
		Testnet: testnet,
		Logging: logging,
	}, nil
}
