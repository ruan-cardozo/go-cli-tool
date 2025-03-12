package cmd

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"os"

	"github.com/spf13/cobra"
)

var filePath string

var countLinesCmd = &cobra.Command{
    Use:   "count-lines",
    Short: "Count total lines and comment lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {
        if !policies.ValidateFilePath(filePath) {
            fmt.Printf("%sPlease provide the path to the JavaScript file using the -f flag or use.%s", utils.Red, utils.ResetColor)
            os.Exit(1)
        }

        if !policies.IsJSFileExtension(filePath) {
            fmt.Printf("%sOnly JavaScript files are accepted.%s", utils.Red, utils.ResetColor)
            os.Exit(1)
        }

        result := analyzer.CountLines(filePath)

        fmt.Printf("%sTotal lines:%s %d\n", utils.Blue, utils.ResetColor, result.TotalLines)
    },
}

func init() {
    countLinesCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the JavaScript file or directory")
    countLinesCmd.MarkFlagRequired("file")
    rootCmd.AddCommand(countLinesCmd)
}