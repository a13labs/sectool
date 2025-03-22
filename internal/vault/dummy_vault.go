package vault

import (
	"errors"
	"sync"

	"github.com/a13labs/sectool/internal/crypto"
)

type DummyVault struct {
	VaultProvider
	data   map[string]string
	mu     sync.RWMutex
	backup bool
}

func NewDummyVault() *DummyVault {
	return &DummyVault{
		data: make(map[string]string),
	}
}

func (v *DummyVault) Initialize() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data = make(map[string]string)
	return nil
}

func (v *DummyVault) VaultHasKey(key string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	_, exists := v.data[key]
	return exists
}

func (v *DummyVault) VaultGetValue(key string) (string, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	value, exists := v.data[key]
	if !exists {
		return "", errors.New("key not found in vault")
	}
	return value, nil
}

func (v *DummyVault) VaultListKeys() []string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	keys := make([]string, 0, len(v.data))
	for key := range v.data {
		keys = append(keys, key)
	}
	return keys
}

func (v *DummyVault) VaultSetValue(key, value string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.data[key] = value
	return nil
}

func (v *DummyVault) VaultDelKey(key string) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, exists := v.data[key]; !exists {
		return errors.New("key not found in vault")
	}
	delete(v.data, key)
	return nil
}

func (v *DummyVault) VaultEnableBackup(value bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.backup = value
}

func (v *DummyVault) SetSensitiveStrings(kv *crypto.SecureKVStore) {
}

func (v *DummyVault) Lock() error {
	return nil
}

func (v *DummyVault) Unlock() error {
	return nil
}

func (v *DummyVault) VaultGetMultipleValues(keys []string, kv *crypto.SecureKVStore) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	for _, key := range keys {
		if value, exists := v.data[key]; exists {
			err := kv.Put(key, value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
