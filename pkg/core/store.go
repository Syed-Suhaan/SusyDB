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
// The keys slice enables O(1) random access for probabilistic GC.
type KVStore struct {
	mu        sync.RWMutex     // Protects all fields
	data      map[string]Entry // Key -> Entry
	keys      []string         // Slice of all keys for O(1) random access
	keyIndex  map[string]int   // Key -> index in keys slice (for O(1) deletion)
	Hub       *Hub             // Pub/Sub Hub
	startTime time.Time        // Used for uptime calculation
	MaxKeys   int              // Maximum number of keys allowed (0 = unlimited)
}

// NewKVStore initializes a new, empty Key-Value Store with no key limit.
func NewKVStore() *KVStore {
	return &KVStore{
		data:      make(map[string]Entry),
		keys:      make([]string, 0),
		keyIndex:  make(map[string]int),
		Hub:       NewHub(),
		startTime: time.Now(),
		MaxKeys:   0, // 0 = unlimited
	}
}

// NewKVStoreWithLimit initializes a KVStore with a maximum key limit.
// When the limit is reached, Set operations will return an error.
func NewKVStoreWithLimit(maxKeys int) *KVStore {
	return &KVStore{
		data:      make(map[string]Entry),
		keys:      make([]string, 0),
		keyIndex:  make(map[string]int),
		Hub:       NewHub(),
		startTime: time.Now(),
		MaxKeys:   maxKeys,
	}
}

// addKey adds a key to the keys slice (call only when key is new).
// Must be called while holding the write lock.
func (s *KVStore) addKey(key string) {
	s.keyIndex[key] = len(s.keys)
	s.keys = append(s.keys, key)
}

// removeKey removes a key from the keys slice using swap-and-pop for O(1).
// Must be called while holding the write lock.
func (s *KVStore) removeKey(key string) {
	idx, exists := s.keyIndex[key]
	if !exists {
		return
	}

	// Swap with last element
	lastIdx := len(s.keys) - 1
	if idx != lastIdx {
		lastKey := s.keys[lastIdx]
		s.keys[idx] = lastKey
		s.keyIndex[lastKey] = idx
	}

	// Pop last element
	s.keys = s.keys[:lastIdx]
	delete(s.keyIndex, key)
}

// Info returns a formatted string containing server statistics.
// This is used by the INFO command for observability.
func (s *KVStore) Info() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	uptime := time.Since(s.startTime).Seconds()
	keyCount := len(s.data)

	// Format mimics the Redis INFO command output style
	return fmt.Sprintf("# Server\r\nsubydb_version:1.2.0\r\nuptime_in_seconds:%.0f\r\n\r\n# Stats\r\nkeys:%d\r\n", uptime, keyCount)
}
