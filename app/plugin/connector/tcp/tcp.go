package tcp

const InterfaceName = "TCPConnector"

type ITCPConnector interface {
	Connect(address string) error
}

type TCPConnectorFunc func() ITCPConnector

type TCPConnector struct {
	F TCPConnectorFunc
}
