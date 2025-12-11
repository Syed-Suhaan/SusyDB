package core

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
)

// ErrMaxKeysExceeded is returned when the store has reached its key limit.
var ErrMaxKeysExceeded = fmt.Errorf("ERR max number of keys exceeded")

// Set stores a key-value pair with an optional Time-To-Live (TTL).
func (s *KVStore) Set(key string, value string, ttlSeconds int64) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	_, exists := shard.data[key]
	if !exists {
		// Strictly enforce atomic MaxKeys check
		if s.MaxKeys > 0 {
			currentKeys := atomic.LoadInt64(&s.keyCount)
			if currentKeys >= int64(s.MaxKeys) {
				return ErrMaxKeysExceeded
			}
		}
		// Increment count for new key
		atomic.AddInt64(&s.keyCount, 1)
		shard.addKey(key)
	}

	var expiresAt int64 = 0
	if ttlSeconds > 0 {
		expiresAt = time.Now().Add(time.Duration(ttlSeconds) * time.Second).UnixNano()
	}

	shard.data[key] = Entry{
		Value:     value,
		ExpiresAt: expiresAt,
	}
	return nil
}

// Get retrieves a value by its key.
func (s *KVStore) Get(key string) (string, bool, error) {
	shard := s.getShard(key)
	shard.mu.RLock()
	entry, ok := shard.data[key]
	shard.mu.RUnlock()

	if !ok {
		return "", false, nil
	}

	// Check for expiration
	if entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		shard.mu.Lock()
		defer shard.mu.Unlock()

		currentEntry, currentOk := shard.data[key]
		if currentOk && currentEntry.ExpiresAt > 0 && time.Now().UnixNano() > currentEntry.ExpiresAt {
			delete(shard.data, key)
			shard.removeKey(key)
			atomic.AddInt64(&s.keyCount, -1) // Decrement count
		}
		return "", false, nil
	}

	strVal, isString := entry.Value.(string)
	if !isString {
		return "", true, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	return strVal, true, nil
}

// IncrBy atomically increments the integer value of a key by delta.
func (s *KVStore) IncrBy(key string, delta int64) (int64, error) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	entry, ok := shard.data[key]

	if ok && entry.ExpiresAt > 0 && time.Now().UnixNano() > entry.ExpiresAt {
		delete(shard.data, key)
		shard.removeKey(key)
		atomic.AddInt64(&s.keyCount, -1)
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
		// New key creation
		if s.MaxKeys > 0 {
			currentKeys := atomic.LoadInt64(&s.keyCount)
			if currentKeys >= int64(s.MaxKeys) {
				return 0, ErrMaxKeysExceeded
			}
		}
		entry = Entry{ExpiresAt: 0}
		shard.addKey(key)
		atomic.AddInt64(&s.keyCount, 1)
	}

	newVal := currentVal + delta
	entry.Value = fmt.Sprintf("%d", newVal)
	shard.data[key] = entry

	return newVal, nil
}

// Delete explicitly removes a key from the store.
func (s *KVStore) Delete(key string) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	if _, exists := shard.data[key]; exists {
		delete(shard.data, key)
		shard.removeKey(key)
		atomic.AddInt64(&s.keyCount, -1)
	}
}
