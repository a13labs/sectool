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

	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

// addCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a key/value from the vault.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing key.")
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

		err = vaultProvider.VaultDelKey(args[0])
		if err != nil {
			fmt.Println("Error deleting key/value.")
			os.Exit(1)
		}
		fmt.Println("Key/value deleted")
		os.Exit(0)
	},
}

func init() {
	vaultCmd.AddCommand(delCmd)
}
