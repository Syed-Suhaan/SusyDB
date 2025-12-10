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

// sampleAndClean samples up to gcSampleSize random keys and deletes expired ones.
// Uses the pre-maintained keys slice for TRUE O(1) sampling.
// Returns the number of expired keys found.
func (s *KVStore) sampleAndClean() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyCount := len(s.keys)
	if keyCount == 0 {
		return 0
	}

	// Sample random keys from the pre-maintained slice - TRUE O(1)
	sampleSize := gcSampleSize
	if keyCount < sampleSize {
		sampleSize = keyCount
	}

	now := time.Now().UnixNano()
	expired := 0

	// Sample random indices without building a new slice
	for i := 0; i < sampleSize; i++ {
		// Pick a random key from the slice
		idx := rand.Intn(keyCount)
		key := s.keys[idx]

		// Check if this key is expired
		entry, exists := s.data[key]
		if exists && entry.ExpiresAt > 0 && now > entry.ExpiresAt {
			delete(s.data, key)
			s.removeKey(key)
			expired++
			// Adjust keyCount since we removed a key
			keyCount = len(s.keys)
			if keyCount == 0 {
				break
			}
		}
	}

	return expired
}
