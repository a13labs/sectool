package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ProviderType represents the type of vault provider
type ProviderType string

const (
	FileProvider          ProviderType = "file"
	BitwardenProvider     ProviderType = "bitwarden"
	ObjectStorageProvider ProviderType = "object_storage"
)

// Config represents the configuration structure
type Config struct {
	Provider           ProviderType         `json:"provider"`
	FileVault          *FileConfig          `json:"file,omitempty"`
	BitwardenVault     *BitwardenConfig     `json:"bitwarden,omitempty"`
	ObjectStorageVault *ObjectStorageConfig `json:"object_storage,omitempty"`
	SSHPasswordKey     string               `json:"ssh_password_key,omitempty"`
}

// FileConfig represents the configuration for the file provider
type FileConfig struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path"`
}

// BitwardenConfig represents the configuration for the Bitwarden provider
type BitwardenConfig struct {
	APIURL         string `json:"api_url,omitempty"`
	IdentityURL    string `json:"identity_url,omitempty"`
	AccessToken    string `json:"access_token,omitempty"`
	OrganizationId string `json:"organization,omitempty"`
	ProjectId      string `json:"project,omitempty"`
}

type ObjectStorageConfig struct {
	Region   string `json:"region"`
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Backup   bool   `json:"backup"`
}

var (
	// DefaultConfigFile is the default configuration file
	defaultConfig = Config{
		Provider: "file",
		FileVault: &FileConfig{
			Path: "repository.vault",
		},
	}
)

// ReadConfig reads the configuration from a JSON file
func ReadConfig(config_file string) (*Config, error) {

	if config_file == "" {
		var exist bool
		config_file, exist = os.LookupEnv("SECTOOL_CONFIG_FILE")
		if !exist {
			config_file = "sectool.json"
		}
	}

	if _, err := os.Stat(config_file); os.IsNotExist(err) {
		return &defaultConfig, nil
	}

	file, err := os.Open(config_file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}
