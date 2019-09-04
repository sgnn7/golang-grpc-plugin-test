package tcp

import (
	"github.com/hashicorp/go-plugin"
)

const InterfaceName = "TCPConnector"

type ITCPConnector interface {
	Connect(address string) plugin.BasicError
	PluginInfo() map[string]string
}

type TCPConnectorFunc func() ITCPConnector

type TCPConnector struct {
	F TCPConnectorFunc
}
