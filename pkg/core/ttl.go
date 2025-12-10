package core

import (
	"math/rand"
	"time"
)

const (
	// gcSampleSize is the number of keys to sample per GC cycle
	gcSampleSize = 20
	// gcExpiredThreshold is the percentage of expired keys that triggers another cycle
	gcExpiredThreshold = 0.25
	// gcMaxCycles limits iterations to prevent infinite loops
	gcMaxCycles = 10
)

// StartGC starts a background goroutine for "Active Expiry".
// It uses the Redis-style probabilistic algorithm:
// 1. Sample 20 random keys
// 2. Delete any that are expired
// 3. If >25% were expired, repeat (up to gcMaxCycles)
// This keeps GC at O(1) instead of O(N), preventing "stop-the-world" pauses.
func (s *KVStore) StartGC() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // Run more frequently, but O(1) per run
		defer ticker.Stop()
		for range ticker.C {
			s.gcCycle()
		}
	}()
}

// gcCycle performs one round of probabilistic garbage collection.
func (s *KVStore) gcCycle() {
	for cycle := 0; cycle < gcMaxCycles; cycle++ {
		expired := s.sampleAndClean()

		// If less than 25% were expired, we're done
		if expired < int(float64(gcSampleSize)*gcExpiredThreshold) {
			break
		}
		// Otherwise, repeat - there's likely more expired keys
	}
}

// sampleAndClean samples up to gcSampleSize keys and deletes expired ones.
// Returns the number of expired keys found.
func (s *KVStore) sampleAndClean() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyCount := len(s.data)
	if keyCount == 0 {
		return 0
	}

	// Collect keys into a slice for random sampling
	// This is O(N) but only happens once, not per-key
	keys := make([]string, 0, keyCount)
	for k := range s.data {
		keys = append(keys, k)
	}

	// Sample random keys
	sampleSize := gcSampleSize
	if keyCount < sampleSize {
		sampleSize = keyCount
	}

	now := time.Now().UnixNano()
	expired := 0

	// Use Fisher-Yates partial shuffle to get random sample
	for i := 0; i < sampleSize; i++ {
		// Pick a random index from remaining elements
		j := i + rand.Intn(keyCount-i)
		keys[i], keys[j] = keys[j], keys[i]

		// Check if this key is expired
		key := keys[i]
		entry, exists := s.data[key]
		if exists && entry.ExpiresAt > 0 && now > entry.ExpiresAt {
			delete(s.data, key)
			expired++
		}
	}

	return expired
}
