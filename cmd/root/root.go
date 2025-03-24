package root

import (
	"fmt"
	count_class_and_functions "go-cli-tool/cmd/count-class-and-functions"
	count_comments "go-cli-tool/cmd/count-comments"
	count_lines "go-cli-tool/cmd/count-lines"

	"github.com/spf13/cobra"
)


var RootCmd = &cobra.Command{
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

func RootCommand() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(RootCmd.OutOrStdout(), "%v\n", err)
		return
	}
}

func init() {
	RootCmd.AddCommand(count_lines.CountLinesAnalyzer)
	RootCmd.AddCommand(count_comments.CountCommentsCmd)
	RootCmd.AddCommand(count_class_and_functions.CountClassAndFunctionsCmd)
}