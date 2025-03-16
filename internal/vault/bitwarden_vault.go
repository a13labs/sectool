package vault

import (
	"errors"
	"log"
	"os"

	"github.com/a13labs/sectool/internal/config"
	"github.com/bitwarden/sdk-go"
)

// BitwardenVault represents a Bitwarden vault provider.
type BitwardenVault struct {
	VaultProvider
	apiURL         string
	identityURL    string
	accessToken    string
	organizationId string
	projectId      string
}

// NewBitwardenVault creates a new BitwardenVault instance.
func NewBitwardenVault(config *config.BitwardenConfig) (*BitwardenVault, error) {

	if config == nil {
		return nil, errors.New("bitwarden configuration is nil")
	}

	v := BitwardenVault{
		apiURL:         config.APIURL,
		identityURL:    config.IdentityURL,
		accessToken:    config.AccessToken,
		organizationId: config.OrganizationId,
		projectId:      config.ProjectId,
	}

	var err error
	var exist bool
	if v.apiURL == "" {
		v.apiURL, exist = os.LookupEnv("BW_API_URL")
		if !exist {
			v.apiURL = "https://api.bitwarden.com"
		}
	}

	if v.identityURL == "" {
		v.identityURL, exist = os.LookupEnv("BW_IDENTITY_URL")
		if !exist {
			v.identityURL = "https://identity.bitwarden.com"
		}
	}

	if v.accessToken == "" {
		v.accessToken, exist = os.LookupEnv("BW_ACCESS_TOKEN")
		if !exist {
			log.Printf("Error parsing access token: %v", err)
			return nil, errors.New("access token is not defined")
		}
	}

	if v.projectId == "" {
		v.projectId, exist = os.LookupEnv("BW_PROJECT_ID")
		if !exist {
			log.Printf("Error parsing project ID: %v", err)
			return nil, errors.New("project ID is not defined")
		}
	}

	if v.organizationId == "" {
		v.organizationId, exist = os.LookupEnv("BW_ORGANIZATION_ID")
		if !exist {
			log.Printf("Error parsing organization ID: %v", err)
			return nil, errors.New("organization ID is not defined")
		}
	}

	return &v, nil
}

// VaultListKeys lists all keys in the Bitwarden vault.
func (v *BitwardenVault) VaultListKeys() []string {

	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return []string{}
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return []string{}
	}

	secretIdentifiers, err := client.Secrets().List(v.organizationId)
	if err != nil {
		return []string{}
	}

	// Get secrets with a list of IDs that belong to the specified project
	secretKeys := make([]string, 0, len(secretIdentifiers.Data))
	for _, identifier := range secretIdentifiers.Data {
		if identifier.OrganizationID != v.organizationId {
			continue
		}

		secret, err := client.Secrets().Get(identifier.ID)
		if err != nil {
			continue
		}

		if *secret.ProjectID != v.projectId {
			continue
		}

		secretKeys = append(secretKeys, identifier.Key)
	}

	return secretKeys
}

// VaultSetValue sets the value of a key in the Bitwarden vault.
func (v *BitwardenVault) VaultSetValue(key, value string) error {
	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return nil
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return err
	}

	_, err = client.Secrets().Create(key, value, "Sectool managed secret", v.organizationId, []string{v.projectId})
	if err != nil {
		return err
	}

	return nil
}

// VaultGetValue returns the value of a key from the Bitwarden vault.
func (v *BitwardenVault) VaultGetValue(key string) (string, error) {
	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return "", nil
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return "", err
	}

	secretIdentifiers, err := client.Secrets().List(v.organizationId)
	if err != nil {
		return "", err
	}

	// Get secrets with a list of IDs
	secretId := ""
	for _, identifier := range secretIdentifiers.Data {
		if identifier.Key == key {

			// Check if the secret belongs to the specified project
			secret, err := client.Secrets().Get(identifier.ID)
			if err != nil {
				continue
			}

			if *secret.ProjectID != v.projectId {
				continue
			}

			secretId = identifier.ID
			break
		}
	}

	if secretId == "" {
		return "", errors.New("key not found in vault")
	}

	secret, err := client.Secrets().Get(secretId)
	if err != nil {
		return "", err
	}

	if *secret.ProjectID != v.projectId || secret.OrganizationID != v.organizationId {
		return "",
			errors.New("secret does not belong to the specified project")
	}

	return secret.Value, nil
}

// VaultDelKey deletes a key from the Bitwarden vault.
func (v *BitwardenVault) VaultDelKey(key string) error {

	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return nil
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return err
	}

	secretIdentifiers, err := client.Secrets().List(v.organizationId)
	if err != nil {
		return err
	}

	// Get secrets with a list of IDs
	secretId := ""
	for _, identifier := range secretIdentifiers.Data {
		if identifier.Key == key {
			secretId = identifier.ID
			break
		}
	}

	if secretId == "" {
		return errors.New("key not found in vault")
	}

	_, err = client.Secrets().Delete([]string{secretId})
	if err != nil {
		return err
	}

	return nil
}

// VaultHasKey checks if the Bitwarden vault contains the specified key.
func (v *BitwardenVault) VaultHasKey(key string) bool {
	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return false
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return false
	}

	secretIdentifiers, err := client.Secrets().List(v.organizationId)
	if err != nil {
		return false
	}

	// Get secrets with a list of IDs
	for _, identifier := range secretIdentifiers.Data {
		if identifier.Key == key {
			return true
		}
	}

	return false
}

// VaultEnableBackup is not applicable for BitwardenVault.
func (v *BitwardenVault) VaultEnableBackup(value bool) {
	// No-op for BitwardenVault
}

// GetSensitiveStrings returns sensitive strings from the Bitwarden vault.
func (v *BitwardenVault) GetSensitiveStrings() []string {
	return []string{v.accessToken, v.organizationId, v.projectId}
}

func (v *BitwardenVault) Lock() error {
	return nil
}

func (v *BitwardenVault) Unlock() error {
	return nil
}

// VaultGetMultipleValues returns multiple values from the Bitwarden vault.
func (v *BitwardenVault) VaultGetMultipleValues(keys []string) (map[string]string, error) {

	client, err := sdk.NewBitwardenClient(&v.apiURL, &v.identityURL)
	if err != nil {
		log.Printf("Error creating Bitwarden client: %v", err)
		return nil, err
	}
	defer client.Close()

	err = client.AccessTokenLogin(v.accessToken, nil)
	if err != nil {
		log.Printf("Error logging in with access token: %v", err)
		return nil, err
	}

	secretIdentifiers, err := client.Secrets().List(v.organizationId)
	if err != nil {
		return nil, err
	}

	// Get secrets with a list of IDs that belong to the specified project
	secretKeysMap := make(map[string]string, len(secretIdentifiers.Data))
	for _, identifier := range secretIdentifiers.Data {
		if identifier.OrganizationID != v.organizationId {
			continue
		}

		secret, err := client.Secrets().Get(identifier.ID)
		if err != nil {
			continue
		}

		if *secret.ProjectID != v.projectId {
			continue
		}

		secretKeysMap[identifier.Key] = secret.ID
	}

	values := make(map[string]string)
	for _, key := range keys {

		id, ok := secretKeysMap[key]
		if !ok {
			continue
		}

		secret, err := client.Secrets().Get(id)
		if err != nil {
			continue
		}

		if *secret.ProjectID != v.projectId {
			continue
		}

		values[key] = secret.Value
	}

	return values, nil
}
