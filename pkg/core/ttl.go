package core

import (
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	gcSampleSize       = 20   // Number of keys to sample per GC cycle
	gcExpiredThreshold = 0.25 // Repeat if >25% expired
	gcMaxCycles        = 10   // Max cycles per tick
)

// StartGC starts a background goroutine for "Active Expiry".
// Uses Redis-style probabilistic algorithm with TRUE O(1) complexity.
// Iterates over all shards to ensure global cleanup.
func (s *KVStore) StartGC() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			// Run GC on all shards
			// We can do this sequentially to avoid CPU spikes
			for _, shard := range s.shards {
				s.gcShardCycle(shard)
			}
		}
	}()
}

// gcShardCycle performs probabilistic garbage collection on a SINGLE shard.
func (s *KVStore) gcShardCycle(shard *Shard) {
	for cycle := 0; cycle < gcMaxCycles; cycle++ {
		expired := s.sampleAndCleanShard(shard)
		if expired < int(float64(gcSampleSize)*gcExpiredThreshold) {
			break
		}
	}
}

// sampleAndCleanShard samples random keys from a specific shard.
func (s *KVStore) sampleAndCleanShard(shard *Shard) int {
	shard.mu.Lock()
	defer shard.mu.Unlock()

	keyCount := len(shard.keys)
	if keyCount == 0 {
		return 0
	}

	sampleSize := gcSampleSize
	if keyCount < sampleSize {
		sampleSize = keyCount
	}

	now := time.Now().UnixNano()
	expired := 0

	for i := 0; i < sampleSize; i++ {
		// O(1) random access into keys slice
		idx := rand.Intn(keyCount)
		key := shard.keys[idx]

		entry, exists := shard.data[key]
		if exists && entry.ExpiresAt > 0 && now > entry.ExpiresAt {
			delete(shard.data, key)
			shard.removeKey(key)
			atomic.AddInt64(&s.keyCount, -1) // Decrement count
			expired++
			keyCount = len(shard.keys)
			if keyCount == 0 {
				break
			}
		}
	}

	return expired
}
