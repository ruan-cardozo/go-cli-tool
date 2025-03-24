package count_class_and_functions

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"

	"github.com/spf13/cobra"
)

var countClassAndFunctions analyzer.CountClassesAndFunctionsAnalyzer

var CountClassAndFunctionsCmd = &cobra.Command{
    Use:   "count-class-and-functions",
    Short: "Count classes and functions in a JavaScript file",
    Run: func(cmd *cobra.Command, args []string) {
        
        err := policies.ValidateUserInput(cmd)

        if err {
            return
        }

        if utils.FilePath != "" {
            result := countClassAndFunctions.CountClassesAndFunctionsByFilePath(utils.FilePath)
            fmt.Fprintf(cmd.OutOrStdout(),"%sFunctions:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.Functions)
            fmt.Fprintf(cmd.OutOrStdout(),"%sClasses:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.Classes)
            return
        }

        if utils.DirectoryPath != "" {
            if !policies.ValidateDirectoryPath(utils.DirectoryPath) {
                fmt.Fprintf(cmd.OutOrStdout(),"%sPlease provide a valid directory path.%s", utils.RED, utils.RESET_COLOR)
                return
            }
        }

        result, totalClassesAndFunctions := countClassAndFunctions.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)

        if len(result) == 0 {
            fmt.Fprintf(cmd.OutOrStdout(),"%sNo JavaScript files found in the provided directory.%s", utils.RED, utils.RESET_COLOR)
            return
        }

        concatFuncsAndClasses := fmt.Sprintf("Total Classes: %d, Total Functions: %d", totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions)

        if utils.OutputFilePath != "" {
            templates.SaveResultsToHTML(result, concatFuncsAndClasses, utils.OutputFilePath, utils.COUNT_CLASS_AND_FUNCTIONS,cmd, true, true)
        } else {
            printResults(result, totalClassesAndFunctions,cmd)
        }
    },
}


func printResults(result analyzer.ClassesAndFunctionsMap, totalClassesAndFuncByDirectory analyzer.ClassFuncResult, cmd *cobra.Command) {
    for fileName, result := range result {

        fmt.Fprintf(cmd.OutOrStdout(),"%sFunctions:%s %s %d\n", utils.BLUE, fileName,utils.RESET_COLOR, result.Functions)
        fmt.Fprintf(cmd.OutOrStdout(),"%sClasses:%s %s %d\n", utils.BLUE, fileName, utils.RESET_COLOR, result.Classes)
    }
    fmt.Fprintf(cmd.OutOrStdout(),"%sTotal Classes in directory%s:%d\n", utils.BLUE, utils.RESET_COLOR, totalClassesAndFuncByDirectory.Classes)
    fmt.Fprintf(cmd.OutOrStdout(),"%sTotal Functions in directory%s:%d\n", utils.BLUE, utils.RESET_COLOR, totalClassesAndFuncByDirectory.Functions)
}

func init() {
    countClassAndFunctions = &analyzer.CountClassAndFunctionsImpl{}
    CountClassAndFunctionsCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file")
    CountClassAndFunctionsCmd.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files")
    CountClassAndFunctionsCmd.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output file")
}