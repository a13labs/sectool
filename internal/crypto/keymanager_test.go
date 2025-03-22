package crypto

import (
	"testing"
)

func TestKeyManager_GenerateKey(t *testing.T) {
	km := NewKeyManager()
	id := "test-key"

	key, err := km.GenerateKey(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(key) != 32 {
		t.Errorf("expected key length 32, got %d", len(key))
	}

	retrievedKey, err := km.GetKey(id)
	if err != nil {
		t.Fatalf("expected no error retrieving key, got %v", err)
	}

	if string(key) != string(retrievedKey) {
		t.Errorf("expected retrieved key to match generated key")
	}
}

func TestKeyManager_GetKey_NotFound(t *testing.T) {
	km := NewKeyManager()
	_, err := km.GetKey("non-existent-key")
	if err == nil {
		t.Errorf("expected error for non-existent key, got nil")
	}
}

func TestKeyManager_DeleteKey(t *testing.T) {
	km := NewKeyManager()
	id := "test-key"

	_, err := km.GenerateKey(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	km.DeleteKey(id)

	_, err = km.GetKey(id)
	if err == nil {
		t.Errorf("expected error for deleted key, got nil")
	}
}
