package ssh

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Algorithm represents the SSH key algorithm.
type Algorithm int

const (
	ECDSA Algorithm = iota
	RSA
)

// GenerateKeyPair generates a password protected SSH key pair.
func GenerateKeyPair(algo Algorithm, keyLength int, password string) (privateKey, publicKey string, err error) {
	var privateKeyBytes []byte
	var privKey *pem.Block

	switch algo {
	case ECDSA:
		privKey, privateKeyBytes, err = generateECDSAKey(keyLength)
		if err != nil {
			return "", "", err
		}
	case RSA:
		privKey, privateKeyBytes, err = generateRSAKey(keyLength)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("unsupported algorithm")
	}

	encryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader, privKey.Type, privateKeyBytes, []byte(password), x509.PEMCipherAES256)
	if err != nil {
		return "", "", err
	}

	privateKey = string(pem.EncodeToMemory(encryptedPEMBlock))

	pubKey, err := getPublicKeyFromPrivateKey(privKey)
	if err != nil {
		return "", "", err
	}

	publicKey = string(ssh.MarshalAuthorizedKey(pubKey))

	return privateKey, publicKey, nil
}

// IsValidPrivateKey checks if a given string represents a valid private key.
func IsValidPrivateKey(key string) bool {
	block, _ := pem.Decode([]byte(key))
	return block != nil
}

// sshutils.go

// ... (other code)

// GetPublicKey extracts the public key from a private key string using the given password.
func GetPublicKey(privateKey string, password string) (string, error) {
	privKeyBlock, _ := pem.Decode([]byte(privateKey))
	if privKeyBlock == nil {
		return "", errors.New("failed to decode private key")
	}

	privKeyBytes, err := x509.DecryptPEMBlock(privKeyBlock, []byte(password))
	if err != nil {
		return "", err
	}

	privKeyBlock.Bytes = privKeyBytes

	pubKey, err := getPublicKeyFromPrivateKey(privKeyBlock)
	if err != nil {
		return "", err
	}

	return string(ssh.MarshalAuthorizedKey(pubKey)), nil
}

func generateECDSAKey(keyLength int) (*pem.Block, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	privBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	privBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	}

	return privBlock, privBytes, nil
}

func generateRSAKey(keyLength int) (*pem.Block, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return nil, nil, err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(key)

	privBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}

	return privBlock, privBytes, nil
}

func getPublicKeyFromPrivateKey(privKey *pem.Block) (ssh.PublicKey, error) {
	switch privKey.Type {
	case "EC PRIVATE KEY":
		key, err := x509.ParseECPrivateKey(privKey.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, err := ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		return publicKey, nil
	case "RSA PRIVATE KEY":
		key, err := x509.ParsePKCS1PrivateKey(privKey.Bytes)
		if err != nil {
			return nil, err
		}
		publicKey, err := ssh.NewPublicKey(&key.PublicKey)
		if err != nil {
			return nil, err
		}
		return publicKey, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", privKey.Type)
	}
}

func GetAvailableAlgorithms() []Algorithm {
	return []Algorithm{ECDSA, RSA}
}

func AlgorithmFromString(str string) Algorithm {

	if str == "ecdsa" {
		return ECDSA
	}

	if str == "rsa" {
		return RSA
	}

	return ECDSA
}

func AlgorithmToString(algo Algorithm) string {

	if algo == ECDSA {
		return "ecdsa"
	}

	if algo == RSA {
		return "rsa"
	}

	return "ecdsa"
}

func AlgorithmDefaultKeyLength(algo Algorithm) int32 {
	if algo == ECDSA {
		return 256
	}

	if algo == RSA {
		return 3072
	}

	return 256
}
