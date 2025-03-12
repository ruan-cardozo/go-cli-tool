package cmd

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"os"

	"github.com/spf13/cobra"
)

var countCommentsCmd = &cobra.Command{
    Use:   "count-comments",
    Short: "Count total comment lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {
        if !policies.ValidateFilePath(filePath) {
            fmt.Println("\033[1;31mPlease provide the path to the JavaScript file using the -f flag.\033[0m")
            os.Exit(1)
        }

        if !policies.IsJSFileExtension(filePath) {
            fmt.Println("\033[1;31mOnly JavaScript files are accepted.\033[0m")
            os.Exit(1)
        }

        result := analyzer.CountComments(filePath)

        fmt.Printf("\033[1;34mComment lines:\033[0m %d\n", result.CommentLines)
    },
}


func init() {
    countCommentsCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the JavaScript file")
    rootCmd.AddCommand(countCommentsCmd)
}