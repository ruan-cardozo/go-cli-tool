package cmd

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"
	"os"

	"github.com/spf13/cobra"
)

var countLinesCmd = &cobra.Command{
    Use:   "count-lines",
    Short: "Count total lines and comment lines in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {

        validateUserInput()

        if filePath != "" {
            result := analyzer.CountLinesByFilePath(filePath)
            fmt.Printf("%sTotal lines:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.TotalLines)
        }

        result, totalLinesByDirectory := analyzer.CountLinesByDirectory(directoryPath)

        if outputFilePath != "" {
            templates.SaveResultsToHTML(result, totalLinesByDirectory, outputFilePath, utils.COUNT_LINES)
        } else {
            printResults(result, totalLinesByDirectory)
        }
    },
}

func printResults(result analyzer.FilesNameCountLineMap, totalLinesByDirectory analyzer.LineResult) {
    for fileName, result := range result {

        fmt.Printf("%s Total lines in %s:%s %d\n", utils.BLUE, fileName, utils.RESET_COLOR, result.TotalLines)
    }
    fmt.Printf("%sTotal lines in directory:%s %d\n", utils.BLUE, utils.RESET_COLOR, totalLinesByDirectory.TotalLines)
}

func validateUserInput() {

    if filePath == "" && directoryPath == "" {
        fmt.Printf("Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.\n")
        os.Exit(1)
    }

    if filePath != "" {
        validateFilePath()
    }
}

func validateFilePath() {
    if !policies.ValidateFilePath(filePath) {
        fmt.Printf("%sPlease provide the path to the JavaScript file using the -f flag or use.%s", utils.RED, utils.RESET_COLOR)
        os.Exit(1)
    }

    if !policies.IsJSFileExtension(filePath) {
        fmt.Printf("%sOnly JavaScript files are accepted.%s", utils.RED, utils.RESET_COLOR)
        os.Exit(1)
    }
}

func init() {
    countLinesCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    countLinesCmd.Flags().StringVarP(&directoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    countLinesCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Path to the output file. The tool will generate an HTML report with the results. If not provided, the tool will print the results to the console.")
    rootCmd.AddCommand(countLinesCmd)
}