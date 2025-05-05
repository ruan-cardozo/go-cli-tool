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
- Indentation analysis

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
        
        indentationAnalyzer := &analyzer.IdentationAnalyzerImpl{}
        var indentResults map[string]interface{}
        
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
            
            tempFilePath := utils.FilePath
            tempDirPath := utils.DirectoryPath

            utils.DirectoryPath = ""
            indentResults, _ = indentationAnalyzer.IdentationByFilePath()
            utils.FilePath = tempFilePath
            utils.DirectoryPath = tempDirPath

            if utils.OutputFilePath == "" {
                printFileResults(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines, 
                    classAndFunctionResult.Classes, classAndFunctionResult.Functions, indentResults)
                return
            } else {
                generateJSONOutput(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines, 
                    classAndFunctionResult.Classes, classAndFunctionResult.Functions, indentResults)
                return    
            }
        } else {

            lineResults, totalLines = lineAnalyzer.CountLinesByDirectory(utils.DirectoryPath)
            commentResults, totalComments = commentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
            classFuncResults, totalClassesAndFunctions = classFuncAnalyzer.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)
            
            tempFilePath := utils.FilePath
            tempDirPath := utils.DirectoryPath
            utils.FilePath = ""

            indentResults, _ = indentationAnalyzer.IdentationByFilePath()
            utils.FilePath = tempFilePath
            utils.DirectoryPath = tempDirPath
            
            if utils.SummaryOnly {

                if utils.OutputFilePath != "" {

                    var indentSummary map[string]interface{} = nil
                    
                    if files, ok := indentResults["files"].([]map[string]interface{}); ok && len(files) > 0 {

                        totalMaxIndent := 0
                        totalAvgIndent := 0.0
                        spacesCount := 0
                        tabsCount := 0
                        mixedCount := 0
                        
                        for _, file := range files {
                            if stats, ok := file["stats"].(analyzer.IndentResult); ok {
                                totalMaxIndent += stats.MaxIndentLevel
                                totalAvgIndent += stats.AverageIndentLevel
                                if stats.UsesSpaces {
                                    spacesCount++
                                }
                                if stats.UsesTabs {
                                    tabsCount++
                                }
                                if stats.MixedIndentation {
                                    mixedCount++
                                }
                            }
                        }
                        
                        fileCount := len(files)
                        if fileCount > 0 {
                            indentSummary = map[string]interface{}{
                                "avgMaxIndentLevel": float64(totalMaxIndent) / float64(fileCount),
                                "avgIndentLevel":    totalAvgIndent / float64(fileCount),
                                "filesUsingSpaces":  spacesCount,
                                "filesUsingTabs":    tabsCount,
                                "filesMixedIndent":  mixedCount,
                            }
                        }
                    }
                    
                    generateSummaryJSONOutput(cmd, utils.DirectoryPath, 
                        totalLines.TotalLines, 
                        totalComments.CommentLines,
                        totalClassesAndFunctions.Classes, 
                        totalClassesAndFunctions.Functions, 
                        indentSummary)
                    return
                }

                printDirectoryResults(cmd, totalLines.TotalLines, totalComments.CommentLines,
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions, indentResults)
                return
            }

            if utils.Detailed {
                generateDetailedJSONOutput(cmd, utils.DirectoryPath, lineResults, commentResults, 
                    classFuncResults, indentResults)
                return
            }

            if utils.OutputFilePath == "" {
                printDirectoryResults(cmd, totalLines.TotalLines, totalComments.TotalComments, 
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions, indentResults)
                return
            } else {
                generateJSONOutput(cmd, utils.DirectoryPath, totalLines.TotalLines, totalComments.CommentLines, 
                    totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions, indentResults)
            }
        }
    },
}

