package count_methods

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"

	"github.com/spf13/cobra"
)

var methodCountAnalyzer analyzer.MethodCountAnalyzer

var CountMethodsAnalyzer = &cobra.Command{
	Use:   "count-methods",
	Short: "Count public and private methods in JavaScript files",
	Run: func(cmd *cobra.Command, args []string) {
		err := policies.ValidateUserInput(cmd)
		if err {
			return
		}

		if utils.FilePath != "" {
			result := methodCountAnalyzer.AnalyzeFile(utils.FilePath)
			printSingleFileResult(result, cmd)
			return
		}

		if utils.DirectoryPath != "" && !policies.ValidateDirectoryPath(utils.DirectoryPath) {
			fmt.Fprintf(cmd.OutOrStdout(), "%sInvalid directory path.%s\n", utils.RED, utils.RESET_COLOR)
			return
		}

		result, total := methodCountAnalyzer.AnalyzeDirectory(utils.DirectoryPath)

		if len(result) == 0 {
			fmt.Fprintf(cmd.OutOrStdout(), "%sNo JavaScript files found in the provided directory.%s\n", utils.RED, utils.RESET_COLOR)
			return
		}

		if utils.OutputFilePath != "" {
			// Converte o resultado para o formato usado no template HTML
			filesNameCountLineMap := make(analyzer.FilesNameCountLineMap)
			for file, data := range result {
				filesNameCountLineMap[file] = analyzer.LineResult{
					TotalLines: int(data.Public + data.Private),
				}
			}

			summary := fmt.Sprintf("Total public: %d, private: %d", total.Public, total.Private)

			templates.SaveResultsToHTML(filesNameCountLineMap, summary, utils.OutputFilePath, "COUNT_METHODS", cmd, true, true)
		} else {
			printDirectoryResults(result, total, cmd)
		}
	},
}

func printSingleFileResult(result analyzer.MethodCountResult, cmd *cobra.Command) {
	fmt.Fprintf(cmd.OutOrStdout(), "%sPublic methods:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.Public)
	fmt.Fprintf(cmd.OutOrStdout(), "%sPrivate methods:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.Private)
}

func printDirectoryResults(result analyzer.MethodCountMap, total analyzer.MethodCountResult, cmd *cobra.Command) {
	for file, data := range result {
		fmt.Fprintf(cmd.OutOrStdout(), "%s%s:%s public=%d, private=%d\n", utils.BLUE, file, utils.RESET_COLOR, data.Public, data.Private)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "%sTotal:%s public=%d, private=%d\n", utils.BLUE, utils.RESET_COLOR, total.Public, total.Private)
}

func init() {
	methodCountAnalyzer = &analyzer.MethodCountAnalyzerImpl{}
	CountMethodsAnalyzer.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to a JavaScript file")
	CountMethodsAnalyzer.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to a directory")
	CountMethodsAnalyzer.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "HTML output file (optional)")
}
