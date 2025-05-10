package count_percent_test

import (
	"bytes"
	count_percent "go-cli-tool/cmd/count-percent-lines"
	"go-cli-tool/tests"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code before running tests
	tests.ResetGlobals()

	// Run tests
	code := m.Run()

	// Cleanup code after running tests
	tests.ResetGlobals()

	// Exit with the test code
	os.Exit(code)
}

func TestCountPercentCommandWithFilePath(t *testing.T) {
	tests.ResetGlobals()

	cmd := count_percent.CountPercentCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	cmd.SetArgs([]string{"-f", "../../javascript-tests/test.js"})

	cmd.Execute()

	expectedOutput := "\x1b[34mTotal comments:\x1b[0m 3\n\x1b[34mComment Percentage:\x1b[0m 30.00%\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountPercentCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}
