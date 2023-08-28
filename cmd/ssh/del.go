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

	"github.com/spf13/cobra"
)

// delCmd represents the list command
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "Delete a SSH key pair",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Missing key name")
			os.Exit(1)
		}

		_, err := os.Stat("ssh-keys")
		if os.IsNotExist(err) {
			fmt.Println("Missing keys directory, no key can be deleted.")
			os.Exit(1)
		}

		key_path := filepath.Join("ssh-keys", args[0])
		_, err = os.Stat(key_path)
		if os.IsNotExist(err) {
			fmt.Println("No key to be deleted.")
			os.Exit(1)
		}

		err = os.RemoveAll(key_path)
		if err != nil {
			fmt.Println("Error deleting key pair:", err)
			os.Exit(1)
		}

		fmt.Println("Key pair successfully removed.")
		os.Exit(0)
	},
}

func init() {
	sshCmd.AddCommand(delCmd)
}
