package crypto

import (
	"io"
	"os"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	key := []byte("mysecretkey")

	// Test Encrypt and Decrypt round-trip
	input := "Hello, this is a test message!"
	encrypted, err := Encrypt(input, key)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != input {
		t.Errorf("Expected decrypted data to be '%s', but got '%s'", input, decrypted)
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	key := []byte("mysecretkey")
	sourceFilePath := "testdata/source.txt"
	encryptedFilePath := "testdata/encrypted.txt"
	decryptedFilePath := "testdata/decrypted.txt"

	// Clean up any previous test artifacts
	defer func() {
		_ = os.Remove(encryptedFilePath)
		_ = os.Remove(decryptedFilePath)
	}()

	// Read test data from source file, Encrypt, and write to encrypted file
	err := EncryptFile(sourceFilePath, encryptedFilePath, key)
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt the encrypted file and write to decrypted file
	err = DecryptFile(encryptedFilePath, decryptedFilePath, key)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the decrypted data with the original source file content
	originalData, err := os.ReadFile(sourceFilePath)
	if err != nil {
		t.Fatal(err)
	}

	decryptedData, err := os.ReadFile(decryptedFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(decryptedData) != string(originalData) {
		t.Errorf("Decrypted data does not match original data")
	}
}

// Add more tests for other functions...
func TestEncryptAndDecryptStdin(t *testing.T) {
	key := []byte("mysecretkey")

	// Test Encrypt and Decrypt using standard input
	input := "Hello, this is a test message!"
	encryptedOutput, err := testStdin(EncryptStdin, key, input)
	if err != nil {
		t.Fatal(err)
	}

	decryptedOutput, err := testStdin(DecryptStdin, key, encryptedOutput)
	if err != nil {
		t.Fatal(err)
	}

	if decryptedOutput != input {
		t.Errorf("Expected decrypted data to be '%s', but got '%s'", input, decryptedOutput)
	}
}

func testStdin(fn func([]byte) (string, error), key []byte, input string) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	defer r.Close()
	defer w.Close()

	// Redirect os.Stdin to the write end of the pipe
	tmp := os.Stdin
	os.Stdin = r

	// Write input to the pipe
	go func() {
		defer w.Close()
		_, _ = w.Write([]byte(input))
	}()

	// Capture the output of the function
	output, err := fn(key)
	// Restore os.Stdin
	os.Stdin = tmp

	if err != nil {
		return "", err
	}

	return output, nil
}

func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old

	return string(out)
}

func TestEncryptFromFile(t *testing.T) {
	key := []byte("mysecretkey")
	inputFilePath := "testdata/source.txt"

	// Test encrypting data from a file
	encrypted, err := EncryptFromFile(inputFilePath, key)
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt the encrypted data and compare with the original content
	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}

	originalData, err := os.ReadFile(inputFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != string(originalData) {
		t.Errorf("Decrypted data does not match original data")
	}
}

func TestEncryptToFile(t *testing.T) {
	key := []byte("mysecretkey")
	input := "Hello, this is a test message!"
	outputFilePath := "testdata/output.txt"

	// Clean up any previous test artifacts
	defer func() {
		_ = os.Remove(outputFilePath)
	}()

	// Test encrypting data to a file
	err := EncryptToFile(input, outputFilePath, key, false)
	if err != nil {
		t.Fatal(err)
	}

	// Decrypt the encrypted data from the output file and compare with the original content
	decrypted, err := DecryptFromFile(outputFilePath, key)
	if err != nil {
		t.Fatal(err)
	}

	if decrypted != input {
		t.Errorf("Decrypted data does not match original data")
	}

	encrypted, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Test decrypting data to a file
	err = DecryptToFile(string(encrypted), outputFilePath, key, false)
	if err != nil {
		t.Fatal(err)
	}

	originalData, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatal(err)
	}

	if input != string(originalData) {
		t.Errorf("Decrypted data does not match original data")
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
