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
	"testing"

	"github.com/a13labs/sectool/cmd/exec"
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

func TestComposeEnv(t *testing.T) {
	envFile := "test.env"
	expectedEnv := []string{"ENV_VAR1=value1", "ENV_VAR2=value2"}
	env := exec.ComposeEnv(envFile)
	if len(env) != len(expectedEnv) {
		t.Errorf("Expected %d env variables, but got %d", len(expectedEnv), len(env))
	}
	for i := range env {
		if env[i] != expectedEnv[i] {
			t.Errorf("Expected env[%d] %q, but got %q", i, expectedEnv[i], env[i])
		}
	}
}
