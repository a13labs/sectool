package vault

import (
	"testing"

	"github.com/a13labs/sectool/internal/crypto"
)

func TestNewDummyVault(t *testing.T) {
	vault := NewDummyVault()
	if vault == nil {
		t.Fatal("Expected non-nil DummyVault instance")
	}
}

func TestDummyVault_Initialize(t *testing.T) {
	vault := NewDummyVault()
	err := vault.Initialize()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(vault.data) != 0 {
		t.Fatal("Expected empty vault data after initialization")
	}
}

func TestDummyVault_VaultHasKey(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultSetValue("testKey", "testValue")
	if !vault.VaultHasKey("testKey") {
		t.Fatal("Expected key to be present in vault")
	}
	if vault.VaultHasKey("nonExistentKey") {
		t.Fatal("Expected key to be absent in vault")
	}
}

func TestDummyVault_VaultGetValue(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultSetValue("testKey", "testValue")
	value, err := vault.VaultGetValue("testKey")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "testValue" {
		t.Fatalf("Expected value 'testValue', got %v", value)
	}
	_, err = vault.VaultGetValue("nonExistentKey")
	if err == nil {
		t.Fatal("Expected error for non-existent key")
	}
}

func TestDummyVault_VaultListKeys(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultSetValue("testKey1", "testValue1")
	vault.VaultSetValue("testKey2", "testValue2")
	keys := vault.VaultListKeys()
	if len(keys) != 2 {
		t.Fatalf("Expected 2 keys, got %d", len(keys))
	}
}

func TestDummyVault_VaultSetValue(t *testing.T) {
	vault := NewDummyVault()
	err := vault.VaultSetValue("testKey", "testValue")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !vault.VaultHasKey("testKey") {
		t.Fatal("Expected key to be present in vault")
	}
}

func TestDummyVault_VaultDelKey(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultSetValue("testKey", "testValue")
	err := vault.VaultDelKey("testKey")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if vault.VaultHasKey("testKey") {
		t.Fatal("Expected key to be absent in vault")
	}
	err = vault.VaultDelKey("nonExistentKey")
	if err == nil {
		t.Fatal("Expected error for non-existent key")
	}
}

func TestDummyVault_VaultEnableBackup(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultEnableBackup(true)
	if !vault.backup {
		t.Fatal("Expected backup to be enabled")
	}
	vault.VaultEnableBackup(false)
	if vault.backup {
		t.Fatal("Expected backup to be disabled")
	}
}

func TestDummyVault_VaultGetMultipleValues(t *testing.T) {
	vault := NewDummyVault()
	vault.VaultSetValue("key1", "value1")
	vault.VaultSetValue("key2", "value2")
	keys := []string{"key1", "key2", "key3"}
	km := crypto.NewKeyManager()
	kv := crypto.NewSecureKVStore(km)
	err := vault.VaultGetMultipleValues(keys, kv)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if kv.Len() != 2 {
		t.Fatalf("Expected 2 values, got %d", kv.Len())
	}
	if value, err := kv.Get("key1"); err != nil || value != "value1" {
		t.Fatalf("Expected value 'value1', got %v", value)
	}
	if value, err := kv.Get("key2"); err != nil || value != "value2" {
		t.Fatalf("Expected value 'value2', got %v", value)
	}
}
