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
	osExec "os/exec"
	"sort"
	"testing"

	"github.com/a13labs/sectool/cmd/exec"
	"github.com/a13labs/sectool/internal/crypto"
	"github.com/a13labs/sectool/internal/vault"
)

func TestHideSensitiveInfo(t *testing.T) {
	input := "This is a secret: $SECRET_KEY$"
	km := crypto.NewKeyManager()
	kv := crypto.NewSecureKVStore(km)
	kv.Put("SECRET", "$SECRET_KEY$")
	expected := "This is a secret: [HIDDEN]"
	output := exec.HideSensitiveInfo(input, kv)
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
	km := crypto.NewKeyManager()
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	vaultProvider.VaultSetValue("SECRET3", "This is a\nmultiline secret")
	expectedEnv := map[string]string{
		"ENV_VAR1":         "value1",
		"ENV_VAR2":         "value2",
		"SECRET_VAR1":      "$SECRET1",
		"SECRET_VAR2":      "$SECRET2",
		"COMPOSED_SECRET":  "$SECRET1.$SECRET2",
		"MULTILINE_SECRET": "$SECRET3",
	}
	env, _, err := exec.ParseEnvFile(envFile, vaultProvider, km)
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

func TestComposeEnv(t *testing.T) {
	envFile := "test.env"
	km := crypto.NewKeyManager()
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	vaultProvider.VaultSetValue("SECRET3", "This is a\nmultiline secret")
	expectedEnv := []string{"ENV_VAR1=value1", "ENV_VAR2=value2", "SECRET_VAR1=secret_value1", "SECRET_VAR2=secret_value2", "COMPOSED_SECRET=secret_value1.secret_value2", "MULTILINE_SECRET=This is a\nmultiline secret"}
	env, kv, err := exec.ParseEnvFile(envFile, vaultProvider, km)
	if err != nil {
		t.Errorf("Error parsing env file: %v", err)
	}
	composedEnv, err := exec.ComposeEnv(env, kv)
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

func TestExec(t *testing.T) {
	cmdExec := osExec.Command("/bin/sh", "-c", "echo $SECRET_VAR1 $SECRET_VAR2 \"$MULTILINE_SECRET\"")
	envFile := "test.env"
	km := crypto.NewKeyManager()
	vaultProvider := vault.NewDummyVault()
	vaultProvider.VaultSetValue("SECRET1", "secret_value1")
	vaultProvider.VaultSetValue("SECRET2", "secret_value2")
	vaultProvider.VaultSetValue("SECRET3", "This is a\nmultiline secret")
	env, kv, err := exec.ParseEnvFile(envFile, vaultProvider, km)
	if err != nil {
		t.Errorf("Error parsing env file: %v", err)
	}
	composedEnv, err := exec.ComposeEnv(env, kv)
	if err != nil {
		t.Errorf("Error composing env: %v", err)
	}
	cmdExec.Env = append(cmdExec.Env, composedEnv...)
	output, err := cmdExec.Output()
	if err != nil {
		t.Errorf("Error running command: %v", err)
	}
	expectedOutput := "secret_value1 secret_value2 This is a\nmultiline secret\n"
	if string(output) != expectedOutput {
		t.Errorf("Expected %q, but got %q", expectedOutput, string(output))
	}
}
