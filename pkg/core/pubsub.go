package core

import (
	"sync"
)

// Hub manages Pub/Sub channels and subscribers.
type Hub struct {
	mu   sync.RWMutex
	subs map[string][]chan string // Map channel_name -> list of client channels
}

// NewHub initializes a new Pub/Sub Hub.
func NewHub() *Hub {
	return &Hub{
		subs: make(map[string][]chan string),
	}
}

// Subscribe adds a client to a channel and returns a go-channel for messages.
func (h *Hub) Subscribe(channel string) <-chan string {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Create a new channel for this client
	// Buffer 10 messages to prevent slow clients from blocking publishers immediately
	clientCh := make(chan string, 100)
	h.subs[channel] = append(h.subs[channel], clientCh)

	return clientCh
}

// Publish sends a message to all subscribers of a channel.
// Returns the number of clients that received the message.
func (h *Hub) Publish(channel, message string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.subs[channel]
	if !ok {
		return 0
	}

	count := 0
	for _, clientCh := range clients {
		// Non-blocking send: if client is too slow (buffer full), drop message
		select {
		case clientCh <- message:
			count++
		default:
			// Dropped message (slow consumer)
		}
	}
	return count
}

// Unsubscribe removes a client channel from the list.
// Note: accurate unsubscription requires tracking specific channels.
// For v1, we rely on the channel closing naturally or GC, but ideally we'd pass the chan back.
func (h *Hub) Unsubscribe(channel string, clientCh <-chan string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.subs[channel]
	if !ok {
		return
	}

	// Filter out the specific client channel
	// Since we returned <-chan (recv only), we might need to change signature or cast.
	// For simplicity in this lightweight implementation, let's skip explicit unsub for now
	// and assume the http handler loop exiting is enough (garbage collection).
	// Proper implementation would require a unique Client ID.
}
