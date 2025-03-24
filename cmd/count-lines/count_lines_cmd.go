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

        err := validateUserInput(cmd)

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

        if len(result) == 0 {
            fmt.Fprintf(cmd.OutOrStdout(),"%sNo JavaScript files found in the provided directory.%s", utils.RED, utils.RESET_COLOR)
            return
        }

        if utils.OutputFilePath != "" {
            templates.SaveResultsToHTML(result, totalLinesByDirectory, utils.OutputFilePath, utils.COUNT_LINES,cmd)
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

func validateUserInput(cmd *cobra.Command) bool {

    err := false

    if utils.FilePath == "" && utils.DirectoryPath == "" {
        cmd.Println("Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.")
        err = true
    }

    if utils.FilePath != "" {
       err = validateFilePath(err,cmd)
    }

    return err
}

func validateFilePath(err bool, cmd *cobra.Command) bool {

    if !policies.IsJSFileExtension(utils.FilePath) {
        fmt.Fprintf(cmd.OutOrStdout(), "%sOnly JavaScript files are accepted.%s", utils.RED, utils.RESET_COLOR)
        err = true
    }

    return err
}

func init() {
    countLinesAnalyzer = &analyzer.CountLinesAnalyzerImpl{}
    CountLinesAnalyzer.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    CountLinesAnalyzer.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    CountLinesAnalyzer.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output file. The tool will generate an HTML report with the results if a directory is provided. If not provided, the tool will print the results to the console. Note: This flag is ignored when a single file is provided.")
}