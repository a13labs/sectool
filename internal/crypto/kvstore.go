package crypto

import (
	"errors"
	"sync"
)

// SecureKVStore is an in-memory encrypted key-value store
type SecureKVStore struct {
	store      map[string]string
	keyManager *KeyManager
	mu         sync.RWMutex
}

// NewSecureKVStore initializes a new secure in-memory store
func NewSecureKVStore(km *KeyManager) *SecureKVStore {
	return &SecureKVStore{
		store:      make(map[string]string),
		keyManager: km,
	}
}

// Put securely stores an encrypted key-value pair in memory
func (s *SecureKVStore) Put(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	encryptionKey, err := s.keyManager.GenerateKey(key)
	if err != nil {
		return err
	}

	encryptedValue, err := Encrypt(value, encryptionKey)
	if err != nil {
		return err
	}

	s.store[key] = encryptedValue
	return nil
}

// Get retrieves and decrypts a value from memory
func (s *SecureKVStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	encryptionKey, err := s.keyManager.GetKey(key)
	if err != nil {
		return "", err
	}

	encryptedValue, exists := s.store[key]
	if !exists {
		return "", errors.New("key not found")
	}

	return Decrypt(encryptedValue, encryptionKey)
}

// Delete removes a key from the store
func (s *SecureKVStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}

// Clear removes all values from the store
func (s *SecureKVStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store = make(map[string]string)
}

func (s *SecureKVStore) ListKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.store))
	for key := range s.store {
		keys = append(keys, key)
	}
	return keys
}

func (s *SecureKVStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.store)
}

func (s *SecureKVStore) MatchValue(value string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, encryptedValue := range s.store {

		encryptionKey, err := s.keyManager.GetKey(key)
		if err != nil {
			continue
		}

		decryptedValue, err := Decrypt(encryptedValue, encryptionKey)
		if err != nil {
			continue
		}

		if decryptedValue == value {
			return true
		}

	}

	return false
}
