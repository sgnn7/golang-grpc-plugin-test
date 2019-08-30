package echoer

import (
	"errors"
	// "log"

	"github.com/hashicorp/go-plugin"
)

type echoerServer struct {
	Broker *plugin.MuxBroker
	IEcho  IEcho
}

func (echoPlugin *Echoer) Server(pluginBroker *plugin.MuxBroker) (interface{}, error) {
	if echoPlugin.F == nil {
		return nil, errors.New("Echoer interface not implemented")
	}

	return &echoerServer{
		Broker: pluginBroker,
		IEcho:  echoPlugin.F(),
	}, nil
}

func (echoPlugin *echoerServer) Reply(sentString string, result *string) error {
	// log.Printf("In: Server Reply: %s", sentString)
	*result = echoPlugin.IEcho.Reply(sentString)
	// log.Printf("In: Server Reply Response: %s", *result)
	return nil
}
