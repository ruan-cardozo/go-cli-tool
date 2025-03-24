package count_lines

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"

	"github.com/spf13/cobra"
)

var countLinesAnalyzer analyzer.CountLinesAnalyzer

var CountLinesAnalyzer = &cobra.Command{
    Use:   "count-lines",
    Short: "Count total lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {

        err := policies.ValidateUserInput(cmd)

        if err {
            return
        }

        if utils.FilePath != "" {
            result := countLinesAnalyzer.CountLinesByFilePath(utils.FilePath)
            fmt.Fprintf(cmd.OutOrStdout(),"%sTotal lines:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.TotalLines)
            return
        }

        if utils.DirectoryPath != "" {
            if !policies.ValidateDirectoryPath(utils.DirectoryPath) {
                fmt.Fprintf(cmd.OutOrStdout(),"%sPlease provide a valid directory path.%s", utils.RED, utils.RESET_COLOR)
                return
            }
        }

        result, totalLinesByDirectory := countLinesAnalyzer.CountLinesByDirectory(utils.DirectoryPath)

        totalLinesStr := fmt.Sprintf("Total lines: %d", totalLinesByDirectory.TotalLines)

        if len(result) == 0 {
            fmt.Fprintf(cmd.OutOrStdout(),"%sNo JavaScript files found in the provided directory.%s", utils.RED, utils.RESET_COLOR)
            return
        }

        if utils.OutputFilePath != "" {
            templates.SaveResultsToHTML(result, totalLinesStr, utils.OutputFilePath, utils.COUNT_LINES, cmd, false, false)
        } else {
            printResults(result, totalLinesByDirectory,cmd)
        }
    },
}

func printResults(result analyzer.FilesNameCountLineMap, totalLinesByDirectory analyzer.LineResult, cmd *cobra.Command) {
    for fileName, result := range result {

        fmt.Fprintf(cmd.OutOrStdout(),"%s Total lines in %s:%s %d\n", utils.BLUE, fileName, utils.RESET_COLOR, result.TotalLines)
    }
    fmt.Fprintf(cmd.OutOrStdout(),"%sTotal lines in directory:%s %d\n", utils.BLUE, utils.RESET_COLOR, totalLinesByDirectory.TotalLines)
}

func init() {
    countLinesAnalyzer = &analyzer.CountLinesAnalyzerImpl{}
    CountLinesAnalyzer.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    CountLinesAnalyzer.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    CountLinesAnalyzer.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output file. The tool will generate an HTML report with the results if a directory is provided. If not provided, the tool will print the results to the console. Note: This flag is ignored when the -f flag is used to provide a single file.")
}