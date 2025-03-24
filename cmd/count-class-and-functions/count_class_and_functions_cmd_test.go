package count_class_and_functions_test

import (
	"bytes"
	count_class_and_functions "go-cli-tool/cmd/count-class-and-functions"
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

func TestCountClassAndFunctionsCommandWithNoArgs(t *testing.T) {

	tests.ResetGlobals()

    cmd := count_class_and_functions.CountClassAndFunctionsCmd

    var stdout bytes.Buffer
    cmd.SetOut(&stdout)

    cmd.SetArgs([]string{})

	cmd.Execute()

    expectedOutput := "Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountClassAndFunctionsWithFilePath(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_class_and_functions.CountClassAndFunctionsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)

	cmd.SetArgs([]string{"-f", "../../javascript-tests/test.js"})

	cmd.Execute()

	expectedOutput := "\x1b[34mFunctions:\x1b[0m 20\n\x1b[34mClasses:\x1b[0m 1\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountClassAndFunctionsCommandWithWrongFilePathExtension(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_class_and_functions.CountClassAndFunctionsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-f", "../../main.go"})

	cmd.Execute()

	expectedOutput := "\x1b[31mOnly JavaScript files are accepted.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountClassAndFunctionsCommandWithDirectoryPath(t *testing.T) {

	tests.ResetGlobals()

	cmd := count_class_and_functions.CountClassAndFunctionsCmd

    var stdout bytes.Buffer
    cmd.SetOut(&stdout) 
    cmd.SetErr(&stdout)

    cmd.SetArgs([]string{"-d", "../../javascript-tests"})

    if err := cmd.Execute(); err != nil {
        t.Errorf("CountClassAndFunctions() error = %v, want nil", err)
    }

    expectedOutput := "\x1b[34mFunctions:test.js \x1b[0m 20\n\x1b[34mClasses:test.js \x1b[0m 1\n\x1b[34mTotal Classes in directory\x1b[0m:1\n\x1b[34mTotal Functions in directory\x1b[0m:20\n"

    actualOutput := stdout.String()

    if actualOutput != expectedOutput {
        t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestCountClassAndFunctionsCommandWithWrongDirectoryPath(t *testing.T) {
	
	tests.ResetGlobals()

    cmd := count_class_and_functions.CountClassAndFunctionsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../../main.go"})

	cmd.Execute()

	expectedOutput := "\x1b[31mPlease provide a valid directory path.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountClassAndFunctionsCommandWithDirectoryWihtoutJavascriptFiles(t *testing.T) {

	tests.ResetGlobals()

    cmd := count_class_and_functions.CountClassAndFunctionsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../root"})

	cmd.Execute()

	expectedOutput := "\x1b[31mNo JavaScript files found in the provided directory.\x1b[0m"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
		t.Logf("Expected Output: %q", expectedOutput)
	}
}

func TestCountClassAndFunctionsCommandWithValidDirectoryAndGeneratingReportHTML(t *testing.T) {

	tests.ResetGlobals()

    cmd := count_class_and_functions.CountClassAndFunctionsCmd

	var stdout bytes.Buffer
	cmd.SetOut(&stdout) 

	cmd.SetArgs([]string{"-d", "../../javascript-tests", "-o", "."})

	cmd.Execute()

	expectedOutput := "\x1b[1;34mReport generated successfully at report.html\x1b[0m\n"

	actualOutput := stdout.String()

	if actualOutput != expectedOutput {
		t.Errorf("CountClassAndFunctionsCommand() = %v, want %v", actualOutput, expectedOutput)
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