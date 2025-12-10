package core

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	store := NewKVStore()

	tests := []struct {
		name    string
		key     string
		value   string
		wantOk  bool
		wantVal string
	}{
		{"simple string", "name", "Suhaan", true, "Suhaan"},
		{"with spaces", "greeting", "Hello World", true, "Hello World"},
		{"json value", "data", `{"user":"suhaan"}`, true, `{"user":"suhaan"}`},
		{"empty value", "empty", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(tt.key, tt.value, 0)
			if err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}

			got, ok, err := store.Get(tt.key)
			if err != nil {
				t.Errorf("Get() error = %v", err)
				return
			}
			if ok != tt.wantOk {
				t.Errorf("Get() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.wantVal {
				t.Errorf("Get() = %v, want %v", got, tt.wantVal)
			}
		})
	}
}

func TestGetNonExistent(t *testing.T) {
	store := NewKVStore()

	val, ok, err := store.Get("nonexistent")
	if err != nil {
		t.Errorf("Get() error = %v", err)
	}
	if ok {
		t.Errorf("Get() ok = true, want false for nonexistent key")
	}
	if val != "" {
		t.Errorf("Get() = %v, want empty string", val)
	}
}

func TestTTLExpiry(t *testing.T) {
	store := NewKVStore()

	// Set with 1 second TTL
	err := store.Set("volatile", "will_expire", 1)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Should exist immediately
	_, ok, _ := store.Get("volatile")
	if !ok {
		t.Error("Key should exist immediately after set")
	}

	// Wait for expiry
	time.Sleep(1100 * time.Millisecond)

	// Should be gone
	_, ok, _ = store.Get("volatile")
	if ok {
		t.Error("Key should have expired after 1 second")
	}
}

func TestIncrBy(t *testing.T) {
	store := NewKVStore()

	// Increment non-existent key (should start at 0)
	val, err := store.IncrBy("counter", 1)
	if err != nil {
		t.Fatalf("IncrBy() error = %v", err)
	}
	if val != 1 {
		t.Errorf("IncrBy() = %v, want 1", val)
	}

	// Increment again
	val, err = store.IncrBy("counter", 5)
	if err != nil {
		t.Fatalf("IncrBy() error = %v", err)
	}
	if val != 6 {
		t.Errorf("IncrBy() = %v, want 6", val)
	}

	// Negative increment
	val, err = store.IncrBy("counter", -2)
	if err != nil {
		t.Fatalf("IncrBy() error = %v", err)
	}
	if val != 4 {
		t.Errorf("IncrBy() = %v, want 4", val)
	}
}

func TestIncrByNonNumeric(t *testing.T) {
	store := NewKVStore()

	// Set a non-numeric value
	store.Set("text", "hello", 0)

	// Try to increment - should error
	_, err := store.IncrBy("text", 1)
	if err == nil {
		t.Error("IncrBy() should error on non-numeric value")
	}
}

func TestDelete(t *testing.T) {
	store := NewKVStore()

	store.Set("toDelete", "value", 0)
	store.Delete("toDelete")

	_, ok, _ := store.Get("toDelete")
	if ok {
		t.Error("Key should not exist after delete")
	}
}

func TestMaxKeysLimit(t *testing.T) {
	store := NewKVStoreWithLimit(3)

	// Should succeed
	if err := store.Set("k1", "v1", 0); err != nil {
		t.Errorf("Set() error = %v", err)
	}
	if err := store.Set("k2", "v2", 0); err != nil {
		t.Errorf("Set() error = %v", err)
	}
	if err := store.Set("k3", "v3", 0); err != nil {
		t.Errorf("Set() error = %v", err)
	}

	// Should fail - at capacity
	if err := store.Set("k4", "v4", 0); err == nil {
		t.Error("Set() should error when MaxKeys exceeded")
	}

	// Updating existing key should still work
	if err := store.Set("k1", "updated", 0); err != nil {
		t.Errorf("Set() error updating existing key = %v", err)
	}
}
