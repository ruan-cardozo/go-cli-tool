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

        lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer := initializeAnalyzers()

        if utils.FilePath != "" {
            handleFileAnalysis(cmd, lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer)
        } else {
            handleDirectoryAnalysis(cmd, lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer)
        }
    },
}

// Helper function to initialize analyzers
func initializeAnalyzers() (*analyzer.CountLinesAnalyzerImpl, *analyzer.CountCommentsAnalyzerImpl, *analyzer.CountClassAndFunctionsImpl, *analyzer.IdentationAnalyzerImpl, *analyzer.CountDependenciesAnalyzerImpl) {
    return &analyzer.CountLinesAnalyzerImpl{},
        &analyzer.CountCommentsAnalyzerImpl{},
        &analyzer.CountClassAndFunctionsImpl{},
        &analyzer.IdentationAnalyzerImpl{},
        &analyzer.CountDependenciesAnalyzerImpl{}
}

// Handles file-level analysis
func handleFileAnalysis(cmd *cobra.Command, lineAnalyzer *analyzer.CountLinesAnalyzerImpl, commentAnalyzer *analyzer.CountCommentsAnalyzerImpl, classFuncAnalyzer *analyzer.CountClassAndFunctionsImpl, indentationAnalyzer *analyzer.IdentationAnalyzerImpl, dependenciesAnalyzer *analyzer.CountDependenciesAnalyzerImpl) {
    lineResults := make(analyzer.FilesNameCountLineMap)
    commentResults := make(analyzer.CommentsMap)
    classFuncResults := make(analyzer.ClassesAndFunctionsMap)

    lineCount := lineAnalyzer.CountLinesByFilePath(utils.FilePath)
    lineResults[utils.FilePath] = lineCount

    commentCount := commentAnalyzer.CountCommentsByFilePath(utils.FilePath)
    commentResults[utils.FilePath] = commentCount

    classAndFunctionResult := classFuncAnalyzer.CountClassesAndFunctionsByFilePath(utils.FilePath)
    classFuncResults[utils.FilePath] = classAndFunctionResult

    tempFilePath := utils.FilePath
    tempDirPath := utils.DirectoryPath
    utils.DirectoryPath = ""
    indentResults, _ := indentationAnalyzer.IdentationByFilePath()
    utils.FilePath = tempFilePath
    utils.DirectoryPath = tempDirPath

    dependencieResultMap, _ := dependenciesAnalyzer.CountDependenciesByFilePath(utils.FilePath)

    if utils.OutputFilePath == "" {
        printFileResults(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines,
            classAndFunctionResult.Classes, classAndFunctionResult.Functions, indentResults, dependencieResultMap)
    } else {
        generateJSONOutput(cmd, utils.FilePath, lineCount.TotalLines, commentCount.CommentLines,
            classAndFunctionResult.Classes, classAndFunctionResult.Functions, indentResults, dependencieResultMap)
    }
}

// Handles directory-level analysis
func handleDirectoryAnalysis(cmd *cobra.Command, lineAnalyzer *analyzer.CountLinesAnalyzerImpl, commentAnalyzer *analyzer.CountCommentsAnalyzerImpl, classFuncAnalyzer *analyzer.CountClassAndFunctionsImpl, indentationAnalyzer *analyzer.IdentationAnalyzerImpl, dependenciesAnalyzer *analyzer.CountDependenciesAnalyzerImpl) {
    lineResults, totalLines := lineAnalyzer.CountLinesByDirectory(utils.DirectoryPath)
    commentResults, totalComments := commentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
    classFuncResults, totalClassesAndFunctions := classFuncAnalyzer.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)

    tempFilePath := utils.FilePath
    tempDirPath := utils.DirectoryPath
    utils.FilePath = ""
    indentResults, _ := indentationAnalyzer.IdentationByFilePath()
    utils.FilePath = tempFilePath
    utils.DirectoryPath = tempDirPath

    dependencieResultMap, _ := dependenciesAnalyzer.CountDependenciesByDirectory(utils.DirectoryPath)

    if utils.Detailed {
        generateDetailedJSONOutput(cmd, utils.DirectoryPath, lineResults, commentResults,
            classFuncResults, indentResults, dependencieResultMap)
    } else if utils.OutputFilePath == "" {
        printDirectoryResults(cmd, totalLines.TotalLines, totalComments.TotalComments,
            totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions, indentResults)
    } else {
        generateJSONOutput(cmd, utils.DirectoryPath, totalLines.TotalLines, totalComments.TotalComments,
            totalClassesAndFunctions.Classes, totalClassesAndFunctions.Functions, indentResults, dependencieResultMap)
    }
}

