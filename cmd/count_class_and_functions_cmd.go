package cmd

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"os"

	"github.com/spf13/cobra"
)

var countClassAndFunctionsCmd = &cobra.Command{
    Use:   "count-class-and-functions",
    Short: "Count classes and functions in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {
        if !policies.ValidateFilePath(filePath) {
            fmt.Println("\033[1;31mPlease provide the path to the JavaScript file using the -f flag.\033[0m")
            os.Exit(1)
        }
        result := analyzer.CountClassesAndFunctions(filePath)

        fmt.Printf("\033[1;34mFunctions:\033[0m %d\n", result.Functions)
        fmt.Printf("\033[1;34mClasses:\033[0m %d\n", result.Classes)
    },
}

func init() {
    countClassAndFunctionsCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the JavaScript file")
    rootCmd.AddCommand(countClassAndFunctionsCmd)
}