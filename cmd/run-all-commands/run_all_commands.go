package run_all_commands

import (
	"encoding/json"
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/utils"
	"os"

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
                generateJSONOutput(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines, 
                    classAndFunctionResult.Classes, classAndFunctionResult.Functions)
                return    
            }
        } else {

            lineResults, totalLines = lineAnalyzer.CountLinesByDirectory(utils.DirectoryPath)
            commentResults, totalComments = commentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
            classFuncResults, totalClassesAndFunctions = classFuncAnalyzer.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)
            
            if utils.SummaryOnly {
                printDirectoryResults(cmd, totalLines.TotalLines, totalComments.CommentLines,
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions)
                return
            }

            if utils.Detailed {
                generateDetailedJSONOutput(cmd, utils.DirectoryPath, lineResults, commentResults, classFuncResults)
                return
            }

            if utils.OutputFilePath == "" {
                printDirectoryResults(cmd, totalLines.TotalLines, totalComments.TotalComments, 
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions)
                return
            } else {
                generateJSONOutput(cmd, utils.DirectoryPath, totalLines.TotalLines, totalComments.CommentLines, totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions)
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

func generateJSONOutput(cmd *cobra.Command, filePath string, lineCount int, commentCount int, classes int, functions int) {
    result := map[string]interface{}{
        "file_path": filePath,
        "lines":     lineCount,
        "comments":  commentCount,
        "classes":   classes,
        "functions": functions,
    }

    if utils.OutputFilePath != "" {
        outputPath := utils.OutputFilePath

        fileInfo, err := os.Stat(outputPath)
        if err == nil && fileInfo.IsDir() {
            outputPath = fmt.Sprintf("%s/analysis_report.json", outputPath)
        }

        file, err := os.Create(outputPath)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error creating JSON file: %s\n", err)
            return
        }
        defer file.Close()

        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ") // Format JSON with indentation
        if err := encoder.Encode(result); err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error writing JSON to file: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "JSON report saved to %s\n", outputPath)
    } else {
        formattedJSON, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error formatting JSON: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(formattedJSON))
    }
}

func generateDetailedJSONOutput(cmd *cobra.Command, directoryPath string, lineResults analyzer.FilesNameCountLineMap, commentResults analyzer.CommentsMap, classFuncResults analyzer.ClassesAndFunctionsMap) {
    detailedResult := map[string]interface{}{
        "directory_path": directoryPath,
        "files":          []map[string]interface{}{},
    }

    for filePath, lineResult := range lineResults {
        fileDetails := map[string]interface{}{
            "file_path": filePath,
            "lines":     lineResult.TotalLines,
            "comments":  commentResults[filePath].CommentLines,
            "classes":   classFuncResults[filePath].Classes,
            "functions": classFuncResults[filePath].Functions,
        }
        detailedResult["files"] = append(detailedResult["files"].([]map[string]interface{}), fileDetails)
    }

    formattedJSON, err := json.MarshalIndent(detailedResult, "", "  ")
    if err != nil {
        fmt.Fprintf(cmd.OutOrStderr(), "Error formatting detailed JSON: %s\n", err)
        return
    }

    fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(formattedJSON))
}

func init() {
    RunAllCommand.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    RunAllCommand.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    RunAllCommand.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Specify the output file path. If omitted, results will be displayed in the terminal.")
    RunAllCommand.Flags().BoolVar(&utils.SummaryOnly, "summary", false, "Show only a summary of the analysis.")
    RunAllCommand.Flags().BoolVar(&utils.Detailed, "detailed", false, "Show detailed analysis, including per-file information.")
}