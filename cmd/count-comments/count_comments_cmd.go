package count_comments

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"

	"github.com/spf13/cobra"
)

var countCommentsAnalyzer analyzer.CountCommentsAnalyzer

var CountCommentsCmd = &cobra.Command{
    Use:   "count-comments",
    Short: "Count total comment lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {

        err := policies.ValidateUserInput(cmd)

        if err {
            return
        }

        if utils.FilePath != "" {
            result := countCommentsAnalyzer.CountCommentsByFilePath(utils.FilePath)
            fmt.Fprintf(cmd.OutOrStdout(),"%sTotal comments:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.CommentLines)
            return
        }

        if utils.DirectoryPath != "" {
            if !policies.ValidateDirectoryPath(utils.DirectoryPath) {
                fmt.Fprintf(cmd.OutOrStdout(),"%sPlease provide a valid directory path.%s", utils.RED, utils.RESET_COLOR)
                return
            }
        }

        result, totalCommentsByDirectory := countCommentsAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)

        if len(result) == 0 {
            fmt.Fprintf(cmd.OutOrStdout(),"%sNo JavaScript files found in the provided directory.%s", utils.RED, utils.RESET_COLOR)
            return
        }

        totalLinesStr := fmt.Sprintf("Total comments: %d", totalCommentsByDirectory.TotalComments)

        if utils.OutputFilePath != "" {
            templates.SaveResultsToHTML(result, totalLinesStr, utils.OutputFilePath, utils.COUNT_COMMENTS, cmd, false, false)
        } else {
            printResults(result, totalCommentsByDirectory,cmd)
        }
    },
}

func printResults(result analyzer.CommentsMap, totalCommentsByDirectory analyzer.CommentResult, cmd *cobra.Command) {
    for fileName, result := range result {
        fmt.Fprintf(cmd.OutOrStdout(),"%s Comment lines in %s:%s %d\n", utils.BLUE, fileName, utils.RESET_COLOR, result.CommentLines)
    }
    fmt.Fprintf(cmd.OutOrStdout(),"%sTotal Comments in directory:%s %d\n", utils.BLUE, utils.RESET_COLOR, totalCommentsByDirectory.TotalComments)
}

func init() {
    countCommentsAnalyzer = &analyzer.CountCommentsAnalyzerImpl{}
    CountCommentsCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    CountCommentsCmd.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    CountCommentsCmd.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output file. The tool will generate an HTML report with the results if a directory is provided. If not provided, the tool will print the results to the console. Note: This flag is ignored when the -f flag is used to provide a single file.")
}