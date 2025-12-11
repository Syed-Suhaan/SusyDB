package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/Syed-Suhaan/SusyDB/pkg/core"
)

// Start initializes the TCP server and listens for incoming connections.
func Start(store *core.KVStore, addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to bind to port %s: %v\n", addr, err)
		return
	}
	defer listener.Close()

	// Connection Semaphore to limit concurrency
	// 5000 concurrent clients is a safe upper bound to prevent OOM
	maxClients := 5000
	sem := make(chan struct{}, maxClients)

	fmt.Printf("SusyDB started on %s (Max Clients: %d)\n", addr, maxClients)
	fmt.Println("Ready to accept connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Connection error: %v\n", err)
			continue
		}

		// Try to acquire a token
		select {
		case sem <- struct{}{}:
			// Acquired
			go func() {
				defer func() { <-sem }() // Release token
				handleClient(conn, store)
			}()
		default:
			// Server full - reject immediately
			fmt.Println("Max connections reached, rejecting client.")
			conn.Close()
		}
	}
}

func handleClient(conn net.Conn, store *core.KVStore) {
	defer conn.Close()

	// Panic Recovery Middleware
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[%s] Panic recovered: %v\n", conn.RemoteAddr(), r)
		}
	}()

	reader := bufio.NewReader(conn)

	for {
		// Set deadline to prevent hanging connections (5 minute timeout)
		conn.SetDeadline(time.Now().Add(300 * time.Second))

		// Peek at the first byte to determine protocol
		peek, err := reader.Peek(1)
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}

		var parts []string
		if peek[0] == '*' {
			// Binary RESP protocol
			parts, err = ParseRESP(reader)
			if err != nil {
				// Protocol error
				conn.Write([]byte(fmt.Sprintf("-ERR Protocol error: %v\r\n", err)))
				return
			}
		} else {
			// Inline text protocol (Telnet)
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error:", err)
				}
				return
			}
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			parts = parseCommand(line)
		}

		if len(parts) == 0 {
			continue
		}

		// Dispatch command
		cmdName := strings.ToUpper(parts[0])
		if handler, exists := Handlers[cmdName]; exists {
			response := handler(conn, store, parts)
			if response != nil {
				conn.Write(response)
			}
		} else {
			conn.Write([]byte("-ERR unknown command\r\n"))
		}
	}
}
