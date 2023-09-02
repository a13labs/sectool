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
package exec

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/a13labs/sectool/cmd"
	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

var vault_file = ""

func hideSensitiveInfo(input string, sensitiveStrings []string) string {
	result := input
	for _, s := range sensitiveStrings {
		result = strings.Replace(result, s, "<sensitive data>", -1)
	}
	return result
}

const pattern = `\s*\$([a-zA-Z_][a-zA-Z0-9_]*)`

func processArgs(args []string) (string, []string) {

	var arguments []string
	cmd := ""
	foundFirstArg := false

	i := 0 // Initialize a loop variable
	for i < len(args) {
		arg := args[i]

		if !foundFirstArg {
			if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
				if arg == "--vault" || arg == "-v" {
					i++
					if i == len(args) {
						fmt.Println("Missing vault file location.")
						os.Exit(1)
					}
					vault_file = args[i]
					if strings.HasPrefix(vault_file, "--") || strings.HasPrefix(vault_file, "-") {
						fmt.Println("Invalid value for vault file location.")
						os.Exit(1)
					}
				}
			} else {
				// The first non-option argument is encountered
				foundFirstArg = true
				cmd = arg
			}
		} else {
			// Handle arguments after the first argument
			arguments = append(arguments, arg)
		}

		i++ // Increment the loop variable
	}

	return cmd, arguments
}

func readVaultLocation() (string, error) {
	// Open the file for reading
	file, err := os.Open(".vault")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a bufio scanner to read lines from the file
	scanner := bufio.NewScanner(file)

	// Check if there is at least one line in the file
	if scanner.Scan() {
		// Read the first line
		firstLine := scanner.Text()
		return firstLine, nil
	} else if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("The file is empty")
}

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a recipe",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: sectool exec <cmd> <args>")
			os.Exit(0)
		}
		masterPwd, exist := os.LookupEnv("VAULT_MASTER_PASSWORD")
		if !exist {
			fmt.Println("VAULT_MASTER_PASSWORD it's not defined, aborting.")
			os.Exit(1)
		}

		cmdToRun, cmdArgs := processArgs(args)

		if vault_file == "" {
			vault_file, _ = readVaultLocation()
		}

		if vault_file == "" {
			vault_file = "repository.vault"
		}

		v := vault.NewVault(vault_file, []byte(masterPwd))

		cmdExec := exec.Command(cmdToRun, cmdArgs...)
		stdoutPipe, _ := cmdExec.StdoutPipe()
		stderrPipe, _ := cmdExec.StderrPipe()
		cmdExec.Env = append(os.Environ(), "SECTOOL_ENV=1")
		sensitiveStrings := []string{masterPwd}

		env_file := "sectool.env"
		_, err := os.Stat(env_file)
		if err == nil {
			contents, err := os.ReadFile(env_file)
			if err != nil {
				fmt.Println("File 'sectool.env' is invalid, ignoring.")
			} else {
				// Compile the regular expression pattern
				regex := regexp.MustCompile(pattern)
				lines := strings.Split(string(contents), "\n")
				lineNr := 0
				for _, line := range lines {
					lineNr += 1
					parts := strings.SplitN(line, "=", 2)
					if len(parts) != 2 {
						continue
					}
					envName := parts[0]
					// Find all matches in the input string
					matches := regex.FindAllStringSubmatch(parts[1], -1)
					result := parts[1]
					for _, match := range matches {
						keyName := match[1]
						keyValue, err := v.VaultGetValue(keyName)
						if err != nil {
							fmt.Printf("'sectool.env': Error getting value '%s', line %d.", keyName, lineNr)
							os.Exit(1)
						}
						result = strings.Replace(result, fmt.Sprintf("$%s", keyName), keyValue, -1)

						escaped := strings.Replace(keyValue, "\\", "\\\\", -1)
						escaped = strings.Replace(escaped, "$", "\\$", -1)

						sensitiveStrings = append(sensitiveStrings, escaped)
					}
					cmdExec.Env = append(cmdExec.Env, fmt.Sprintf("%s=%s", envName, result))
				}
			}
		}

		if err := cmdExec.Start(); err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			os.Exit(1)
		}

		// Create goroutines to read and hide output from stdout and stderr
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(hideSensitiveInfo(line, sensitiveStrings))
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Println(hideSensitiveInfo(line, sensitiveStrings))
			}
		}()

		err = cmdExec.Wait()
		if err != nil {
			exitError := err.(*exec.ExitError)
			exitCode := exitError.ExitCode()
			os.Exit(exitCode)
		}
		os.Exit(0)
	},
}

func init() {
	cmd.RootCmd.AddCommand(execCmd)
	execCmd.DisableFlagParsing = true
}
