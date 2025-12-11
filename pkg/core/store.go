package core

import (
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

const ShardCount = 32

// Entry represents a single record in the database.
type Entry struct {
	Value     interface{}
	ExpiresAt int64
}

// Shard reduces lock contention by splitting the DB.
type Shard struct {
	mu       sync.RWMutex
	data     map[string]Entry
	keys     []string
	keyIndex map[string]int
}

// KVStore is the core in-memory database structure.
// Uses Sharding to allow high concurrent throughput.
type KVStore struct {
	shards    []*Shard
	Hub       *Hub
	startTime time.Time
	MaxKeys   int
	keyCount  int64 // Atomic counter for total keys
}

// NewKVStore initializes a new sharded Key-Value Store.
func NewKVStore() *KVStore {
	s := &KVStore{
		shards:    make([]*Shard, ShardCount),
		Hub:       NewHub(),
		startTime: time.Now(),
		MaxKeys:   0,
		keyCount:  0,
	}
	for i := 0; i < ShardCount; i++ {
		s.shards[i] = &Shard{
			data:     make(map[string]Entry),
			keys:     make([]string, 0),
			keyIndex: make(map[string]int),
		}
	}
	return s
}

// NewKVStoreWithLimit initializes a KVStore with a maximum key limit.
func NewKVStoreWithLimit(maxKeys int) *KVStore {
	s := &KVStore{
		shards:    make([]*Shard, ShardCount),
		Hub:       NewHub(),
		startTime: time.Now(),
		MaxKeys:   maxKeys,
		keyCount:  0,
	}
	for i := 0; i < ShardCount; i++ {
		s.shards[i] = &Shard{
			data:     make(map[string]Entry),
			keys:     make([]string, 0),
			keyIndex: make(map[string]int),
		}
	}
	return s
}

// getShard returns the specific shard for a given key.
func (s *KVStore) getShard(key string) *Shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return s.shards[h.Sum32()%ShardCount]
}

// addKey adds a key to the shard's keys slice.
// Must be called while holding the SHARD'S write lock.
func (s *Shard) addKey(key string) {
	s.keyIndex[key] = len(s.keys)
	s.keys = append(s.keys, key)
}

// removeKey removes a key from the shard's keys slice.
// Must be called while holding the SHARD'S write lock.
func (s *Shard) removeKey(key string) {
	idx, exists := s.keyIndex[key]
	if !exists {
		return
	}
	lastIdx := len(s.keys) - 1
	if idx != lastIdx {
		lastKey := s.keys[lastIdx]
		s.keys[idx] = lastKey
		s.keyIndex[lastKey] = idx
	}
	s.keys = s.keys[:lastIdx]
	delete(s.keyIndex, key)
}

// Info aggregates stats.
func (s *KVStore) Info() string {
	uptime := time.Since(s.startTime).Seconds()
	totalKeys := atomic.LoadInt64(&s.keyCount)
	return fmt.Sprintf("# Server\r\nsubydb_version:1.3.0\r\nuptime_in_seconds:%.0f\r\n\r\n# Stats\r\nkeys:%d\r\n", uptime, totalKeys)
}
