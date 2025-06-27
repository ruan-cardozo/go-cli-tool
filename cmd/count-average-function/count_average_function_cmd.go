package count_average_function_size

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"

	"github.com/spf13/cobra"
)

var averageFunctionAnalyzer analyzer.AverageFunctionAnalyzer = &analyzer.AverageFunctionAnalyzerImpl{}

var CountAverageFunctionSizeCmd = &cobra.Command{
	Use:   "count-average-function-size",
	Short: "Calculate the average function size in a JavaScript file or directory",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")
		directoryPath, _ := cmd.Flags().GetString("directory")

		if filePath == "" && directoryPath == "" {
			fmt.Fprintf(cmd.OutOrStdout(), "%sInvalid input. Please provide a file (-f) or directory (-d).%s\n", utils.RED, utils.RESET_COLOR)
			return
		}

		utils.FilePath = filePath
		utils.DirectoryPath = directoryPath

		if utils.FilePath != "" {
			average := averageFunctionAnalyzer.CalculateAverageFunctionSize(utils.FilePath)
			fmt.Fprintf(cmd.OutOrStdout(), "%sAverage function size:%s %.2f lines\n", utils.BLUE, utils.RESET_COLOR, average)
			return
		}

		if utils.DirectoryPath != "" {
			if !policies.ValidateDirectoryPath(utils.DirectoryPath) {
				fmt.Fprintf(cmd.OutOrStdout(), "%sPlease provide a valid directory path.%s\n", utils.RED, utils.RESET_COLOR)
				return
			}

			results, overallAverage := averageFunctionAnalyzer.CalculateAverageFunctionSizeByDirectory(utils.DirectoryPath)

			if len(results) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%sNo JavaScript files found in the provided directory.%s\n", utils.RED, utils.RESET_COLOR)
				return
			}

			for file, avg := range results {
				fmt.Fprintf(cmd.OutOrStdout(), "%sAverage function size in %s:%s %.2f lines\n", utils.BLUE, file, utils.RESET_COLOR, avg)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%sOverall average function size in directory:%s %.2f lines\n", utils.BLUE, utils.RESET_COLOR, overallAverage)
		}
	},
}

func init() {
	CountAverageFunctionSizeCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file")
	CountAverageFunctionSizeCmd.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files")
}
