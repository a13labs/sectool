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

	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/crypto"
	"github.com/a13labs/sectool/internal/ssh"
	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

// unlockCmd represents the list command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "A Unlock SSH key pairs",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		_, err := os.Stat("ssh-keys")
		if os.IsNotExist(err) {
			fmt.Println("Missing keys directory, no keys can be listed.")
			os.Exit(1)
		}

		cfg, err := config.ReadConfig(config_file)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}

		vaultProvider, err := vault.NewVaultProvider(*cfg)
		if err != nil {
			fmt.Println("Error initializing vault provider.")
			os.Exit(1)
		}

		key := cfg.SSHPasswordKey
		if key == "" {
			key = defaultSSHPasswordKey
		}

		keys, err := listKeys("ssh-keys")
		if err != nil {
			fmt.Println("Error listing keys.")
			os.Exit(1)
		}

		ssh_master_password, err := vaultProvider.VaultGetValue(key)
		if err != nil {
			fmt.Printf("Error reading '%s', aborting.", key)
			os.Exit(1)
		}

		for _, key := range keys {
			key_root_path := filepath.Join("ssh-keys", key)

			for _, algo := range ssh.GetAvailableAlgorithms() {
				prefix := ssh.AlgorithmToString(algo)
				key_path := filepath.Join(key_root_path, "id_"+prefix)
				_, err := os.Stat(key_path + ".key")

				if err != nil {
					fmt.Printf("No encrypted private key (%s) in '%s', skipping.\n", prefix, key)
					continue
				}

				err = crypto.DecryptFile(key_path+".key", key_path, []byte(ssh_master_password))
				if err != nil {
					fmt.Printf("Error decrypting private (%s) data in '%s', skipping.\n", prefix, key)
					continue
				}
				fmt.Printf("Private key (%s) in '%s' unlocked.\n", prefix, key)
			}
		}
		os.Exit(0)
	},
}

func init() {
	sshCmd.AddCommand(unlockCmd)
}
