package crypto

import (
	"testing"
)

func TestSecureKVStore_PutAndGet(t *testing.T) {
	mk := NewKeyManager()
	store := NewSecureKVStore(mk)

	err := store.Put("key1", "value1")
	if err != nil {
		t.Fatalf("failed to put value: %v", err)
	}

	value, err := store.Get("key1")
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}

	if value != "value1" {
		t.Errorf("expected value1, got %s", value)
	}
}

func TestSecureKVStore_GetNonExistentKey(t *testing.T) {
	mk := NewKeyManager()
	store := NewSecureKVStore(mk)

	_, err := store.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent key, got nil")
	}
}

func TestSecureKVStore_Delete(t *testing.T) {
	mk := NewKeyManager()
	store := NewSecureKVStore(mk)

	err := store.Put("key1", "value1")
	if err != nil {
		t.Fatalf("failed to put value: %v", err)
	}

	store.Delete("key1")

	_, err = store.Get("key1")
	if err == nil {
		t.Fatal("expected error for deleted key, got nil")
	}
}

func TestSecureKVStore_Clear(t *testing.T) {
	mk := NewKeyManager()
	store := NewSecureKVStore(mk)

	err := store.Put("key1", "value1")
	if err != nil {
		t.Fatalf("failed to put value: %v", err)
	}

	err = store.Put("key2", "value2")
	if err != nil {
		t.Fatalf("failed to put value: %v", err)
	}

	store.Clear()

	_, err = store.Get("key1")
	if err == nil {
		t.Fatal("expected error for cleared key, got nil")
	}

	_, err = store.Get("key2")
	if err == nil {
		t.Fatal("expected error for cleared key, got nil")
	}
}
