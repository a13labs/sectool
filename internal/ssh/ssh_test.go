package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyPair(t *testing.T) {
	privateKeyECDSA, publicKeyECDSA, err := GenerateKeyPair(ECDSA, 256, "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, privateKeyECDSA)
	assert.NotEmpty(t, publicKeyECDSA)

	privateKeyRSA, publicKeyRSA, err := GenerateKeyPair(RSA, 2048, "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, privateKeyRSA)
	assert.NotEmpty(t, publicKeyRSA)
}

func TestGetPublicKeyFromPrivateKey(t *testing.T) {
	privateKeyECDSA, _, err := generateECDSAKey(256)
	assert.NoError(t, err)
	pubKeyECDSA, err := getPublicKeyFromPrivateKey(privateKeyECDSA)
	assert.NoError(t, err)
	assert.NotNil(t, pubKeyECDSA)

	privateKeyRSA, _, err := generateRSAKey(2048)
	assert.NoError(t, err)
	pubKeyRSA, err := getPublicKeyFromPrivateKey(privateKeyRSA)
	assert.NoError(t, err)
	assert.NotNil(t, pubKeyRSA)
}

func TestIsValidPrivateKey(t *testing.T) {
	privateKey, _, err := GenerateKeyPair(ECDSA, 256, "password")
	assert.NoError(t, err)
	assert.True(t, IsValidPrivateKey(privateKey))

	invalidPrivateKey := "invalid-private-key"
	assert.False(t, IsValidPrivateKey(invalidPrivateKey))
}

func TestGetPublicKey(t *testing.T) {
	privateKey, _, err := GenerateKeyPair(ECDSA, 256, "password")
	assert.NoError(t, err)

	publicKey, err := GetPublicKey(privateKey, "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, publicKey)
}