// Consolidates dependencies into a single map
func consolidateDependencies(dependenciesResults map[string]interface{}) map[string]interface{} {
    if dependenciesResults != nil {
        consolidatedDeps := map[string]interface{}{
            "total_dependencies": 0,
            "dependencies":       []string{},
            "native_modules":     []string{},
        }

        uniqueDeps := make(map[string]struct{})
        uniqueNativeModules := make(map[string]struct{})

        for _, fileDeps := range dependenciesResults {
            if deps, ok := fileDeps.(map[string]interface{}); ok {
                if depList, exists := deps["dependencies"].([]string); exists {
                    for _, dep := range depList {
                        uniqueDeps[dep] = struct{}{}
                    }
                }
                if nativeList, exists := deps["native_modules"].([]string); exists {
                    for _, native := range nativeList {
                        uniqueNativeModules[native] = struct{}{}
                    }
                }
            }
        }

        for dep := range uniqueDeps {
            consolidatedDeps["dependencies"] = append(consolidatedDeps["dependencies"].([]string), dep)
        }
        for native := range uniqueNativeModules {
            consolidatedDeps["native_modules"] = append(consolidatedDeps["native_modules"].([]string), native)
        }

        consolidatedDeps["total_dependencies"] = len(uniqueDeps)
        return consolidatedDeps
    }
    return nil
}

// Outputs JSON to a file or terminal
func outputJSON(cmd *cobra.Command, outputPath string, data interface{}) {
    if outputPath != "" {
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
        if err := encoder.Encode(data); err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error writing JSON to file: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "JSON report saved to %s\n", outputPath)
    } else {
        formattedJSON, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error formatting JSON: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(formattedJSON))
    }
}

func generateJSONOutput(cmd *cobra.Command, filePath string, lineCount int, commentCount int, classes int, functions int, indentResults map[string]interface{}, dependencieResults map[string]interface{}) {
    // Limpeza de dados desnecessários em indentResults
    if indentResults != nil {
        delete(indentResults, "path")
        delete(indentResults, "filename")
    }

    // Estrutura do JSON reorganizada
    result := map[string]interface{}{
        "directory": filePath,
        "summary": map[string]interface{}{
            "lines":      lineCount,
            "comments":   commentCount,
            "classes":    classes,
            "functions":  functions,
            "dependencies": consolidateDependencies(dependencieResults),
            "indentation": indentResults,
        },
    }

    // Verifica se o caminho de saída foi especificado
    if utils.OutputFilePath != "" {
        outputPath := utils.OutputFilePath

        // Ajusta o caminho se for um diretório
        fileInfo, err := os.Stat(outputPath)
        if err == nil && fileInfo.IsDir() {
            outputPath = fmt.Sprintf("%s/analysis_report.json", outputPath)
        }

        // Cria o arquivo JSON
        file, err := os.Create(outputPath)
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error creating JSON file: %s\n", err)
            return
        }
        defer file.Close()

        // Escreve o JSON no arquivo
        encoder := json.NewEncoder(file)
        encoder.SetIndent("", "  ")
        if err := encoder.Encode(result); err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error writing JSON to file: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "JSON report saved to %s\n", outputPath)
    } else {
        // Formata e exibe o JSON no terminal
        formattedJSON, err := json.MarshalIndent(result, "", "  ")
        if err != nil {
            fmt.Fprintf(cmd.OutOrStderr(), "Error formatting JSON: %s\n", err)
            return
        }

        fmt.Fprintf(cmd.OutOrStdout(), "%s\n", string(formattedJSON))
    }
}

func generateDetailedJSONOutput(cmd *cobra.Command, directoryPath string, lineResults analyzer.FilesNameCountLineMap, commentResults analyzer.CommentsMap, classFuncResults analyzer.ClassesAndFunctionsMap, indentResults map[string]interface{}, dependenciesResults map[string]interface{}) {
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
        "dependencies": consolidateDependencies(dependenciesResults),
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

func printFileResults(cmd *cobra.Command, filePath string, lines, comments, classes, functions int, indentResults map[string]interface{}, dependenciesResults map[string]interface{}) {
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

    hasDeps := dependenciesResults["total_dependencies"] != nil && dependenciesResults["dependencies"] != nil && dependenciesResults["native_modules"] != nil

    if hasDeps  {
        fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Dependencies Analysis ===%s\n", utils.BLUE, utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Total Dependencies: %s%d%s\n", utils.GREEN, dependenciesResults["total_dependencies"], utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Dependencies: %s%v%s\n", utils.GREEN, dependenciesResults["dependencies"], utils.RESET_COLOR)
        fmt.Fprintf(cmd.OutOrStdout(), "Native Modules: %s%v%s\n", utils.GREEN, dependenciesResults["native_modules"], utils.RESET_COLOR)
    }
}

func printDirectoryResults(cmd *cobra.Command, lines, comments, classes, functions int, indentResults map[string]interface{}, ) {
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

func init() {
    RunAllCommand.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    RunAllCommand.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    RunAllCommand.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Specify the output file path. If omitted, results will be displayed in the terminal.")
    RunAllCommand.Flags().BoolVar(&utils.Detailed, "detailed", false, "Show detailed analysis, including per-file information.")
}