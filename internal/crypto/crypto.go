package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

func Encrypt(input string, key []byte) (string, error) {
	// Create a new SHA-256 hash of the key
	hash := sha256.Sum256(key)

	// Create a new AES cipher block using the hashed key
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	// Generate a random nonce
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Create a new GCM cipher instance
	aesGCM, err := cipher.NewGCMWithNonceSize(block, 12)
	if err != nil {
		return "", err
	}

	// Encrypt the input using AES-GCM
	cipherText := aesGCM.Seal(nil, nonce, []byte(input), nil)

	// Combine nonce and cipherText and base64 encode the result
	encryptedData := append(nonce, cipherText...)
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	return encodedData, nil
}

func Decrypt(encryptedBase64 string, key []byte) (string, error) {
	// Create a new SHA-256 hash of the key
	hash := sha256.Sum256(key)

	// Decode the base64-encoded data
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	// Extract nonce and cipherText from the encrypted data
	nonce := encryptedData[:12]
	cipherText := encryptedData[12:]

	// Create a new AES cipher block using the hashed key
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return "", err
	}

	// Create a new GCM cipher instance
	aesGCM, err := cipher.NewGCMWithNonceSize(block, 12)
	if err != nil {
		return "", err
	}

	// Decrypt the cipherText using AES-GCM and nonce
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// Encrypt data from the stdin pipe
func EncryptStdin(key []byte) (string, error) {
	// Read data from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	// Encrypt the data
	encryptedData, err := Encrypt(string(data), key)
	if err != nil {
		return "", err
	}

	// Print the encrypted data to stdout
	return encryptedData, nil
}

// Encrypt data read from a text file
func EncryptFromFile(inputFilePath string, key []byte) (string, error) {
	// Read data from the input file
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", nil
	}

	// Encrypt the data
	encryptedData, err := Encrypt(string(data), key)
	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

// Encrypt an input string to a file, with an option to append or overwrite
func EncryptToFile(input string, outputFilePath string, key []byte, append bool) error {
	encryptedData, err := Encrypt(input, key)
	if err != nil {
		return err
	}

	// Open the output file in append or overwrite mode
	var file *os.File
	if append {
		file, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		file, err = os.Create(outputFilePath)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the encrypted data to the file
	_, err = file.WriteString(encryptedData)
	if err != nil {
		return err
	}

	return nil
}

// Read from a source text file, encrypt, and write to a target file
func EncryptFile(sourceFilePath string, targetFilePath string, key []byte) error {
	encryptedData, err := EncryptFromFile(sourceFilePath, key)
	if err != nil {
		return err
	}

	// Write the encrypted data to the target file
	err = os.WriteFile(targetFilePath, []byte(encryptedData), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Decrypt data from the stdin pipe
func DecryptStdin(key []byte) (string, error) {
	// Read data from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	decryptedData, err := Decrypt(string(data), key)
	if err != nil {
		return "", err
	}

	// Print the decrypted data to stdout
	return decryptedData, nil
}

// Decrypt data read from a text file
func DecryptFromFile(inputFilePath string, key []byte) (string, error) {
	// Read data from the input file
	data, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", nil
	}

	// Decrypt the data
	decryptedData, err := Decrypt(string(data), key)
	if err != nil {
		return "", err
	}

	return decryptedData, nil
}

// Decrypt an input string to a file, with an option to append or overwrite
func DecryptToFile(input string, outputFilePath string, key []byte, append bool) error {
	decryptedData, err := Decrypt(input, key)
	if err != nil {
		return err
	}

	// Open the output file in append or overwrite mode
	var file *os.File
	if append {
		file, err = os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	} else {
		file, err = os.Create(outputFilePath)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the decrypted data to the file
	_, err = file.WriteString(decryptedData)
	if err != nil {
		return err
	}

	return nil
}

// Read from a source text file, decrypt, and write to a target file
func DecryptFile(sourceFilePath string, targetFilePath string, key []byte) error {
	decryptedData, err := DecryptFromFile(sourceFilePath, key)
	if err != nil {
		return err
	}

	// Write the decrypted data to the target file
	err = os.WriteFile(targetFilePath, []byte(decryptedData), 0644)
	if err != nil {
		return err
	}

	return nil
}
