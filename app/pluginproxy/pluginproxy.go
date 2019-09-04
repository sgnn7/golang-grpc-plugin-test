package pluginproxy

import (
	"io"
	"log"
	"net"
	"sync"
)

func handleConnection(clientConn net.Conn, fromPluginConn net.Conn, targetAddr string) {
	log.Printf("Got a new connection to passthrough socket from %v",
		fromPluginConn.RemoteAddr().String())

	log.Printf("Starting backend server connection to %s...", targetAddr)
	toBackendConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Backend server connect dial Error: %v", err)
		return
	}

	pluginExit := false

	defer func() {
		pluginExit = true
		log.Println("Closing backend connection...")
		toBackendConn.Close()
	}()

	log.Println("Shuttling bytes...")

	wg := &sync.WaitGroup{}
	wg.Add(2)

	log.Println("Linking connection from backend to plugin socket...")
	go proxyConnection(toBackendConn, clientConn, &pluginExit, wg)

	log.Println("Linking connection from plugin socket to backend...")
	go proxyConnection(fromPluginConn, toBackendConn, &pluginExit, wg)

	log.Println("Shuttling data...")
	wg.Wait()

	log.Println("Closing client connection...")
	clientConn.Close()
	toBackendConn.Close()

	log.Printf("Connection closed")
	log.Printf("Passthrough subroutine done")
}

func proxyConnection(fromConn net.Conn, toConn net.Conn, pluginExit *bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if !*pluginExit {
		if _, err := io.Copy(toConn, fromConn); err != nil {
			if !*pluginExit {
				log.Printf("Error copying: %v", err)
			} else {
				log.Printf("Proxy io.Copy done from %s to %s",
					fromConn.LocalAddr().String(), toConn.LocalAddr().String())
			}
			return
		} else {
			log.Printf("Copied some data from %s to %s...",
				fromConn.LocalAddr().String(), toConn.LocalAddr().String())
		}
	}

	log.Printf("Proxy io.Copy done from %s to %s",
		fromConn.LocalAddr().String(), toConn.LocalAddr().String())

	*pluginExit = true
}

func CreatePassthroughProxy(clientConn net.Conn, targetAddr string) (string, error) {
	listenSocket, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}

	log.Printf("Created passthrough local ephemeral connection for a plugin at %s",
		listenSocket.Addr().String())

	go func() {
		// One time use - no need to loop
		fromPluginConn, err := listenSocket.Accept()
		if err != nil {
			log.Printf("Listen error: %v", err)
		}

		log.Println("New connection from plugin to broker accepted")
		handleConnection(clientConn, fromPluginConn, targetAddr)

		log.Printf("Closing passthrough connection of %s...", listenSocket.Addr().String())
		fromPluginConn.Close()
		log.Printf("Passthrough connection closed")

		log.Printf("Closing passthrough listener of %s...", listenSocket.Addr().String())
		listenSocket.Close()
		log.Printf("Passthrough listener closed")
	}()

	return listenSocket.Addr().String(), nil
}
