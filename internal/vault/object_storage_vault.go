// filepath: /home/alexandre/Projects/sectool/internal/vault/object_storage_vault.go
package vault

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/crypto"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ObjectStorageVault represents a secure key-value store stored in an S3 bucket.
type ObjectStorageVault struct {
	VaultProvider
	client   *s3.Client
	bucket   string
	key      []byte
	fileName string
	backup   bool
}

// NewObjectStorageVault creates a new ObjectStorageVault instance.
func NewObjectStorageVault(c *config.ObjectStorageConfig) (*ObjectStorageVault, error) {
	if c == nil {
		return nil, errors.New("object storage configuration is nil")
	}

	vaultKey := c.Key
	if vaultKey == "" {
		vaultKey, _ = os.LookupEnv("FILE_VAULT_KEY")
		if vaultKey == "" {
			return nil, errors.New("FILE_VAULT_KEY is not defined")
		}
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(c.Region),
		awsconfig.LoadOptionsFunc(func(o *awsconfig.LoadOptions) error {
			o.EndpointResolver = aws.EndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: c.Endpoint}, nil
			}))
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsConfig)

	return &ObjectStorageVault{
		client:   client,
		bucket:   c.Bucket,
		key:      []byte(vaultKey),
		fileName: "repository.vault",
		backup:   false,
	}, nil
}

// readVault reads the vault file contents from the S3 bucket.
func (v *ObjectStorageVault) readVault() (string, error) {
	output, err := v.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(v.bucket),
		Key:    aws.String(v.fileName),
	})
	if err != nil {
		if isNotFoundError(err) {
			return "", nil
		}
		return "", err
	}
	defer output.Body.Close()

	decryptedContents, err := crypto.DecryptFromReader(output.Body, v.key)
	if err != nil {
		return "", err
	}

	return decryptedContents, nil
}

// writeVault writes encrypted data to the vault file in the S3 bucket.
func (v *ObjectStorageVault) writeVault(contents string) error {
	if v.backup {
		backupName := v.vaultBackupName()
		err := v.backupVault(backupName)
		if err != nil {
			return err
		}
	}

	encryptedData, err := crypto.EncryptToBytes(contents, v.key)
	if err != nil {
		return err
	}

	_, err = v.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(v.bucket),
		Key:    aws.String(v.fileName),
		Body:   bytes.NewReader(encryptedData),
	})
	return err
}

// backupVault creates a backup of the vault file in the S3 bucket.
func (v *ObjectStorageVault) backupVault(backupName string) error {
	output, err := v.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(v.bucket),
		Key:    aws.String(v.fileName),
	})
	if err != nil {
		return err
	}
	defer output.Body.Close()

	_, err = v.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(v.bucket),
		Key:    aws.String(backupName),
		Body:   output.Body,
	})
	return err
}

// vaultBackupName generates a backup filename with a timestamp.
func (v *ObjectStorageVault) vaultBackupName() string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s", v.fileName, timestamp)
}

// vaultFileExists checks if the vault file exists in the S3 bucket.
func (v *ObjectStorageVault) vaultFileExists() bool {
	_, err := v.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(v.bucket),
		Key:    aws.String(v.fileName),
	})
	return err == nil
}

// Initialize creates a new vault file in the S3 bucket if it doesn't exist.
func (v *ObjectStorageVault) Initialize() error {
	if v.vaultFileExists() {
		return nil
	}

	return v.writeVault("")
}

// VaultHasKey checks if the vault contains the specified key.
func (v *ObjectStorageVault) VaultHasKey(key string) bool {
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
func (v *ObjectStorageVault) VaultGetValue(key string) (string, error) {
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
func (v *ObjectStorageVault) VaultListKeys() []string {
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
func (v *ObjectStorageVault) VaultSetValue(key, value string) error {
	contents, err := v.readVault()
	if err != nil {
		return err
	}

	lines := strings.Split(contents, "\n")
	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
			updatedContents := strings.Join(lines, "\n")
			return v.writeVault(updatedContents)
		}
	}

	lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	updatedContents := strings.Join(lines, "\n")
	return v.writeVault(updatedContents)
}

// VaultDelKey deletes a key from the vault.
func (v *ObjectStorageVault) VaultDelKey(key string) error {
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
func (v *ObjectStorageVault) VaultEnableBackup(value bool) {
	v.backup = value
}

// GetSensitiveStrings returns the sensitive strings in the vault.
func (v *ObjectStorageVault) SetSensitiveStrings(kv *crypto.SecureKVStore) {
	kv.Put("SECTOOL_OS_SENSITIVE_1", string(v.key))
}

// isNotFoundError checks if an error is an S3 not found error.
func isNotFoundError(err error) bool {
	var notFoundErr *types.NoSuchKey
	return errors.As(err, &notFoundErr)
}

// VaultGetMultipleValues returns the values of multiple keys from the vault.
func (v *ObjectStorageVault) VaultGetMultipleValues(keys []string, kv *crypto.SecureKVStore) error {

	contents, err := v.readVault()
	if err != nil {
		return err
	}

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			for _, key := range keys {
				if parts[0] == key {
					kv.Put(parts[0], parts[1])
				}
			}
		}
	}

	return nil
}
