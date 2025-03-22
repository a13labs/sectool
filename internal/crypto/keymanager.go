package crypto

import (
	"crypto/rand"
	"errors"
	"sync"
)

// KeyManager stores encryption keys externally
type KeyManager struct {
	keys map[string][]byte
	mu   sync.RWMutex
}

// NewKeyManager initializes an external key manager
func NewKeyManager() *KeyManager {
	return &KeyManager{keys: make(map[string][]byte)}
}

// GenerateKey creates a new AES-256 key
func (km *KeyManager) GenerateKey(id string) ([]byte, error) {
	key := make([]byte, 32) // AES-256 key
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	km.mu.Lock()
	km.keys[id] = key
	km.mu.Unlock()
	return key, nil
}

// GetKey retrieves a key from the key manager
func (km *KeyManager) GetKey(id string) ([]byte, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	key, exists := km.keys[id]
	if !exists {
		return nil, errors.New("key not found")
	}
	return key, nil
}

// DeleteKey removes a key from the key manager
func (km *KeyManager) DeleteKey(id string) {
	km.mu.Lock()
	delete(km.keys, id)
	km.mu.Unlock()
}
