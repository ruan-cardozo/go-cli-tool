package run_all_commands

import (
	"encoding/json"
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/utils"
	"os"

	"github.com/spf13/cobra"
)

type AnalysisParams struct {
	FilePath              string
	DirectoryPath         string
	OutputFilePath        string
	Detailed              bool
	LineCount             int64
	CommentCount          int
	Classes               int
	Functions             int
	CommentPercentage     float64
	AverageFunctionSize   float64
	OverallAverageSize    float64
	IndentResults         map[string]interface{}
	DependenciesResults   map[string]interface{}
	MethodCountResult     analyzer.MethodCountResult
	TotalMethodCount      analyzer.MethodCountResult
	LineResults           analyzer.FilesNameCountLineMap
	CommentResults        analyzer.CommentsMap
	ClassFuncResults      analyzer.ClassesAndFunctionsMap
	PercentResults        analyzer.PercentResult
	MethodCountResults    analyzer.MethodCountMap
}

var RunAllCommand = &cobra.Command{
	Use:   "analyze",
	Short: "Comprehensive JavaScript code analysis in a single operation",
	Long: `Execute a complete analysis of JavaScript code with a single command.
    
This command performs multiple analyses simultaneously including:
- Line count analysis
- Comment count analysis
- Function and class count analysis
- Indentation analysis
- Code Comment Percentage Analysis
- Method count analysis (public/private)
- Average Function Size Analysis
- Dependency analysis

Results are presented in terminal or json output, providing a complete overview
of your JavaScript codebase. Use flags to customize the analysis and output format.`,
	Run: func(cmd *cobra.Command, args []string) {
		if utils.FilePath == "" && utils.DirectoryPath == "" {
			fmt.Fprintf(cmd.OutOrStderr(), "%sError: You must provide either a file path (-f) or directory path (-d).%s\n",
				utils.RED, utils.RESET_COLOR)
			return
		}

		lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer, methodCountAnalyzer, averageFunctionAnalyzer := initializeAnalyzers()
		percentAnalyzer := &analyzer.CountPercentAnalyzerImpl{}

		if utils.FilePath != "" {
			handleFileAnalysis(cmd, lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer, percentAnalyzer, methodCountAnalyzer, averageFunctionAnalyzer)
		} else {
			handleDirectoryAnalysis(cmd, lineAnalyzer, commentAnalyzer, classFuncAnalyzer, indentationAnalyzer, dependenciesAnalyzer, percentAnalyzer, methodCountAnalyzer, averageFunctionAnalyzer)
		}
	},
}

func initializeAnalyzers() (*analyzer.CountLinesAnalyzerImpl, *analyzer.CountCommentsAnalyzerImpl, *analyzer.CountClassAndFunctionsImpl, *analyzer.IdentationAnalyzerImpl, *analyzer.CountDependenciesAnalyzerImpl, *analyzer.MethodCountAnalyzerImpl, *analyzer.AverageFunctionAnalyzerImpl) {
	return &analyzer.CountLinesAnalyzerImpl{},
		&analyzer.CountCommentsAnalyzerImpl{},
		&analyzer.CountClassAndFunctionsImpl{},
		&analyzer.IdentationAnalyzerImpl{},
		&analyzer.CountDependenciesAnalyzerImpl{},
		&analyzer.MethodCountAnalyzerImpl{},
		&analyzer.AverageFunctionAnalyzerImpl{}
}

func handleFileAnalysis(cmd *cobra.Command, lineAnalyzer *analyzer.CountLinesAnalyzerImpl, commentAnalyzer *analyzer.CountCommentsAnalyzerImpl, classFuncAnalyzer *analyzer.CountClassAndFunctionsImpl, indentationAnalyzer *analyzer.IdentationAnalyzerImpl, dependenciesAnalyzer *analyzer.CountDependenciesAnalyzerImpl, percentAnalyzer *analyzer.CountPercentAnalyzerImpl, methodCountAnalyzer *analyzer.MethodCountAnalyzerImpl, averageFunctionAnalyzer *analyzer.AverageFunctionAnalyzerImpl) {
	lineCount := lineAnalyzer.CountLinesByFilePath(utils.FilePath)
	commentCount := commentAnalyzer.CountCommentsByFilePath(utils.FilePath)
	classAndFunctionResult := classFuncAnalyzer.CountClassesAndFunctionsByFilePath(utils.FilePath)
	percentResult := percentAnalyzer.CountPercentByFilePath(utils.FilePath)
	methodCountResult := methodCountAnalyzer.AnalyzeFile(utils.FilePath)
	averageFunctionSize := averageFunctionAnalyzer.CalculateAverageFunctionSize(utils.FilePath)

	tempFilePath := utils.FilePath
	tempDirPath := utils.DirectoryPath
	utils.DirectoryPath = ""
	indentResults, _ := indentationAnalyzer.IdentationByFilePath()
	utils.FilePath = tempFilePath
	utils.DirectoryPath = tempDirPath

	dependencieResultMap, _ := dependenciesAnalyzer.CountDependenciesByFilePath(utils.FilePath)

	params := AnalysisParams{
		FilePath:            utils.FilePath,
		DirectoryPath:       utils.DirectoryPath,
		OutputFilePath:      utils.OutputFilePath,
		Detailed:            utils.Detailed,
		LineCount:           int64(lineCount.TotalLines),
		CommentCount:        commentCount.CommentLines,
		Classes:             classAndFunctionResult.Classes,
		Functions:           classAndFunctionResult.Functions,
		CommentPercentage:   percentResult.CommentPercentage,
		IndentResults:       indentResults,
		DependenciesResults: dependencieResultMap,
		MethodCountResult:   methodCountResult,
		AverageFunctionSize: averageFunctionSize,
	}

	if utils.OutputFilePath == "" {
		printFileResults(cmd, params)
	} else {
		generateJSONOutput(cmd, params)
	}
}

