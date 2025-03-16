package vault

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/crypto"
)

// FileVault represents a secure key-value store.
type FileVault struct {
	VaultProvider
	path   string
	key    []byte
	backup bool
}

// NewVault creates a new FileVault instance.
func NewFileVault(config *config.FileConfig) (*FileVault, error) {

	if config == nil {
		fmt.Println("File configuration is nil")
		return nil, errors.New("file configuration is nil")
	}

	key := config.Key
	if key == "" {
		key, _ = os.LookupEnv("FILE_VAULT_KEY")
		if key == "" {
			fmt.Println("FILE_VAULT_KEY it's not defined, aborting.")
			return nil, errors.New("file vault key is not defined")
		}
	}

	path := config.Path
	if path == "" {
		path, _ := os.LookupEnv("FILE_VAULT_PATH")
		if path == "" {
			path = "repository.vault"
		}
	}

	return &FileVault{
		path:   path,
		key:    []byte(key),
		backup: false,
	}, nil
}

// readVault reads the vault file contents.
func (v *FileVault) readVault() (string, error) {
	if !v.vaultFileExists() {
		f, err := os.Create(v.path)
		if err != nil {
			return "", errors.New("vault file does not exist")
		}
		f.Close()
	}

	decryptedContents, err := crypto.DecryptFromFile(v.path, v.key)
	if err != nil {
		return "", err
	}

	return decryptedContents, nil
}

// writeVault writes encrypted data to the vault file.
func (v *FileVault) writeVault(contents string) error {

	if v.backup {
		// Create a backup of the existing vault
		backupName := v.vaultBackupName()
		err := v.backupVault(backupName)
		if err != nil {
			return err
		}
	}

	// Encrypt the data and write to the vault file
	err := crypto.EncryptToFile(contents, v.path, v.key, false)
	if err != nil {
		return err
	}

	return nil
}

// backupVault creates a backup of the vault file.
func (v *FileVault) backupVault(backupName string) error {
	contents, err := os.ReadFile(v.path)
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
func (v *FileVault) vaultBackupName() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s", v.path, timestamp)
}

// vaultFileExists checks if the vault file exists.
func (v *FileVault) vaultFileExists() bool {
	_, err := os.Stat(v.path)
	return !os.IsNotExist(err)
}

// Initialize creates a new vault file if it doesn't exist.
func (v *FileVault) Initialize() error {
	if v.vaultFileExists() {
		return nil
	}

	// Create an empty vault file
	err := os.WriteFile(v.path, []byte(""), 0600)
	if err != nil {
		return err
	}

	return nil
}

// VaultHasKey checks if the vault contains the specified key.
func (v *FileVault) VaultHasKey(key string) bool {
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
func (v *FileVault) VaultGetValue(key string) (string, error) {
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
func (v *FileVault) VaultListKeys() []string {
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
func (v *FileVault) VaultSetValue(key, value string) error {
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
func (v *FileVault) VaultDelKey(key string) error {
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

// VaultEnableBackup enables or disables vault backups.
func (v *FileVault) VaultEnableBackup(value bool) {
	v.backup = value
}

// GetSensitiveStrings returns the sensitive strings in the vault.
func (v *FileVault) GetSensitiveStrings() []string {
	return []string{string(v.key)}
}

// Lock encrypts the vault file.
func (v *FileVault) Lock() error {

	unlocked_vault := v.path + ".unlocked"
	locked_vault := v.path

	_, err := os.Stat(unlocked_vault)
	if os.IsNotExist(err) {
		return nil
	}

	err = crypto.EncryptFile(unlocked_vault, locked_vault, []byte(v.key))
	if err != nil {
		return err
	}

	os.Remove(unlocked_vault)

	return nil
}

// Unlock decrypts the vault file.
func (v *FileVault) Unlock() error {

	unlocked_vault := v.path + ".unlocked"
	locked_vault := v.path

	_, err := os.Stat(locked_vault)
	if os.IsNotExist(err) {
		return nil
	}

	err = crypto.DecryptFile(locked_vault, unlocked_vault, []byte(v.key))
	if err != nil {
		return err
	}

	return nil
}

// VaultGetMultipleValues returns the values of multiple keys from the vault.
func (v *FileVault) VaultGetMultipleValues(keys []string) (map[string]string, error) {
	contents, err := v.readVault()
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			for _, key := range keys {
				if parts[0] == key {
					values[key] = parts[1]
				}
			}
		}
	}

	return values, nil
}
