package main

import (
	"log"
	"time"

	"github.com/sgnn7/golang-grpc-plugin-test/app/plugin"
	"github.com/sgnn7/golang-grpc-plugin-test/app/plugin/echoer"
)

type echoerImpl struct {
	startTime time.Time
}

var echoerStartTime = time.Now()

func Echoer() echoer.IEcho {
	return &echoerImpl{}
}

func (p *echoerImpl) Reply(sentString string) string {
	return sentString + " plus plugin-added string"
}

func main() {
	pluginOpts := &plugin.PluginOpts{
		Echoer:      Echoer,
		RunAsPlugin: true,
	}

	plugin.StartPlugin(pluginOpts, make(chan bool))

	echoPlugin := Echoer()
	for {
		log.Printf("Plugin self-test: %s\n", echoPlugin.Reply("self-test"))
		time.Sleep(1000 * time.Millisecond)
	}
}
