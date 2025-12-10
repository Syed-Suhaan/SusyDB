package core

import (
	"time"
)

// StartGC starts a background goroutine for "Active Expiry".
// It periodically sweeps the data store to remove expired keys, ensuring memory is freed
// even if keys are never accessed (solving the "Lazy Expiry" limitation).
func (s *KVStore) StartGC() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			// In a Production DB, we would use a specialized probabilistic algorithm
			// to avoid locking the entire DB. For simplicity, we lock fully here.
			s.mu.Lock()
			now := time.Now().UnixNano()
			for key, entry := range s.data {
				if entry.ExpiresAt > 0 && now > entry.ExpiresAt {
					delete(s.data, key)
				}
			}
			s.mu.Unlock()
		}
	}()
}
