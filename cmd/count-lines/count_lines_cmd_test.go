package count_lines_test

import (
	"bytes"
	count_lines "go-cli-tool/cmd/count-lines"
	"go-cli-tool/internal/utils"
	"os"
	"testing"
)

func resetGlobals() {
	utils.FilePath = ""
	utils.DirectoryPath = ""
	utils.OutputFilePath = ""
}

func TestMain(m *testing.M) {
    // Setup code before running tests
	resetGlobals()

    // Run tests
    code := m.Run()

    // Cleanup code after running tests
	resetGlobals()

    // Exit with the test code
    os.Exit(code)
}

func TestCountLinesCommandWithNoArgs(t *testing.T) {

	resetGlobals()

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

	resetGlobals()

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

func TestCountLinesCommandWithWrongFilePathExtension(t *testing.T) {

	resetGlobals()

	// create the count lines command
	cmd := count_lines.CountLinesAnalyzer

	// redirect the stdout to a buffer to capture the output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	// set the args
	cmd.SetArgs([]string{"-f", "../../main.go"})

	// execute the count lines command w/ args
	cmd.Execute()

	// check the output
	expectedOutput := "\x1b[31mOnly JavaScript files are accepted.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountLinesCommandWithDirectoryPath(t *testing.T) {

	resetGlobals()

    // create the count lines command
    cmd := count_lines.CountLinesAnalyzer

    // redirect the stdout to a buffer to capture the output
    var stdout bytes.Buffer
    cmd.SetOut(&stdout) 
    cmd.SetErr(&stdout)

    // set the args
    cmd.SetArgs([]string{"-d", "../../javascript-tests"})

    // execute the count lines command w/ args
    if err := cmd.Execute(); err != nil {
        t.Errorf("CountLinesCommand() error = %v, want nil", err)
    }

    // check the output
    expectedOutput := "\x1b[34m Total lines in test.js:\x1b[0m 157\n\x1b[34mTotal lines in directory:\x1b[0m 157\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountLinesCommandWithWrongDirectoryPath(t *testing.T) {
	
	resetGlobals()

	// create the count lines command
	cmd := count_lines.CountLinesAnalyzer

	// redirect the stdout to a buffer to capture the output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	// set the args
	cmd.SetArgs([]string{"-d", "../../main.go"})

	// execute the count lines command w/ args
	cmd.Execute()

	// check the output
	expectedOutput := "\x1b[31mPlease provide a valid directory path.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountLinesCommandWithDirectoryWihtoutJavascriptFiles(t *testing.T) {

	resetGlobals()

	// create the count lines command
	cmd := count_lines.CountLinesAnalyzer

	// redirect the stdout to a buffer to capture the output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	// set the args
	cmd.SetArgs([]string{"-d", "../root"})

	// execute the count lines command w/ args
	cmd.Execute()

	// check the output
	expectedOutput := "\x1b[31mNo JavaScript files found in the provided directory.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountLinesCommandWithValidDirectoryAndGeneratingReportHTML(t *testing.T) {

	resetGlobals()

	// create the count lines command
	cmd := count_lines.CountLinesAnalyzer

	// redirect the stdout to a buffer to capture the output
	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	// set the args
	cmd.SetArgs([]string{"-d", "../../javascript-tests", "-o", "."})

	// execute the count lines command w/ args
	cmd.Execute()

	// check the output
	expectedOutput := "\x1b[1;34mReport generated successfully at report.html\x1b[0m\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountLinesCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}

	// verify if the file was created
	reportPath := "./report.html"
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Errorf("Expected report file to be created at %s, but it does not exist", reportPath)
	}

	// clean up the generated report file
	if err := os.Remove(reportPath); err != nil {
		t.Errorf("Failed to remove the report file: %v", err)
	}
}