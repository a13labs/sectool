// vault.go
package vault

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a13labs/sectool/internal/crypto"
)

// Vault represents a secure key-value store.
type Vault struct {
	filePath string
	key      []byte
	backup   bool
}

// NewVault creates a new Vault instance.
func NewVault(filePath string, key []byte) *Vault {
	return &Vault{
		filePath: filePath,
		key:      key,
		backup:   false,
	}
}

// readVault reads the vault file contents.
func (v *Vault) readVault() (string, error) {
	if !v.vaultFileExists() {
		f, err := os.Create(v.filePath)
		if err != nil {
			return "", errors.New("vault file does not exist")
		}
		f.Close()
	}

	decryptedContents, err := crypto.DecryptFromFile(v.filePath, v.key)
	if err != nil {
		return "", err
	}

	return decryptedContents, nil
}

// writeVault writes encrypted data to the vault file.
func (v *Vault) writeVault(contents string) error {

	if v.backup {
		// Create a backup of the existing vault
		backupName := v.vaultBackupName()
		err := v.backupVault(backupName)
		if err != nil {
			return err
		}
	}

	// Encrypt the data and write to the vault file
	err := crypto.EncryptToFile(contents, v.filePath, v.key, false)
	if err != nil {
		return err
	}

	return nil
}

// backupVault creates a backup of the vault file.
func (v *Vault) backupVault(backupName string) error {
	contents, err := os.ReadFile(v.filePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(backupName, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

// vaultBackupName generates a backup filename with a timestamp.
func (v *Vault) vaultBackupName() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s", v.filePath, timestamp)
}

// vaultFileExists checks if the vault file exists.
func (v *Vault) vaultFileExists() bool {
	_, err := os.Stat(v.filePath)
	return !os.IsNotExist(err)
}

// Initialize creates a new vault file if it doesn't exist.
func (v *Vault) Initialize() error {
	if v.vaultFileExists() {
		return nil
	}

	// Create an empty vault file
	err := os.WriteFile(v.filePath, []byte(""), 0600)
	if err != nil {
		return err
	}

	return nil
}

// VaultHasKey checks if the vault contains the specified key.
func (v *Vault) VaultHasKey(key string) bool {
	contents, err := v.readVault()
	if err != nil {
		return false
	}

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return true
		}
	}

	return false
}

// VaultGetValue returns the value of a key from the vault.
func (v *Vault) VaultGetValue(key string) (string, error) {
	contents, err := v.readVault()
	if err != nil {
		return "", err
	}

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			return parts[1], nil
		}
	}

	return "", errors.New("key not found in vault")
}

// VaultListKeys lists all keys in the vault.
func (v *Vault) VaultListKeys() []string {
	contents, err := v.readVault()
	if err != nil {
		return []string{}
	}

	var keys []string
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			keys = append(keys, parts[0])
		}
	}

	return keys
}

// VaultSetValue sets the value of a key in the vault.
func (v *Vault) VaultSetValue(key, value string) error {
	contents, err := v.readVault()
	if err != nil {
		return err
	}

	lines := strings.Split(contents, "\n")
	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			// Key already exists, update the value
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updatedContents := strings.Join(lines, "\n")
			return v.writeVault(updatedContents)
		}
	}

	// Key doesn't exist, add a new key-value pair
	lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	updatedContents := strings.Join(lines, "\n")
	return v.writeVault(updatedContents)
}

// VaultDelKey deletes a key from the vault.
func (v *Vault) VaultDelKey(key string) error {
	contents, err := v.readVault()
	if err != nil {
		return err
	}

	var updatedLines []string
	lines := strings.Split(contents, "\n")
	keyFound := false
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			keyFound = true
			continue
		}
		updatedLines = append(updatedLines, line)
	}

	if !keyFound {
		return errors.New("key not found in vault")
	}

	updatedContents := strings.Join(updatedLines, "\n")
	return v.writeVault(updatedContents)
}

func (v *Vault) VaultEnableBackup(value bool) {
	v.backup = value
}
