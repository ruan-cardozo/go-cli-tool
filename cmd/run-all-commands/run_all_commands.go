package run_all_commands

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/utils"

	"github.com/spf13/cobra"
)

var RunAllCommand = &cobra.Command{
    Use:   "analyze",
    Short: "Comprehensive JavaScript code analysis in a single operation",
    Long: `Execute a complete analysis of JavaScript code with a single command.
    
This command performs multiple analyses simultaneously including:
- Line count analysis
- Comment count analysis
- Function and class count analysis

Results are presented in terminal or json output, providing a complete overview
of your JavaScript codebase. Use flags to customize the analysis and output format.`,
    Run: func(cmd *cobra.Command, args []string) {

        if utils.FilePath == "" && utils.DirectoryPath == "" {
            fmt.Fprintf(cmd.OutOrStderr(), "%sError: You must provide either a file path (-f) or directory path (-d).%s\n",
                utils.RED, utils.RESET_COLOR)
            return
        }
    
        lineAnalyzer := &analyzer.CountLinesAnalyzerImpl{}
        var lineResults analyzer.FilesNameCountLineMap
        var totalLines analyzer.LineResult

        commentAnalyzer := &analyzer.CountCommentsAnalyzerImpl{}
        var commentResults analyzer.CommentsMap
        var totalComments analyzer.CommentResult
        
        classFuncAnalyzer := &analyzer.CountClassAndFunctionsImpl{}
        var classFuncResults analyzer.ClassesAndFunctionsMap
        var totalClassesAndFunctions analyzer.ClassFuncResult
        
        if utils.FilePath != "" {

            lineResults = make(analyzer.FilesNameCountLineMap)
            commentResults = make(analyzer.CommentsMap)
            classFuncResults = make(analyzer.ClassesAndFunctionsMap)
            
            lineCount := lineAnalyzer.CountLinesByFilePath(utils.FilePath)

            lineResults[utils.FilePath] = lineCount
            
            commentCount := commentAnalyzer.CountCommentsByFilePath(utils.FilePath)

            commentResults[utils.FilePath] = commentCount
            classAndFunctionResult := classFuncAnalyzer.CountClassesAndFunctionsByFilePath(utils.FilePath)
            classFuncResults[utils.FilePath] = classAndFunctionResult

            if utils.OutputFilePath == "" {
                printFileResults(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines, 
                    classAndFunctionResult.Classes, classAndFunctionResult.Functions)
                return
            } else {
                fmt.Fprintf(cmd.OutOrStdout(), "%sSaida estilo json em desenvolvimento...%s\n", utils.BLUE, utils.RESET_COLOR)
            }
        } else {

            _, totalLines = lineAnalyzer.CountLinesByDirectory(utils.DirectoryPath)
            _, totalComments = commentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
            _, totalClassesAndFunctions = classFuncAnalyzer.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)
            
            if utils.OutputFilePath == "" {
                printDirectoryResults(cmd, totalLines.TotalLines, totalComments.TotalComments, 
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions)
                return
            } else {
                fmt.Fprintf(cmd.OutOrStdout(), "%sSaida estilo json em desenvolvimento...%s\n", utils.BLUE, utils.RESET_COLOR)
            }
        }
    },
}

func printFileResults(cmd *cobra.Command, filePath string, lines, comments, classes, functions int) {
    fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Analysis Results for %s ===%s\n", 
        utils.BLUE, filePath, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, lines, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, comments, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, classes, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, functions, utils.RESET_COLOR)
}

func printDirectoryResults(cmd *cobra.Command, lines, comments, classes, functions int) {
    fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Directory Analysis Summary ===%s\n", 
        utils.BLUE, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, lines, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, comments, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, classes, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, functions, utils.RESET_COLOR)
}

func init() {
    RunAllCommand.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    RunAllCommand.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    RunAllCommand.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Specify the output file path. If omitted, results will be displayed in the terminal. For directory analysis, only a summary will be shown if no output file is provided.")
}