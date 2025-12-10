package core

import (
	"testing"
)

func TestHSetAndHGet(t *testing.T) {
	store := NewKVStore()

	// Set a field
	err := store.HSet("user:1", "name", "Suhaan")
	if err != nil {
		t.Fatalf("HSet() error = %v", err)
	}

	// Get the field
	val, ok, err := store.HGet("user:1", "name")
	if err != nil {
		t.Fatalf("HGet() error = %v", err)
	}
	if !ok {
		t.Error("HGet() ok = false, want true")
	}
	if val != "Suhaan" {
		t.Errorf("HGet() = %v, want Suhaan", val)
	}
}

func TestHGetNonExistent(t *testing.T) {
	store := NewKVStore()

	// Non-existent hash
	_, ok, _ := store.HGet("noexist", "field")
	if ok {
		t.Error("HGet() should return false for nonexistent hash")
	}

	// Existent hash, non-existent field
	store.HSet("exists", "f1", "v1")
	_, ok, _ = store.HGet("exists", "f2")
	if ok {
		t.Error("HGet() should return false for nonexistent field")
	}
}

func TestHGetAll(t *testing.T) {
	store := NewKVStore()

	store.HSet("session:1", "user", "suhaan")
	store.HSet("session:1", "role", "admin")
	store.HSet("session:1", "active", "true")

	hash, ok, err := store.HGetAll("session:1")
	if err != nil {
		t.Fatalf("HGetAll() error = %v", err)
	}
	if !ok {
		t.Error("HGetAll() ok = false, want true")
	}
	if len(hash) != 3 {
		t.Errorf("HGetAll() returned %d fields, want 3", len(hash))
	}
	if hash["user"] != "suhaan" {
		t.Errorf("HGetAll()[user] = %v, want suhaan", hash["user"])
	}
}

func TestHDel(t *testing.T) {
	store := NewKVStore()

	store.HSet("hash", "f1", "v1")
	store.HSet("hash", "f2", "v2")

	// Delete f1
	deleted, err := store.HDel("hash", "f1")
	if err != nil {
		t.Fatalf("HDel() error = %v", err)
	}
	if !deleted {
		t.Error("HDel() = false, want true")
	}

	// f1 should be gone
	_, ok, _ := store.HGet("hash", "f1")
	if ok {
		t.Error("Field should not exist after HDel")
	}

	// f2 should still exist
	_, ok, _ = store.HGet("hash", "f2")
	if !ok {
		t.Error("Other fields should still exist")
	}
}

func TestHSetWrongType(t *testing.T) {
	store := NewKVStore()

	// Set a string key
	store.Set("string_key", "value", 0)

	// Try to HSet on it - should error
	err := store.HSet("string_key", "field", "value")
	if err == nil {
		t.Error("HSet() should error on wrong type")
	}
}
