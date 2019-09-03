package tcp

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type tcpConnectorClient struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

func (tcpConnectorPlugin *TCPConnector) Client(pluginBroker *plugin.MuxBroker,
	rpcClient *rpc.Client) (interface{}, error) {
	log.Println("TCP plugin client instantiation")

	return &tcpConnectorClient{
		Broker: pluginBroker,
		Client: rpcClient,
	}, nil
}

func (tcpConnectorClient *tcpConnectorClient) Connect(address string) plugin.BasicError {
	log.Printf("In: Client Connect: %s", address)

	var resp_err plugin.BasicError

	err := tcpConnectorClient.Client.Call("Plugin.Connect", &address, &resp_err)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("In: Client Connect Response: %v", resp_err)
	return resp_err
}
