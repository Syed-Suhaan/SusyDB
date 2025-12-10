package core

import (
	"fmt"
	"sync"
	"time"
)

// Entry represents a single record in the database.
// Value is an interface{} to support multiple data types:
// - string: for standard Set/Get
// - map[string]string: for Hash operations
// ExpiresAt is the Unix nanosecond timestamp when the key should expire.
// If ExpiresAt is 0, the key persists indefinitely.
type Entry struct {
	Value     interface{}
	ExpiresAt int64
}

// KVStore is the core in-memory database structure.
// It uses a sync.RWMutex to ensure thread safety for concurrent access.
type KVStore struct {
	mu        sync.RWMutex // Protects the data map
	data      map[string]Entry
	Hub       *Hub      // Pub/Sub Hub
	startTime time.Time // Used for uptime calculation
}

// NewKVStore initializes a new, empty Key-Value Store.
func NewKVStore() *KVStore {
	return &KVStore{
		data:      make(map[string]Entry),
		Hub:       NewHub(),
		startTime: time.Now(),
	}
}

// Info returns a formatted string containing server statistics.
// This is used by the INFO command for observability.
func (s *KVStore) Info() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uptime := time.Since(s.startTime).Seconds()
	keyCount := len(s.data)

	// Format mimics the Redis INFO command output style
	return fmt.Sprintf("# Server\r\nsubydb_version:1.0.0\r\nuptime_in_seconds:%.0f\r\n\r\n# Stats\r\nkeys:%d\r\n", uptime, keyCount)
}
