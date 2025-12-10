package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/Syed-Suhaan/susydb/pkg/core"
)

func handlePublish(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 3 {
		return []byte("-ERR wrong number of arguments for 'publish' command\r\n")
	}
	channel := parts[1]
	message := strings.Join(parts[2:], " ")

	count := store.Hub.Publish(channel, message)

	// Redis returns integer number of clients that received the message
	return []byte(fmt.Sprintf(":%d\r\n", count))
}

func handleSubscribe(conn net.Conn, store *core.KVStore, parts []string) []byte {
	if len(parts) < 2 {
		return []byte("-ERR wrong number of arguments for 'subscribe' command\r\n")
	}
	channel := parts[1]

	// Register subscription
	subCh := store.Hub.Subscribe(channel)
	// Optionally, we could clean up on exit using defer/Unsubscribe.
	// For simplicity, we just drop the channel on the implementation side (Hub weak ref or GC).
	// In a robust implementation, we would register a "close callback".

	// Acknowledge subscription
	// Redis format: *3\r\n$9\r\nsubscribe\r\n$channel_len\r\nchannel\r\n:1\r\n
	// Simplified: just say subscribed
	conn.Write([]byte(fmt.Sprintf("*3\r\n$9\r\nsubscribe\r\n$%d\r\n%s\r\n:1\r\n", len(channel), channel)))

	// Blocking loop
	for msg := range subCh {
		// Redis Pub/Sub message format:
		// *3\r\n$7\r\nmessage\r\n$channel_len\r\nchannel\r\n$msg_len\r\nmessage\r\n
		response := fmt.Sprintf("*3\r\n$7\r\nmessage\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
			len(channel), channel, len(msg), msg)

		_, err := conn.Write([]byte(response))
		if err != nil {
			// Connection likely closed
			return nil
		}
	}

	return nil
}
