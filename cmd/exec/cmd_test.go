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
package exec_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/a13labs/sectool/cmd/exec"
	"github.com/a13labs/sectool/internal/vault"
)

func TestHideSensitiveInfo(t *testing.T) {
	input := "This is a secret: $SECRET_KEY$"
	sensitiveStrings := []string{"$SECRET_KEY$"}
	expected := "This is a secret: [HIDDEN]"
	output := exec.HideSensitiveInfo(input, sensitiveStrings)
	if output != expected {
		t.Errorf("Expected %q, but got %q", expected, output)
	}
}

func TestProcessArgs(t *testing.T) {
	args := []string{"echo", "Hello", "World"}
	expectedCmd := "echo"
	expectedArgs := []string{"Hello", "World"}
	cmdToRun, cmdArgs := exec.ProcessArgs(args)
	if cmdToRun != expectedCmd {
		t.Errorf("Expected cmdToRun %q, but got %q", expectedCmd, cmdToRun)
	}
	if len(cmdArgs) != len(expectedArgs) {
		t.Errorf("Expected %d cmdArgs, but got %d", len(expectedArgs), len(cmdArgs))
	}
	for i := range cmdArgs {
		if cmdArgs[i] != expectedArgs[i] {
			t.Errorf("Expected cmdArgs[%d] %q, but got %q", i, expectedArgs[i], cmdArgs[i])
		}
	}
}

func TestParseEnvFile(t *testing.T) {
	envFile := "test.env"
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	expectedEnv := map[string]string{
		"ENV_VAR1":        "value1",
		"ENV_VAR2":        "value2",
		"SECRET_VAR1":     "$SECRET1",
		"SECRET_VAR2":     "$SECRET2",
		"COMPOSED_SECRET": "$SECRET1.$SECRET2",
	}
	env, _, _, err := exec.ParseEnvFile(envFile, vaultProvider)
	if err != nil {
		t.Errorf("Error parsing env file: %v", err)
	}
	if len(env) != len(expectedEnv) {
		t.Errorf("Expected %d env variables, but got %d", len(expectedEnv), len(env))
	}
	for k, v := range env {
		if expectedEnv[k] != v {
			t.Errorf("Expected env[%q] %q, but got %q", k, expectedEnv[k], v)
		}
	}
}

func TestGetSensitiveStrings(t *testing.T) {

	envFile := "test.env"
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	expectedSensitiveStrings := []string{"value1", "value2", "secret_value1", "secret_value2"}
	_, _, sensitiveStrings, err := exec.ParseEnvFile(envFile, vaultProvider)
	if err != nil {
		t.Errorf("Error parsing env file: %v", err)
	}
	sort.Strings(sensitiveStrings)
	sort.Strings(expectedSensitiveStrings)
	if err != nil {
		t.Errorf("Error getting sensitive strings: %v", err)
	}
	if len(sensitiveStrings) != len(expectedSensitiveStrings) {
		t.Errorf("Expected %d sensitive strings, but got %d", len(expectedSensitiveStrings), len(sensitiveStrings))
	}
	for i := range sensitiveStrings {
		if sensitiveStrings[i] != expectedSensitiveStrings[i] {
			t.Errorf("Expected sensitiveStrings[%d] %q, but got %q", i, expectedSensitiveStrings[i], sensitiveStrings[i])
		}
	}
}

func TestComposeEnv(t *testing.T) {
	envFile := "test.env"
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	expectedEnv := []string{"ENV_VAR1=value1", "ENV_VAR2=value2", "SECRET_VAR1=secret_value1", "SECRET_VAR2=secret_value2", "COMPOSED_SECRET=secret_value1.secret_value2"}
	env, vaultValues, _, err := exec.ParseEnvFile(envFile, vaultProvider)
	if err != nil {
		t.Errorf("Error parsing env file: %v", err)
	}
	composedEnv, err := exec.ComposeEnv(env, vaultValues)
	sort.Strings(expectedEnv)
	sort.Strings(composedEnv)
	if err != nil {
		t.Errorf("Error composing env: %v", err)
	}
	if len(composedEnv) != len(expectedEnv) {
		fmt.Println(composedEnv)
		t.Errorf("Expected %d env variables, but got %d", len(expectedEnv), len(composedEnv))
	}
	for i := range composedEnv {
		if composedEnv[i] != expectedEnv[i] {
			t.Errorf("Expected composedEnv[%d] %q, but got %q", i, expectedEnv[i], composedEnv[i])
		}
	}
}
