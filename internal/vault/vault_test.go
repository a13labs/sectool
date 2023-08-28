package vault

import (
	"errors"
	"os"
	"testing"
)

func TestVault(t *testing.T) {
	key := []byte("mysecretkey")
	vault_path := "testdata/repository.vault"
	vault := NewVault(vault_path, key)

	defer func() {
		_ = os.Remove(vault_path)
	}()

	if len(vault.VaultListKeys()) != 0 {
		t.Error(errors.New("VaultListKeys != 0"))
	}

	if vault.VaultHasKey("KEY1") {
		t.Error(errors.New("Key found"))
	}

	if err := vault.VaultSetValue("KEY1", "VALUE1"); err != nil {
		t.Error(err)
	}

	if len(vault.VaultListKeys()) != 1 {
		t.Error(errors.New("VaultListKeys != 1"))
	}

	if !vault.VaultHasKey("KEY1") {
		t.Error(errors.New("Key found"))
	}

	value, err := vault.VaultGetValue("KEY1")
	if err != nil {
		t.Error(err)
	}

	if value != "VALUE1" {
		t.Error(errors.New("Key found"))
	}

	if err := vault.VaultSetValue("KEY2", "VALUE2"); err != nil {
		t.Error(err)
	}

	if len(vault.VaultListKeys()) != 2 {
		t.Error(errors.New("VaultListKeys != 2"))
	}

	if !vault.VaultHasKey("KEY2") {
		t.Error(errors.New("Key found"))
	}

	value, err = vault.VaultGetValue("KEY2")
	if err != nil {
		t.Error(err)
	}

	if value != "VALUE2" {
		t.Error(errors.New("Key found"))
	}

	vault.VaultDelKey("KEY1")

	if len(vault.VaultListKeys()) != 1 {
		t.Error(errors.New("VaultListKeys != 1"))
	}

	if vault.VaultHasKey("KEY1") {
		t.Error(errors.New("Key found"))
	}

	if !vault.VaultHasKey("KEY2") {
		t.Error(errors.New("Key not found"))
	}
}

// TestMain runs before all tests and can be used for setup or teardown tasks
func TestMain(m *testing.M) {
	// Setup code (if any)

	// Run tests
	exitCode := m.Run()

	// Teardown code (if any)

	// Exit with the appropriate code
	os.Exit(exitCode)
}