func printFileResults(cmd *cobra.Command, filePath string, lines, comments, classes, functions int, indentResults map[string]interface{}) {
    fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Analysis Results for %s ===%s\n", 
        utils.BLUE, filePath, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, lines, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, comments, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, classes, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, functions, utils.RESET_COLOR)
    
    if stats, ok := indentResults["stats"].(analyzer.IndentResult); ok {
        fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Indentation Analysis ===%s\n", utils.BLUE, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Max Indent Level: %s%d%s\n", utils.GREEN, stats.MaxIndentLevel, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Average Indent Level: %s%.2f%s\n", utils.GREEN, stats.AverageIndentLevel, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Uses Spaces: %s%t%s\n", utils.GREEN, stats.UsesSpaces, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Uses Tabs: %s%t%s\n", utils.GREEN, stats.UsesTabs, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Mixed Indentation: %s%t%s\n", utils.GREEN, stats.MixedIndentation, utils.RESET_COLOR)
    }
}

func printDirectoryResults(cmd *cobra.Command, lines, comments, classes, functions int, indentResults map[string]interface{}) {
    fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Directory Analysis Summary ===%s\n", 
        utils.BLUE, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, lines, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, comments, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, classes, utils.RESET_COLOR)
    fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, functions, utils.RESET_COLOR)
    
    fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Indentation Analysis Summary ===%s\n", utils.BLUE, utils.RESET_COLOR)
    if files, ok := indentResults["files"].([]map[string]interface{}); ok && len(files) > 0 {

        totalMaxIndent := 0
        totalAvgIndent := 0.0
        spacesCount := 0
        tabsCount := 0
        mixedCount := 0
        
        for _, file := range files {
            if stats, ok := file["stats"].(analyzer.IndentResult); ok {
                totalMaxIndent += stats.MaxIndentLevel
                totalAvgIndent += stats.AverageIndentLevel
                if stats.UsesSpaces {
                    spacesCount++
                }
                if stats.UsesTabs {
                    tabsCount++
                }
                if stats.MixedIndentation {
                    mixedCount++
                }
            }
        }
        
        fileCount := len(files)
        if fileCount > 0 {
            fmt.Fprintf(cmd.OutOrStdout(), "Avg Max Indent Level: %s%.2f%s\n", 
                utils.GREEN, float64(totalMaxIndent)/float64(fileCount), utils.RESET_COLOR)
            fmt.Fprintf(cmd.OutOrStdout(), "Avg Indent Level: %s%.2f%s\n", 
                utils.GREEN, totalAvgIndent/float64(fileCount), utils.RESET_COLOR)
            fmt.Fprintf(cmd.OutOrStdout(), "Files Using Spaces: %s%d%s\n", 
                utils.GREEN, spacesCount, utils.RESET_COLOR)
            fmt.Fprintf(cmd.OutOrStdout(), "Files Using Tabs: %s%d%s\n", 
                utils.GREEN, tabsCount, utils.RESET_COLOR)
            fmt.Fprintf(cmd.OutOrStdout(), "Files With Mixed Indentation: %s%d%s\n", 
                utils.GREEN, mixedCount, utils.RESET_COLOR)
        }
    }
}

