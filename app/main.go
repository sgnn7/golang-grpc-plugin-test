package main

import (
	"log"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	go_plugin "github.com/hashicorp/go-plugin"

	"github.com/sgnn7/golang-grpc-plugin-test/app/listener"
	app_plugin "github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	tcp_connector "github.com/sgnn7/golang-grpc-plugin-test/app/plugin/connector/tcp"
	"github.com/sgnn7/golang-grpc-plugin-test/app/pluginproxy"
)

const TargetAddress = "localhost:8080"
const ListenerAddress = ":9090"

type PluginManager struct {
}

func printPluginInfo(infoMap map[string]string) {
	log.Println("---------------------------")
	log.Println("Plugin Info:")
	for key, value := range infoMap {
		log.Println("Key:", key, "Value:", value)
	}
	log.Println("---------------------------")
}

func (manager *PluginManager) StartPlugin(pluginName string) error {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pluginFile := filepath.Join(currDir, pluginName+".plugin")

	log.Printf("Starting plugin: %s", pluginFile)

	client := go_plugin.NewClient(&go_plugin.ClientConfig{
		Cmd:        exec.Command(pluginFile),
		Managed:    true,
		SyncStdout: os.Stdout,
		SyncStderr: os.Stderr,

		HandshakeConfig: app_plugin.HandshakeConfig,
		Plugins:         app_plugin.GetPluginMap(nil),
	})

	rpcclient, err := client.Client()

	if err != nil {
		log.Printf("Failed to get RPC Client: %v", err)
		client.Kill()
		return err
	}

	rawPluginInterface, err := rpcclient.Dispense(tcp_connector.InterfaceName)
	if err != nil {
		log.Printf("Failed to get interface: %s error: %v", tcp_connector.InterfaceName, err)
		return err
	}

	tcpConnectorObj := rawPluginInterface.(tcp_connector.ITCPConnector)

	log.Println("PluginInfo:")
	printPluginInfo(tcpConnectorObj.PluginInfo())

	programExitChan := make(chan bool, 1)
	clientConnChannel, err := listener.StartPluginListeningProcess(ListenerAddress, programExitChan)
	if err != nil {
		log.Printf("Failed to open client-facing socket address/file: %v", err)
		return err
	}

	go func() {
		for clientConn := range clientConnChannel {
			log.Println()
			log.Println("Creating passthough socket...")

			localPassthroughAddress, err := pluginproxy.CreatePassthroughProxy(clientConn, TargetAddress)
			if err != nil {
				log.Printf("Failed to open a shared socket address/file: %v", err)
				continue
			}

			log.Println("Activating plugin...")
			log.Printf("RCP call to plugin response: %v",
				tcpConnectorObj.Connect(localPassthroughAddress))
		}
	}()

	return nil
}

func main() {
	pluginMgr := &PluginManager{}
	pluginMgr.StartPlugin("testplugin")

	for {
		//		log.Printf("App main status: OK")
		time.Sleep(1000 * time.Millisecond)
	}
}
