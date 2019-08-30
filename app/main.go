package main

import (
	"log"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-plugin"

	app_plugin "github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	"github.com/sgnn7/golang-grpc-plugin-test/app/plugin/echoer"
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

	rawPlugingInterface, err := rpcclient.Dispense(echoer.InterfaceName)
	if err != nil {
		log.Printf("Failed to get interface: %s error: %s", echoer.InterfaceName, err)
		return err
	}

	echoerObj := rawPlugingInterface.(echoer.IEcho)

	go func() {
		for {
			log.Printf("RCP call to plugin: %s", echoerObj.Reply("Marco!"))
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
