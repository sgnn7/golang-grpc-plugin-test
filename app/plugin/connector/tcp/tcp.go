package tcp

const InterfaceName = "TCPConnector"

type ITCPConnector interface {
	Connect(string) error
}

type TCPConnector struct {
	// TODO: Use TCPConnectorFunc
	F func() ITCPConnector
}

type TCPConnectorFunc func() ITCPConnector
