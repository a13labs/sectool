package vault

import (
	"errors"

	"github.com/a13labs/sectool/internal/config"
)

// VaultProvider defines the interface for a vault provider.
type VaultProvider interface {
	VaultListKeys() []string
	VaultSetValue(key, value string) error
	VaultGetValue(key string) (string, error)
	VaultDelKey(key string) error
	VaultHasKey(key string) bool
	VaultEnableBackup(value bool)
	GetSensitiveStrings() []string
	VaultGetMultipleValues(keys []string) (map[string]string, error)
	Lock() error
	Unlock() error
}

// NewVaultProvider creates a new vault provider based on the configuration.
func NewVaultProvider(cfg config.Config) (VaultProvider, error) {

	switch cfg.Provider {
	case config.FileProvider:
		return NewFileVault(cfg.FileVault)
	case config.BitwardenProvider:
		return NewBitwardenVault(cfg.BitwardenVault)
	default:
		return nil, errors.New("unsupported vault provider")
	}
}
