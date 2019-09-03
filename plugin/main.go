package main

import (
	"log"
	"net"
	"time"

	go_plugin "github.com/hashicorp/go-plugin"
	app_plugin "github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	tcp_connector "github.com/sgnn7/golang-grpc-plugin-test/app/plugin/connector/tcp"
)

type tcpConnectorImpl struct {
	startTime time.Time
}

var tcpConnectorStartTime = time.Now()

func TCPConnector() tcp_connector.ITCPConnector {
	return &tcpConnectorImpl{}
}

func (p *tcpConnectorImpl) Connect(address string) go_plugin.BasicError {
	log.Printf("Plugin Connect: %s", address)

	log.Println("Starting plugin client...")
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Plugin Connect Dial Error: %v", err)
		return *go_plugin.NewBasicError(err)
	}

	if err = connection.Close(); err != nil {
		log.Printf("Plugin Connect Close Error: %v", err)
		return *go_plugin.NewBasicError(err)
	}

	if _, err = connection.Write([]byte("foobar")); err != nil {
		log.Printf("Plugin Connect Write Error: %v", err)
		return *go_plugin.NewBasicError(err)
	}

	log.Printf("Plugin Connect OK!")

	return *go_plugin.NewBasicError(nil)
}

func main() {
	pluginOpts := &app_plugin.PluginOpts{
		TCPConnector: TCPConnector,
		RunAsPlugin:  true,
	}

	app_plugin.StartPlugin(pluginOpts, make(chan bool))

	//	tcpConnectorPlugin := TCPConnector()
	for {
		//	log.Printf("Plugin self-test")
		//		log.Printf("Plugin self-test: %s\n", tcpConnectorPlugin.Connect("tcp://localhost:8080"))
		time.Sleep(1000 * time.Millisecond)
	}
}
