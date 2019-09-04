package listener

import (
	"log"
	"net"
)

func StartPluginListeningProcess(pluginExit chan bool) (chan net.Conn, error) {
	listenSocket, err := net.Listen("tcp", ":9090")
	if err != nil {
		return nil, err
	}

	isListening := true

	go func() {
		<-pluginExit
		log.Println("Closing client-facing listening socket")
		isListening = false
		listenSocket.Close()
	}()

	clientConnChan := make(chan net.Conn, 100)

	go func() {
		for {
			conn, err := listenSocket.Accept()
			if err != nil {
				if !isListening {
					return
				}

				log.Printf("Listen error on client-facing socket: %v", err)
				continue
			}

			clientConnChan <- conn
		}
	}()

	return clientConnChan, nil
}
