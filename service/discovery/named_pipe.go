package discovery

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

func getPipePath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\uni-token`
	}
	return "/tmp/uni-token.sock"
}

func SetupNamedPipe() error {
	pipePath := getPipePath()

	// Remove existing socket file on Unix systems
	if runtime.GOOS != "windows" {
		os.Remove(pipePath)
	}

	var listener net.Listener
	var err error

	if runtime.GOOS == "windows" {
		// On Windows, use named pipes
		listener, err = net.Listen("unix", pipePath)
	} else {
		// On Unix systems, use Unix domain sockets
		listener, err = net.Listen("unix", pipePath)
	}

	if err != nil {
		return err
	}

	fmt.Printf("Named pipe server listening at %s\n", pipePath)

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}

			go handleConnection(conn)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected")

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		fmt.Printf("Received data: %s\n", string(buffer[:n]))
	}
}
