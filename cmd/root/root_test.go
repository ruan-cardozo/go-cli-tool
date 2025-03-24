package root_test

import (
	"bytes"
	"go-cli-tool/cmd/root"
	"strings"
	"testing"
)

func TestRootCommand(t *testing.T) {
    // create the root command
    rootCmd := root.RootCmd

    // redirect the stdout to a buffer to capture the output
    var stdout bytes.Buffer
    rootCmd.SetOut(&stdout)

    // set the args

    // execute the root command w/ args
    err := rootCmd.Execute()

    // check the error
    if err != nil {
        t.Errorf("RootCommand() error = %v, want nil", err)
    }

    // check the output
    expectedOutput := "\x1b[34mA CLI tool for JavaScript code analysis!\x1b[0m \nWith this tool, you can easily analyze your JavaScript files and get detailed information, including:\n    - \x1b[32mTotal line count\x1b[0m\n    - \x1b[32mNumber of comment lines\x1b[0m\n    - \x1b[32mNumber of functions\x1b[0m\n    - \x1b[32mNumber of classes\x1b[0m\nSimplify your code analysis. Use go-cli-tool!\n\nUsage:\n  go-cli-tool [command]\n\nAvailable Commands:\n  completion                Generate the autocompletion script for the specified shell\n  count-class-and-functions Count classes and functions in a JavaScript file\n  count-comments            Count total comment lines in a JavaScript file\n  count-lines               Count total lines in a JavaScript file\n  help                      Help about any command\n\nFlags:\n  -h, --help   help for go-cli-tool\n\nUse \"go-cli-tool [command] --help\" for more information about a command."

	actualOutput := strings.TrimSpace(stdout.String())
    expectedOutput = strings.TrimSpace(expectedOutput)

    if actualOutput != expectedOutput {
        t.Errorf("RootCommand() = %v, want %v", actualOutput, expectedOutput)
		t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}

func TestRootCommandWithError(t *testing.T) {
	// Cria um novo comando root
	rootCmd := root.RootCmd

	// Redireciona o stdout para um buffer para capturar a saída
	var stdout bytes.Buffer
	rootCmd.SetOut(&stdout)
    rootCmd.SetErr(&stdout)

	// Define um argumento inválido para simular um erro
	rootCmd.SetArgs([]string{"invalid-command"})

	// Executa o comando root
	root.RootCommand()

	// Verifica se a saída contém a mensagem de erro esperada
	expectedOutput := "Error: unknown command \"invalid-command\" for \"go-cli-tool\"\nRun 'go-cli-tool --help' for usage.\nunknown command \"invalid-command\" for \"go-cli-tool\""
	actualOutput := strings.TrimSpace(stdout.String())

    if actualOutput != expectedOutput {
        t.Errorf("RootCommandWithError() = %v, want %v", actualOutput, expectedOutput)
        t.Logf("Actual Output: %q", actualOutput)
        t.Logf("Expected Output: %q", expectedOutput)
    }
}