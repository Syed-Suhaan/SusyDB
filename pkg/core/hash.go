package core

import (
	"fmt"
	"time"
)

// HSet sets a specific field in a Hash map stored at key.
// Creates the Hash if it doesn't verify existence.
func (s *KVStore) HSet(key, field, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[key]
	if exists {
		// Clean up if expired
		if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
			exists = false
			delete(s.data, key)
			s.removeKey(key)
		} else {
			// Verify type is Map
			_, isMap := entry.Value.(map[string]string)
			if !isMap {
				return fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
			}
		}
	}

	if !exists {
		entry = Entry{
			Value:     make(map[string]string),
			ExpiresAt: 0,
		}
	}

	hash := entry.Value.(map[string]string)
	hash[field] = value
	s.data[key] = entry
	return nil
}

// HGet retrieves a specific field from a Hash.
func (s *KVStore) HGet(key, field string) (string, bool, error) {
	s.mu.RLock()
	entry, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false, nil
	}

	// Double-checked locking for expiry
	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		s.mu.Lock()
		defer s.mu.Unlock()
		currentEntry, currentOk := s.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(s.data, key)
			s.removeKey(key)
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
// Returns a copy of the map to ensure thread safety outside the lock.
func (s *KVStore) HGetAll(key string) (map[string]string, bool, error) {
	s.mu.RLock()
	entry, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return nil, false, nil
	}

	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		s.mu.Lock()
		defer s.mu.Unlock()
		currentEntry, currentOk := s.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(s.data, key)
			s.removeKey(key)
		}
		return nil, false, nil
	}

	hash, isMap := entry.Value.(map[string]string)
	if !isMap {
		return nil, true, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// Return a copy to prevent race conditions on the map reference
	copyMap := make(map[string]string)
	for k, v := range hash {
		copyMap[k] = v
	}
	return copyMap, true, nil
}

// HDel deletes a specific field from a Hash.
// Returns true if the field was present and deleted.
func (s *KVStore) HDel(key, field string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[key]
	if !exists {
		return false, nil
	}

	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		delete(s.data, key)
		s.removeKey(key)
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
