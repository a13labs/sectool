/*
Copyright © 2023 Alexandre Pires

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
package ssh

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/a13labs/sectool/internal/crypto"
	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

// unlockCmd represents the list command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "A Unlock SSH key pairs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		master_pwd, exist := os.LookupEnv("VAULT_MASTER_PASSWORD")
		if !exist {
			fmt.Println("VAULT_MASTER_PASSWORD it's not defined. Aborting")
			os.Exit(1)
		}

		_, err := os.Stat("ssh-keys")
		if os.IsNotExist(err) {
			fmt.Println("Missing keys directory, no keys can be listed.")
			os.Exit(1)
		}

		keys, err := listKeys("ssh-keys")
		if err != nil {
			fmt.Println("Error listing keys.")
			os.Exit(1)
		}

		v := vault.NewVault("repository.vault", []byte(master_pwd))
		ssh_master_password, err := v.VaultGetValue("SSH_MASTER_PASSWORD")
		if err != nil {
			fmt.Println("Error reading SSH_MASTER_PASSWORD, aborting.")
			os.Exit(1)
		}

		for _, key := range keys {
			key_root_path := filepath.Join("ssh-keys", key)
			key_path := filepath.Join(key_root_path, "id_ecdsa")
			_, err := os.Stat(key_path + ".key")
			if !os.IsNotExist(err) {
				key_path := filepath.Join(key_root_path, "id_rsa")
				_, err := os.Stat(key_path + ".key")
				if !os.IsNotExist(err) {
					fmt.Printf("No key data in '%s', skipping.\n", key)
					continue
				}
			}

			err = crypto.DecryptFile(key_path+".key", key_path, []byte(ssh_master_password))
			if err != nil {
				fmt.Printf("Error decrypting key data in '%s', skipping.\n", key)
				continue
			}
			fmt.Printf("Key data in '%s' unlocked.\n", key)
		}
		os.Exit(0)
	},
}

func init() {
	sshCmd.AddCommand(unlockCmd)
}
