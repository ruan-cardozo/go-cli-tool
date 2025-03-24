package count_comments_test

import (
	"bytes"
	count_comments "go-cli-tool/cmd/count-comments"
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

func TestCountCommentsCommandWithNoArgs(t *testing.T) {

	tests.ResetGlobals()

    cmd := count_comments.CountCommentsCmd

    var stdout bytes.Buffer
    cmd.SetOut(&stdout)

    cmd.SetArgs([]string{})

	cmd.Execute()

    expectedOutput := "Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountCommentsCommandWithFilePath(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_comments.CountCommentsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	cmd.SetArgs([]string{"-f", "../../javascript-tests/test.js"})

	cmd.Execute()

	expectedOutput := "\x1b[34mTotal comments:\x1b[0m 3\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountCommentsCommandWithWrongFilePathExtension(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_comments.CountCommentsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-f", "../../main.go"})

	cmd.Execute()

	expectedOutput := "\x1b[31mOnly JavaScript files are accepted.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountCommentsCommandWithDirectoryPath(t *testing.T) {

	tests.ResetGlobals()

    cmd := count_comments.CountCommentsCmd

    var stdout bytes.Buffer
    cmd.SetOut(&stdout) 
    cmd.SetErr(&stdout)

    cmd.SetArgs([]string{"-d", "../../javascript-tests"})

    if err := cmd.Execute(); err != nil {
        t.Errorf("CountCommentsCommand() error = %v, want nil", err)
    }

    expectedOutput := "\x1b[34m Comment lines in test.js:\x1b[0m 3\n\x1b[34mTotal Comments in directory:\x1b[0m 3\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountCommentsCommandWithWrongDirectoryPath(t *testing.T) {
	
	tests.ResetGlobals()

    cmd := count_comments.CountCommentsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../../main.go"})

	cmd.Execute()

	expectedOutput := "\x1b[31mPlease provide a valid directory path.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountCommentsCommandWithDirectoryWihtoutJavascriptFiles(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_comments.CountCommentsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../root"})

	cmd.Execute()

	expectedOutput := "\x1b[31mNo JavaScript files found in the provided directory.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountCommentsCommandWithValidDirectoryAndGeneratingReportHTML(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_comments.CountCommentsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../../javascript-tests", "-o", "."})

	cmd.Execute()

	expectedOutput := "\x1b[1;34mReport generated successfully at report.html\x1b[0m\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountCommentsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}

	reportPath := "./report.html"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Errorf("Expected report file to be created at %s, but it does not exist", reportPath)
	}

	if err := os.Remove(reportPath); err != nil {
		t.Errorf("Failed to remove the report file: %v", err)
	}
}