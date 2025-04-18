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

	"github.com/a13labs/sectool/cmd"
	"github.com/a13labs/sectool/vault"
	"github.com/spf13/cobra"
)

var quoted bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a key/value from the vault.",
	Long:  ``,
	Run: func(c *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Usage: sectool vault get <key>")
			os.Exit(1)
		}

		value, err := vault.GetSecret(cmd.ConfigFile, args[0])
		if err != nil {
			fmt.Println("Error getting value.")
			os.Exit(1)
		}

		if quoted {
			value = "\"" + value + "\""
		}

		fmt.Println(value)
		os.Exit(0)
	},
}

func init() {
	vaultCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVarP(&quoted, "quoted", "q", false, "Quote value")
}
