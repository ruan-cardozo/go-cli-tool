package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
		Use:   "go-cli-tool",
		Short: "CLI for JavaScript code analysis",
		Long: fmt.Sprintf(`%sA CLI tool for JavaScript code analysis!%s 
With this tool, you can easily analyze your JavaScript files and get detailed information, including:
    - %sTotal line count%s
    - %sNumber of comment lines%s
    - %sNumber of functions%s
    - %sNumber of classes%s
Simplify your code analysis. Use go-cli-tool!`, 
			"\033[34m", "\033[0m", // Azul no título
			"\033[32m", "\033[0m", // Verde nos itens
			"\033[32m", "\033[0m",
			"\033[32m", "\033[0m",
			"\033[32m", "\033[0m"),
	}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}