package plugin

// Configuration values for Hashicorp go-plugin framework

import (
	"log"

	"github.com/hashicorp/go-plugin"

	"github.com/sgnn7/golang-grpc-plugin-test/app/plugin/echoer"
)

type PluginOpts struct {
	Echoer      echoer.EchoFunc
	RunAsPlugin bool
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MAGIC_COOKIE",
	MagicCookieValue: "FIXME",
}

// StartPlugin starts the Hashicorp go-plugin system over RPC
// between the agent and the plugin, asynchronously.
func StartPlugin(options *PluginOpts, quit chan bool) {
	if !options.RunAsPlugin {
		log.Println("Starting...")
	}

	go func() {
		log.Println("Starting plugin...")

		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins:         GetPluginMap(options),
		})

		log.Println("Plugin ready...")

		quit <- true
		log.Println("Stopping plugin...")
	}()
}

// GetPluginMap returns the plugin map defined Hashicorp go-plugin.
// The reserved parameter should only be used by the RPC receiver (the plugin).
// Otherwise, reserved should be nil for the RPC sender (the mainapp).
func GetPluginMap(options *PluginOpts) map[string]plugin.Plugin {
	var echoerObj echoer.Echoer

	if options != nil {
		echoerObj.F = options.Echoer
	}

	return map[string]plugin.Plugin{
		echoer.InterfaceName: &echoerObj,
	}
}
