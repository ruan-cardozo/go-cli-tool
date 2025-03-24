package count_comments

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"os"

	"github.com/spf13/cobra"
)

var CountCommentsCmd = &cobra.Command{
    Use:   "count-comments",
    Short: "Count total comment lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {
        // if !policies.ValidateFilePath(utils.FilePath) {
        //     fmt.Println("\033[1;31mPlease provide the path to the JavaScript file using the -f flag.\033[0m")
        //     os.Exit(1)
        // }

        if !policies.IsJSFileExtension(utils.FilePath) {
            fmt.Println("\033[1;31mOnly JavaScript files are accepted.\033[0m")
            os.Exit(1)
        }

        result := analyzer.CountComments(utils.FilePath)

        fmt.Printf("\033[1;34mComment lines:\033[0m %d\n", result.CommentLines)
    },
}


func init() {
    CountCommentsCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file")
}