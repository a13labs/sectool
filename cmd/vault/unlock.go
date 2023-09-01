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
package vault

import (
	"fmt"
	"os"

	"github.com/a13labs/sectool/internal/crypto"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock vault",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		master_pwd, exist := os.LookupEnv("VAULT_MASTER_PASSWORD")
		if !exist {
			fmt.Println("VAULT_MASTER_PASSWORD it's not defined, aborting.")
			os.Exit(1)
		}

		unlocked_vault := "repository.vault.unlocked"
		locked_vault := "repository.vault"

		_, err := os.Stat(locked_vault)
		if os.IsNotExist(err) {
			fmt.Println("Vault does not exists, nothing to be unlocked.")
			os.Exit(1)
		}

		err = crypto.DecryptFile(locked_vault, unlocked_vault, []byte(master_pwd))
		if err != nil {
			fmt.Println("Error unlocking vault.")
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	vaultCmd.AddCommand(unlockCmd)
}
