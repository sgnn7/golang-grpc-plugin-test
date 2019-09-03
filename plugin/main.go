package main

import (
	"log"
	"time"

	"github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	tcp_connector "github.com/sgnn7/golang-grpc-plugin-test/app/plugin/connector/tcp"
)

type tcpConnectorImpl struct {
	startTime time.Time
}

var tcpConnectorStartTime = time.Now()

func TCPConnector() tcp_connector.ITCPConnector {
	return &tcpConnectorImpl{}
}

func (p *tcpConnectorImpl) Connect(address string) error {
	log.Printf("Plugin Connect: %s", address)
	return nil
}

func main() {
	pluginOpts := &plugin.PluginOpts{
		TCPConnector: TCPConnector,
		RunAsPlugin:  true,
	}

	plugin.StartPlugin(pluginOpts, make(chan bool))

	//	tcpConnectorPlugin := TCPConnector()
	for {
		//	log.Printf("Plugin self-test")
		//		log.Printf("Plugin self-test: %s\n", tcpConnectorPlugin.Connect("tcp://localhost:8080"))
		time.Sleep(1000 * time.Millisecond)
	}
}
