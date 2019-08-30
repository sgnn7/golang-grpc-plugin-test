package echoer

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type echoerClient struct {
	Broker *plugin.MuxBroker
	Client *rpc.Client
}

func (echoPlugin *Echoer) Client(pluginBroker *plugin.MuxBroker,
	rpcClient *rpc.Client) (interface{}, error) {

	return &echoerClient{
		Broker: pluginBroker,
		Client: rpcClient,
	}, nil
}

func (echoPlugin *echoerClient) Reply(sentString string) string {
	log.Printf("In: Client Reply: %s", sentString)

	var resp string

	err := echoPlugin.Client.Call("Plugin.Reply", &sentString, &resp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("In: Client Reply Response: %s", resp)
	return resp
}