func handleDirectoryAnalysis(cmd *cobra.Command, lineAnalyzer *analyzer.CountLinesAnalyzerImpl, commentAnalyzer *analyzer.CountCommentsAnalyzerImpl, classFuncAnalyzer *analyzer.CountClassAndFunctionsImpl, indentationAnalyzer *analyzer.IdentationAnalyzerImpl, dependenciesAnalyzer *analyzer.CountDependenciesAnalyzerImpl, percentAnalyzer *analyzer.CountPercentAnalyzerImpl, methodCountAnalyzer *analyzer.MethodCountAnalyzerImpl, averageFunctionAnalyzer *analyzer.AverageFunctionAnalyzerImpl) {
	lineResults, totalLines := lineAnalyzer.CountLinesByDirectory(utils.DirectoryPath)
	commentResults, totalComments := commentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
	classFuncResults, totalClassesAndFunctions := classFuncAnalyzer.CountClassesAndFunctionsByDirectory(utils.DirectoryPath)
	_, percentResults := percentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)
	methodCountResults, totalMethodCount := methodCountAnalyzer.AnalyzeDirectory(utils.DirectoryPath)
	_, overallAverage := averageFunctionAnalyzer.CalculateAverageFunctionSizeByDirectory(utils.DirectoryPath)

	tempFilePath := utils.FilePath
	tempDirPath := utils.DirectoryPath
	utils.FilePath = ""
	indentResults, _ := indentationAnalyzer.IdentationByFilePath()
	utils.FilePath = tempFilePath
	utils.DirectoryPath = tempDirPath

	dependencieResultMap, _ := dependenciesAnalyzer.CountDependenciesByDirectory(utils.DirectoryPath)

	params := AnalysisParams{
		DirectoryPath:       utils.DirectoryPath,
		OutputFilePath:      utils.OutputFilePath,
		Detailed:            utils.Detailed,
		LineCount:           int64(totalLines.TotalLines),
		CommentCount:        totalComments.TotalComments,
		Classes:             totalClassesAndFunctions.Classes,
		Functions:           totalClassesAndFunctions.Functions,
		CommentPercentage:   percentResults.CommentPercentage,
		IndentResults:       indentResults,
		DependenciesResults: dependencieResultMap,
		TotalMethodCount:    totalMethodCount,
		OverallAverageSize:  overallAverage,
		LineResults:         lineResults,
		CommentResults:      commentResults,
		ClassFuncResults:    classFuncResults,
		PercentResults:      percentResults,
		MethodCountResults:  methodCountResults,
	}

	if utils.Detailed {
		generateDetailedJSONOutput(cmd, params)
	} else if utils.OutputFilePath == "" {
		printDirectoryResults(cmd, params)
	} else {
		generateJSONOutput(cmd, params)
	}
}

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

func generateJSONOutput(cmd *cobra.Command, params AnalysisParams) {
	if params.IndentResults != nil {
		delete(params.IndentResults, "path")
		delete(params.IndentResults, "filename")
	}

	summaryData := map[string]interface{}{
		"lines":                 params.LineCount,
		"comments":              params.CommentCount,
		"comment_percentage":    fmt.Sprintf("%.2f%%", params.CommentPercentage),
		"classes":               params.Classes,
		"functions":             params.Functions,
		"public_methods":        params.MethodCountResult.Public,
		"private_methods":       params.MethodCountResult.Private,
		"average_function_size": fmt.Sprintf("%.4f", params.AverageFunctionSize),
		"dependencies":          consolidateDependencies(params.DependenciesResults),
		"indentation":           params.IndentResults,
	}

	result := map[string]interface{}{
		"directory": params.FilePath,
		"summary":   summaryData,
	}

	if params.DirectoryPath != "" && params.FilePath == "" {
		result["directory"] = params.DirectoryPath
		summaryData["public_methods"] = params.TotalMethodCount.Public
		summaryData["private_methods"] = params.TotalMethodCount.Private
		summaryData["average_function_size"] = fmt.Sprintf("%.4f", params.OverallAverageSize)
	}

	outputJSON(cmd, params.OutputFilePath, result)
}

