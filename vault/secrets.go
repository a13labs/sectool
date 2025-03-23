/*
Copyright © 2025 Alexandre Pires

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package vault

import (
	"fmt"
	"io"
	"os"

	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/vault"
)

func GetSecret(path string, key string) (string, error) {
	cfg, err := config.ReadConfig(path)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return "", err
	}

	vaultProvider, err := vault.NewVaultProvider(*cfg)
	if err != nil {
		fmt.Println("Error initializing vault provider.")
		return "", err
	}
	raw_value, err := vaultProvider.VaultGetValue(key)
	if err != nil {
		fmt.Println("Error getting value.")
		return "", err
	}

	return raw_value, nil
}

func DeleteSecret(path string, key string) error {
	cfg, err := config.ReadConfig(path)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return err
	}

	vaultProvider, err := vault.NewVaultProvider(*cfg)
	if err != nil {
		fmt.Println("Error initializing vault provider.")
		return err
	}

	err = vaultProvider.VaultDelKey(key)
	if err != nil {
		fmt.Println("Error deleting key/value.")
		return err
	}

	return nil
}

func ListSecrets(path string) ([]string, error) {
	cfg, err := config.ReadConfig(path)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return nil, err
	}

	vaultProvider, err := vault.NewVaultProvider(*cfg)
	if err != nil {
		fmt.Println("Error initializing vault provider.")
		return nil, err
	}
	return vaultProvider.VaultListKeys(), nil
}

func SetSecret(path string, key string, value string, backup bool) error {
	cfg, err := config.ReadConfig(path)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return err
	}

	vaultProvider, err := vault.NewVaultProvider(*cfg)
	if err != nil {
		fmt.Println("Error initializing vault provider.")
		return err
	}

	vaultProvider.VaultEnableBackup(backup)
	// if key starts with "file://", read from file
	if len(key) > 7 && key[:7] == "file://" {
		if _, err := os.Stat(key[7:]); os.IsNotExist(err) {
			fmt.Println("File does not exist.")
			return err
		}
		v, err := os.ReadFile(key[7:])
		if err != nil {
			fmt.Println("Error reading file.")
			return err
		}
		key = string(v)
	}
	// if key is "stdin://", read from stdin
	if key == "stdin://" {
		v, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading stdin:", err)
			return err
		}
		key = string(v)
	}
	err = vaultProvider.VaultSetValue(key, value)
	if err != nil {
		fmt.Println("Error setting key/value.")
		return err
	}

	return nil
}