func generateJSONOutput(cmd *cobra.Command, filePath string, lineCount int, commentCount int, 
                        classes int, functions int, indentResults map[string]interface{}) {

    if indentResults != nil {
        delete(indentResults, "path")
        delete(indentResults, "filename")
    }              

    result := map[string]interface{}{
        "file_path": filePath,
        "lines":     lineCount,
        "comments":  commentCount,
        "classes":   classes,
        "functions": functions,
        "indentation": indentResults,
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
        encoder.SetIndent("", "  ")
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

func generateDetailedJSONOutput(cmd *cobra.Command, directoryPath string, lineResults analyzer.FilesNameCountLineMap, 
                              commentResults analyzer.CommentsMap, classFuncResults analyzer.ClassesAndFunctionsMap,
                              indentResults map[string]interface{}) {
    
    fileIndentData := make(map[string]interface{})
    if files, ok := indentResults["files"].([]map[string]interface{}); ok {
        for _, fileData := range files {
            if filename, ok := fileData["filename"].(string); ok {
                fileIndentData[filename] = fileData["stats"]
            }
        }
    }

    fileDetails := make([]map[string]interface{}, 0, len(lineResults))
    for filename, lineResult := range lineResults {
        fileInfo := map[string]interface{}{
            "filename":  filename,
            "metrics": map[string]interface{}{
                "lines":     lineResult.TotalLines,
                "comments":  commentResults[filename].CommentLines,
                "classes":   classFuncResults[filename].Classes,
                "functions": classFuncResults[filename].Functions,
            },
        }

        if indentData, ok := fileIndentData[filename]; ok {
            fileInfo["indentation"] = indentData
        }
        
        fileDetails = append(fileDetails, fileInfo)
    }

    totalLines := 0
    totalComments := 0
    totalClasses := 0
    totalFunctions := 0
    
    for filePath := range lineResults {
        totalLines += lineResults[filePath].TotalLines
        totalComments += commentResults[filePath].CommentLines
        totalClasses += classFuncResults[filePath].Classes
        totalFunctions += classFuncResults[filePath].Functions
    }
    
    detailedResult := map[string]interface{}{
        "directory_path": directoryPath,
        "summary": map[string]interface{}{
            "total_files":  len(lineResults),
            "total_lines":  totalLines,
            "total_comments": totalComments,
            "total_classes": totalClasses,
            "total_functions": totalFunctions,
        },
        "files": fileDetails,
    }

    if utils.OutputFilePath != "" {
        outputPath := utils.OutputFilePath
        
        fileInfo, err := os.Stat(outputPath)
        if err == nil && fileInfo.IsDir() {
            outputPath = fmt.Sprintf("%s/detailed_analysis.json", outputPath)
        }
        
        file, err := os.Create(outputPath)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error creating JSON file: %s\n", err)
            return
        }
        defer file.Close()
        
        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(detailedResult); err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error writing JSON to file: %s\n", err)
            return
        }
        
        fmt.Fprintf(cmd.OutOrStdout(), "Detailed JSON report saved to %s\n", outputPath)
    } else {

        formattedJSON, err := json.MarshalIndent(detailedResult, "", "  ")
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error formatting detailed JSON: %s\n", err)
            return
        }
        
        fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(formattedJSON))
    }
}

func generateSummaryJSONOutput(cmd *cobra.Command, directoryPath string, lineCount, commentCount, classes, functions int, indentSummary map[string]interface{}) {

    summaryObject := map[string]interface{}{
        "lines":      lineCount,
        "comments":   commentCount,
        "classes":    classes,
        "functions":  functions,
    }

    if indentSummary != nil {
        summaryObject["indentation"] = indentSummary
    }
    
    summaryResult := map[string]interface{}{
        "directory_path": directoryPath,
        "summary": summaryObject,
    }

    if utils.OutputFilePath != "" {
        outputPath := utils.OutputFilePath

        fileInfo, err := os.Stat(outputPath)
        if err == nil && fileInfo.IsDir() {
            outputPath = fmt.Sprintf("%s/summary_report.json", outputPath)
        }

        file, err := os.Create(outputPath)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error creating JSON file: %s\n", err)
            return
        }
        defer file.Close()

        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(summaryResult); err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error writing JSON to file: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "Summary JSON report saved to %s\n", outputPath)
    }
}

func init() {
    RunAllCommand.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    RunAllCommand.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    RunAllCommand.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Specify the output file path. If omitted, results will be displayed in the terminal.")
    RunAllCommand.Flags().BoolVar(&utils.SummaryOnly, "summary", false, "Show only a summary of the analysis.")
    RunAllCommand.Flags().BoolVar(&utils.Detailed, "detailed", false, "Show detailed analysis, including per-file information.")
}