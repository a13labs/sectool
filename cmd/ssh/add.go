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

	"github.com/a13labs/sectool/internal/ssh"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var key_length int32
var key_algo_str string

// addCmd represents the list command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH key pair",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Missing key name.")
			os.Exit(1)
		}

		_, err := os.Stat("ssh-keys")
		if os.IsNotExist(err) {
			os.Mkdir("ssh-keys", 0755)
		} else if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		key_path := filepath.Join("ssh-keys", args[0])
		_, err = os.Stat(key_path)
		if os.IsNotExist(err) {
			os.Mkdir(key_path, 0755)
		} else if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		algo := ssh.ECDSA
		if key_algo_str == "rsa" {
			algo = ssh.RSA
		}

		key_full_path_priv := filepath.Join(key_path, "id_"+key_algo_str)
		key_full_path_pub := filepath.Join(key_path, "id_"+key_algo_str+".pub")
		_, err = os.Stat(key_full_path_priv)
		if err == nil {
			fmt.Println("Key already exists, delete first")
			os.Exit(1)
		} else if !os.IsNotExist(err) {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		fmt.Print("Enter password: ")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Generating key pair.")
		priv, pub, err := ssh.GenerateKeyPair(algo, int(key_length), string(password))
		if err != nil {
			fmt.Println("Error creating ssh key")
			return
		}

		os.WriteFile(key_full_path_priv, []byte(priv), 0600)
		os.WriteFile(key_full_path_pub, []byte(pub), 0644)

		fmt.Println("Key pair successfully generated.")
		os.Exit(0)
	},
}

func init() {
	sshCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&key_algo_str, "algo", "a", "ecdsa", "Algorithm: ecdsa, rsa, default: ecdsa")
	addCmd.Flags().Int32VarP(&key_length, "length", "l", 0, "Key length, default: ecdsa(256), rsa(3072)")
}
