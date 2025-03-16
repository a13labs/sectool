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
	"github.com/a13labs/sectool/internal/config"
	"github.com/a13labs/sectool/internal/vault"
	"github.com/spf13/cobra"
)

var config_file = ""
var no_output = false
var cfg *config.Config

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute a command",
	Long:  `Execute a command with the environment variables from the .env file and the vault file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: sectool exec <cmd> <args>")
			os.Exit(0)
		}

		cmdToRun, cmdArgs := ProcessArgs(args)

		var err error
		cfg, err = config.ReadConfig(config_file)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			os.Exit(1)
		}

		vaultProvider, err := vault.NewVaultProvider(*cfg)
		if err != nil {
			fmt.Println("Error initializing vault provider.")
			os.Exit(1)
		}

		cmdExec := exec.Command(cmdToRun, cmdArgs...)
		cmdExec.Env = append(os.Environ(), "SECTOOL_ENV=1")
		envMap, vaultValues, sensitiveStrings, err := ParseEnvFile("sectool.env", vaultProvider)
		if err != nil {
			fmt.Printf("Error parsing env file: %v\n", err)
			os.Exit(1)
		}
		if no_output {
			fmt.Println("Command started.")
		} else {
			stdoutPipe, _ := cmdExec.StdoutPipe()
			stderrPipe, _ := cmdExec.StderrPipe()

			go func() {
				scanner := bufio.NewScanner(stdoutPipe)
				for scanner.Scan() {
					line := scanner.Text()
					fmt.Println(HideSensitiveInfo(line, sensitiveStrings))
				}
			}()

			go func() {
				scanner := bufio.NewScanner(stderrPipe)
				for scanner.Scan() {
					line := scanner.Text()
					fmt.Println(HideSensitiveInfo(line, sensitiveStrings))
				}
			}()
		}

		// Append the environment variables from the env file
		envVars, err := ComposeEnv(envMap, vaultValues)
		if err != nil {
			fmt.Printf("Error composing environment variables: %v\n", err)
			os.Exit(1)
		}
		cmdExec.Env = append(cmdExec.Env, envVars...)

		if err := cmdExec.Start(); err != nil {
			fmt.Printf("Error starting command: %v\n", err)
			os.Exit(1)
		}

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
	// Disable flag parsing since we want to handle arguments after the first argument
	execCmd.DisableFlagParsing = true
}

// HideSensitiveInfo replaces sensitive strings with a placeholder
func HideSensitiveInfo(input string, sensitiveStrings []string) string {
	result := input
	for _, s := range sensitiveStrings {
		result = strings.Replace(result, s, "[HIDDEN]", -1)
	}
	return result
}

// Define a regular expression pattern to match environment variables
const pattern = `\s*\$([a-zA-Z_][a-zA-Z0-9_]*)`

// override the default behavior of the flag package to allow for arguments after the first argument
func ProcessArgs(args []string) (string, []string) {

	var arguments []string
	cmd := ""
	foundFirstArg := false

	i := 0 // Initialize a loop variable
	for i < len(args) {
		arg := args[i]

		if !foundFirstArg {
			if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
				if arg == "--config" || arg == "-f" {
					i++
					if i == len(args) {
						fmt.Println("Missing vault file location.")
						os.Exit(1)
					}
					config_file = args[i]
					if strings.HasPrefix(config_file, "--") || strings.HasPrefix(config_file, "-") {
						fmt.Println("Invalid value for vault file location.")
						os.Exit(1)
					}
				}
				if arg == "--no-output" || arg == "-n" {
					no_output = true
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

func ParseEnvFile(envFile string, v vault.VaultProvider) (map[string]string, map[string]string, []string, error) {
	// Read the contents of the file
	contents, err := os.ReadFile(envFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil, nil, nil, err
	}

	// Split the contents of the file into lines
	lines := strings.Split(string(contents), "\n")

	// Create a new map to store the environment variables
	env := make(map[string]string)

	lineNr := 0
	// Iterate over the lines in the file
	for _, line := range lines {
		lineNr++
		// skip empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, nil, nil, fmt.Errorf("invalid line %d: %s", lineNr, line)
		}

		// Extract the environment variable name
		envName := parts[0]
		envValue := parts[1]

		env[envName] = envValue
	}

	// Extract the sensitive strings from the environment variables
	usedKeys := []string{}
	sensitiveStrings := []string{}
	sensitiveValues := []string{}

	regex := regexp.MustCompile(pattern)
	for _, value := range env {

		matches := regex.FindAllStringSubmatch(value, -1)
		if len(matches) == 0 {
			sensitiveValues = append(sensitiveValues, value)
			continue
		}

		// Iterate over the matches
		for _, match := range matches {
			keyName := match[1]
			usedKeys = append(usedKeys, keyName)
		}
	}

	vaultValues, err := v.VaultGetMultipleValues(usedKeys)
	if err != nil {
		fmt.Printf("Error getting multiple values from vault: %v\n", err)
		return nil, nil, nil, err
	}

	for _, value := range vaultValues {
		sensitiveStrings = append(sensitiveStrings, value)
	}

	sensitiveStrings = append(sensitiveStrings, sensitiveValues...)

	// Return the environment map
	return env, vaultValues, sensitiveStrings, nil
}

// ComposeEnv reads the environment variables from the file and replaces the keys with the values from the vault
func ComposeEnv(e map[string]string, vaultValues map[string]string) ([]string, error) {

	// Create a new slice to store the environment variables
	env := []string{}

	// Extract the sensitive strings from the environment variables
	regex := regexp.MustCompile(pattern)
	for key, value := range e {

		composedValue := value
		matches := regex.FindAllStringSubmatch(value, -1)
		if len(matches) > 0 {
			// Iterate over the matches
			for _, match := range matches {
				secretKey := match[1]
				composedValue = strings.Replace(composedValue, "$"+secretKey, vaultValues[secretKey], -1)
			}
		}

		env = append(env, fmt.Sprintf("%s=%s", key, composedValue))
	}

	// Return the environment slice
	return env, nil
}
