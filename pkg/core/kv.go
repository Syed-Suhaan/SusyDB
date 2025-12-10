package core

import (
	"fmt"
	"strconv"
	"time"
)

// Set stores a key-value pair with an optional Time-To-Live (TTL).
// If ttlSeconds > 0, the key will automatically expire after that duration.
func (s *KVStore) Set(key string, value string, ttlSeconds int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt int64 = 0
	if ttlSeconds > 0 {
		expiresAt = time.Now().Add(time.Duration(ttlSeconds) * time.Second).UnixNano()
	}

	s.data[key] = Entry{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

// Get retrieves a value by its key.
// It implements "Lazy Expiry": if a key is accessed after its expiration time,
// it is passively deleted and returned as not found.
func (s *KVStore) Get(key string) (string, bool, error) {
	s.mu.RLock()
	entry, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return "", false, nil
	}

	// Check for expiration
	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		// Key has expired. We need to acquire a Write Lock to delete it.
		// "Double-Checked Locking": verify expiration again after acquiring the lock
		// to prevent race conditions where another goroutine might have already deleted it.
		s.mu.Lock()
		defer s.mu.Unlock()

		currentEntry, currentOk := s.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(s.data, key)
		}
		return "", false, nil
	}

	// Type assertion to ensure we are returning a string value
	strVal, isString := entry.Value.(string)
	if !isString {
		return "", true, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return strVal, true, nil
}

// IncrBy atomically increments the integer value of a key by delta.
// If the key does not exist, it is set to 0 before performing the operation.
// Returns an error if the value cannot be parsed as an integer.
func (s *KVStore) IncrBy(key string, delta int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]

	// Handle expiry inside the critical section
	if ok && entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		delete(s.data, key)
		ok = false
	}

	var currentVal int64 = 0
	if ok {
		strVal, isString := entry.Value.(string)
		if !isString {
			return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		var err error
		currentVal, err = strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("ERR value is not an integer or out of range")
		}
	} else {
		// Initialize new counter
		entry = Entry{ExpiresAt: 0}
	}

	newVal := currentVal + delta
	entry.Value = fmt.Sprintf("%d", newVal) // Store back as string
	s.data[key] = entry

	return newVal, nil
}

// Delete explicitly removes a key from the store.
func (s *KVStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
