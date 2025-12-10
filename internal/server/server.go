package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/Syed-Suhaan/susydb/pkg/core"
)

// Start initializes the TCP server and listens for incoming connections.
func Start(store *core.KVStore, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %s: %v\n", addr, err)
		return
	}
	defer listener.Close()

	fmt.Printf("SusyDB started on %s\n", addr)
	fmt.Println("Ready to accept connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			continue
		}
		go handleClient(conn, store)
	}
}

func handleClient(conn net.Conn, store *core.KVStore) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read until newline (simple text protocol)
		message, err := reader.ReadString('\n')
		if err != nil {
			// Connection closed or error
			return
		}

		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		parts := strings.Split(message, " ")
		if len(parts) == 0 {
			continue
		}

		// Normalize command to uppercase
		command := strings.ToUpper(parts[0])

		// Command Dispatcher
		if handler, ok := Handlers[command]; ok {
			response := handler(conn, store, parts)
			if response != nil {
				conn.Write(response)
			}
		} else {
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
