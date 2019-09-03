package tcp

import (
	"errors"
	"log"

	"github.com/hashicorp/go-plugin"
)

type tcpConnectorServer struct {
	Broker        *plugin.MuxBroker
	ITCPConnector ITCPConnector
}

func (tcpConnectorPlugin *TCPConnector) Server(pluginBroker *plugin.MuxBroker) (interface{}, error) {
	if tcpConnectorPlugin.F == nil {
		return nil, errors.New("TCPConnector interface not implemented")
	}

	return &tcpConnectorServer{
		Broker:        pluginBroker,
		ITCPConnector: tcpConnectorPlugin.F(),
	}, nil
}

func (tcpConnectorPlugin *tcpConnectorServer) Connect(address string, result *error) error {
	log.Printf("In: Server Connect: %s", address)
	*result = tcpConnectorPlugin.ITCPConnector.Connect(address)
	log.Printf("In: Server Reply Response: %v", *result)
	return nil
}