func generateDetailedJSONOutput(cmd *cobra.Command, params AnalysisParams) {
	fileIndentData := make(map[string]interface{})
	if files, ok := params.IndentResults["files"].([]map[string]interface{}); ok {
		for _, fileData := range files {
			if filename, ok := fileData["filename"].(string); ok {
				fileIndentData[filename] = fileData["stats"]
			}
		}
	}

	fileDetails := make([]map[string]interface{}, 0, len(params.LineResults))
	for filename, lineResult := range params.LineResults {
		fileInfo := map[string]interface{}{
			"filename": filename,
			"metrics": map[string]interface{}{
				"lines":           lineResult.TotalLines,
				"comments":        params.CommentResults[filename].CommentLines,
				"classes":         params.ClassFuncResults[filename].Classes,
				"functions":       params.ClassFuncResults[filename].Functions,
				"public_methods":  params.MethodCountResults[filename].Public,
				"private_methods": params.MethodCountResults[filename].Private,
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

	for filePath := range params.LineResults {
		totalLines += int(params.LineResults[filePath].TotalLines)
		totalComments += params.CommentResults[filePath].CommentLines
		totalClasses += params.ClassFuncResults[filePath].Classes
		totalFunctions += params.ClassFuncResults[filePath].Functions
	}

	detailedResult := map[string]interface{}{
		"directory_path": params.DirectoryPath,
		"summary": map[string]interface{}{
			"total_files":       len(params.LineResults),
			"total_lines":       totalLines,
			"total_comments":    totalComments,
			"total_classes":     totalClasses,
			"total_functions":   totalFunctions,
			"total_public_methods":  params.TotalMethodCount.Public,
			"total_private_methods": params.TotalMethodCount.Private,
		},
		"dependencies": consolidateDependencies(params.DependenciesResults),
		"files":        fileDetails,
	}

	if params.OutputFilePath != "" {
		outputPath := params.OutputFilePath

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

func printFileResults(cmd *cobra.Command, params AnalysisParams) {
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Analysis Results for %s ===%s\n",
		utils.BLUE, params.FilePath, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, params.LineCount, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, params.CommentCount, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Comment Percentage: %s%.2f%%%s\n", utils.GREEN, params.CommentPercentage, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, params.Classes, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, params.Functions, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Public Methods: %s%d%s\n", utils.GREEN, params.MethodCountResult.Public, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Private Methods: %s%d%s\n", utils.GREEN, params.MethodCountResult.Private, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Average Function Size: %s%.2f%s\n", utils.GREEN, params.AverageFunctionSize, utils.RESET_COLOR)

	if stats, ok := params.IndentResults["stats"].(analyzer.IndentResult); ok {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Indentation Analysis ===%s\n", utils.BLUE, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Max Indent Level: %s%d%s\n", utils.GREEN, stats.MaxIndentLevel, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Average Indent Level: %s%.2f%s\n", utils.GREEN, stats.AverageIndentLevel, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Uses Spaces: %s%t%s\n", utils.GREEN, stats.UsesSpaces, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Uses Tabs: %s%t%s\n", utils.GREEN, stats.UsesTabs, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Mixed Indentation: %s%t%s\n", utils.GREEN, stats.MixedIndentation, utils.RESET_COLOR)
	}

	hasDeps := params.DependenciesResults["total_dependencies"] != nil && params.DependenciesResults["dependencies"] != nil && params.DependenciesResults["native_modules"] != nil

	if hasDeps {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Dependencies Analysis ===%s\n", utils.BLUE, utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Total Dependencies: %s%d%s\n", utils.GREEN, params.DependenciesResults["total_dependencies"], utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Dependencies: %s%v%s\n", utils.GREEN, params.DependenciesResults["dependencies"], utils.RESET_COLOR)
		fmt.Fprintf(cmd.OutOrStdout(), "Native Modules: %s%v%s\n", utils.GREEN, params.DependenciesResults["native_modules"], utils.RESET_COLOR)
	}
}

func printDirectoryResults(cmd *cobra.Command, params AnalysisParams) {
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Directory Analysis Summary ===%s\n",
		utils.BLUE, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Lines: %s%d%s\n", utils.GREEN, params.LineCount, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Comment Lines: %s%d%s\n", utils.GREEN, params.CommentCount, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Comment Percentage: %s%.2f%%%s\n", utils.GREEN, params.CommentPercentage, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Classes: %s%d%s\n", utils.GREEN, params.Classes, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Functions: %s%d%s\n", utils.GREEN, params.Functions, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Public Methods: %s%d%s\n", utils.GREEN, params.TotalMethodCount.Public, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Total Private Methods: %s%d%s\n", utils.GREEN, params.TotalMethodCount.Private, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "Average function size in directory: %s%.2f lines%s\n", utils.GREEN, params.OverallAverageSize, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "\n%s=== Indentation Analysis Summary ===%s\n", utils.BLUE, utils.RESET_COLOR)

	if files, ok := params.IndentResults["files"].([]map[string]interface{}); ok && len(files) > 0 {

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
}