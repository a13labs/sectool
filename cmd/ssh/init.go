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

	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

// initCmd represents the list command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init SSH key pairs management",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing password.")
			os.Exit(1)
		}

		master_pwd, exist := os.LookupEnv("VAULT_MASTER_PASSWORD")
		if !exist {
			fmt.Println("VAULT_MASTER_PASSWORD it's not defined, aborting.")
			os.Exit(1)
		}

		v := vault.NewVault("repository.vault", []byte(master_pwd))

		if v.VaultHasKey("SSH_MASTER_PASSWORD") {
			fmt.Println("SSH_MASTER_PASSWORD already defined, aborting.")
			os.Exit(1)
		}

		v.VaultSetValue("SSH_MASTER_PASSWORD", args[0])
		fmt.Println("SSH key management successfully initialized.")
	},
}

func init() {
	sshCmd.AddCommand(initCmd)
}
