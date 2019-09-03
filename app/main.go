package main

import (
	"log"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"

	app_plugin "github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	tcp_connector "github.com/sgnn7/golang-grpc-plugin-test/app/plugin/connector/tcp"
)

type PluginManager struct {
}

func (manager *PluginManager) StartPlugin(pluginName string) error {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pluginFile := filepath.Join(currDir, pluginName+".plugin")

	log.Printf("Starting plugin: %s", pluginFile)

	client := plugin.NewClient(&plugin.ClientConfig{
		Cmd:        exec.Command(pluginFile),
		Managed:    true,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,

		HandshakeConfig: app_plugin.HandshakeConfig,
		Plugins:         app_plugin.GetPluginMap(nil),
	})

	rpcclient, err := client.Client()

	if err != nil {
		log.Printf("Failed to get RPC Client: %s", err)
		client.Kill()
		return err
	}

	rawPluginInterface, err := rpcclient.Dispense(tcp_connector.InterfaceName)
	if err != nil {
		log.Printf("Failed to get interface: %s error: %s", tcp_connector.InterfaceName, err)
		return err
	}

	tcpConnectorObj := rawPluginInterface.(tcp_connector.ITCPConnector)

	go func() {
		for {
			log.Println()
			log.Printf("RCP call to plugin response: %v",
				tcpConnectorObj.Connect("tcp://localhost:8080"))
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return nil
}

func main() {
	pluginMgr := &PluginManager{}
	pluginMgr.StartPlugin("testplugin")

	for {
		log.Printf("App main status: OK")
		time.Sleep(1000 * time.Millisecond)
	}
}
