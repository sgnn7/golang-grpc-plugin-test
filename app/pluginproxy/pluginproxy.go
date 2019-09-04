package pluginproxy

import (
	"io"
	"log"
	"net"
	"sync"
)

func handleConnection(clientConn net.Conn, pluginConn net.Conn, targetAddr string) {
	log.Printf("Got a new connection to passthrough socket from %v",
		pluginConn.RemoteAddr().String())

	log.Printf("Starting backend server connection to %s...", targetAddr)
	backendConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Backend server connect dial Error: %v", err)
		return
	}

	pluginExit := false

	defer func() {
		pluginExit = true
		log.Println("Closing backend connection...")
		backendConn.Close()
	}()

	log.Println("Shuttling bytes...")

	wg := &sync.WaitGroup{}
	wg.Add(3)

	log.Println("Linking connection from backend to plugin socket...")
	go proxyConnection(backendConn, clientConn, &pluginExit, wg)

	log.Println("Linking connection from plugin socket to backend...")
	go proxyConnection(pluginConn, backendConn, &pluginExit, wg)

	log.Println("Linking connection from client to plugin...")
	go proxyConnection(clientConn, pluginConn, &pluginExit, wg)

	log.Println("Shuttling data...")
	wg.Wait()

	log.Println("Closing client connection...")
	clientConn.Close()
	backendConn.Close()

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
		pluginConn, err := listenSocket.Accept()
		if err != nil {
			log.Printf("Listen error: %v", err)
		}

		log.Println("New connection from plugin to broker accepted")
		handleConnection(clientConn, pluginConn, targetAddr)

		log.Printf("Closing passthrough connection of %s...", listenSocket.Addr().String())
		pluginConn.Close()
		log.Printf("Passthrough connection closed")

		log.Printf("Closing passthrough listener of %s...", listenSocket.Addr().String())
		listenSocket.Close()
		log.Printf("Passthrough listener closed")
	}()

	return listenSocket.Addr().String(), nil
}
