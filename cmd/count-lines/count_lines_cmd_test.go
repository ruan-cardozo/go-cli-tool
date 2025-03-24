package count_lines_test

import (
	"bytes"
	count_lines "go-cli-tool/cmd/count-lines"
	"testing"
)

func TestCountLinesCommandWithNoArgs(t *testing.T) {
    // create the count lines command
    cmd := count_lines.CountLinesAnalyzer

    // redirect the stdout to a buffer to capture the output
    var stdout bytes.Buffer
    cmd.SetOut(&stdout)

    // set the args
    cmd.SetArgs([]string{})

    // execute the count lines command w/ args
	cmd.Execute()

    // check the output
    expectedOutput := "Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountLinesCommandWithFilePath(t *testing.T) {
	// create the count lines command
	cmd := count_lines.CountLinesAnalyzer

	// redirect the stdout to a buffer to capture the output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 
    cmd.SetErr(&stdout)

	// set the args
	cmd.SetArgs([]string{"-f", "../../javascript-tests/test.js"})

	// execute the count lines command w/ args
	cmd.Execute()

	// check the output
	expectedOutput := "\x1b[34mTotal lines:\x1b[0m 157\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}