package core

import (
	"fmt"
	"sync/atomic"
	"time"
)

// HSet sets a specific field in a Hash map stored at key.
func (s *KVStore) HSet(key, field, value string) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	entry, exists := shard.data[key]
	if exists {
		// Clean up if expired
		if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
			exists = false
			delete(shard.data, key)
			shard.removeKey(key)
			atomic.AddInt64(&s.keyCount, -1)
		} else {
			// Verify type is Map
			_, isMap := entry.Value.(map[string]string)
			if !isMap {
				return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
		}
	}

	if !exists {
		// New key creation
		if s.MaxKeys > 0 {
			currentKeys := atomic.LoadInt64(&s.keyCount)
			if currentKeys >= int64(s.MaxKeys) {
				return ErrMaxKeysExceeded
			}
		}

		entry = Entry{
			Value:     make(map[string]string),
			ExpiresAt: 0,
		}
		shard.addKey(key)
		atomic.AddInt64(&s.keyCount, 1)
	}

	hash := entry.Value.(map[string]string)
	hash[field] = value
	shard.data[key] = entry
	return nil
}

// HGet retrieves a specific field from a Hash.
func (s *KVStore) HGet(key, field string) (string, bool, error) {
	shard := s.getShard(key)
	shard.mu.RLock()
	entry, exists := shard.data[key]
	shard.mu.RUnlock()

	if !exists {
		return "", false, nil
	}

	// Double-checked locking for expiry
	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		shard.mu.Lock()
		defer shard.mu.Unlock()
		currentEntry, currentOk := shard.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(shard.data, key)
			shard.removeKey(key)
			atomic.AddInt64(&s.keyCount, -1)
		}
		return "", false, nil
	}

	hash, isMap := entry.Value.(map[string]string)
	if !isMap {
		return "", true, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	val, ok := hash[field]
	return val, ok, nil
}

// HGetAll returns all fields and values in a Hash.
func (s *KVStore) HGetAll(key string) (map[string]string, bool, error) {
	shard := s.getShard(key)
	shard.mu.RLock()
	entry, exists := shard.data[key]
	shard.mu.RUnlock()

	if !exists {
		return nil, false, nil
	}

	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		shard.mu.Lock()
		defer shard.mu.Unlock()
		currentEntry, currentOk := shard.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(shard.data, key)
			shard.removeKey(key)
			atomic.AddInt64(&s.keyCount, -1)
		}
		return nil, false, nil
	}

	hash, isMap := entry.Value.(map[string]string)
	if !isMap {
		return nil, true, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Return a copy to prevent race conditions
	copyMap := make(map[string]string)
	for k, v := range hash {
		copyMap[k] = v
	}
	return copyMap, true, nil
}

// HDel deletes a specific field from a Hash.
// Returns true if the field was present and deleted.
func (s *KVStore) HDel(key, field string) (bool, error) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	entry, exists := shard.data[key]
	if !exists {
		return false, nil
	}

	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		delete(shard.data, key)
		shard.removeKey(key)
		atomic.AddInt64(&s.keyCount, -1)
		return false, nil
	}

	hash, isMap := entry.Value.(map[string]string)
	if !isMap {
		return false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	_, fieldExists := hash[field]
	if fieldExists {
		delete(hash, field)
		// Optimization: If hash is empty, we could delete the key here,
		// but standard Redis behavior keeps the key until explicitly deleted.
		return true, nil
	}
	return false, nil
}
