package count_percent

import (
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/templates"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var countPercentAnalyzer analyzer.CountPercentAnalyzer

var CountPercentCmd = &cobra.Command{
	Use:   "count-percent",
	Short: "Count total comment lines and calculate the percentage of comments in a JavaScript file",
	Run: func(cmd *cobra.Command, args []string) {
		err := policies.ValidateUserInput(cmd)
		if err {
			return
		}

		// Se o arquivo for passado com o flag -f
		if utils.FilePath != "" {
			result := countPercentAnalyzer.CountPercentByFilePath(utils.FilePath)
			fmt.Fprintf(cmd.OutOrStdout(), "%sTotal comments in file:%s %d\n", utils.BLUE, utils.RESET_COLOR, result.CommentLines)
			fmt.Fprintf(cmd.OutOrStdout(), "%sComment Percentage in file:%s %.2f%%\n", utils.BLUE, utils.RESET_COLOR, result.CommentPercentage)
			return
		}

		// Se o diretório for passado com o flag -d
		if utils.DirectoryPath != "" {
			if !policies.ValidateDirectoryPath(utils.DirectoryPath) {
				fmt.Fprintf(cmd.OutOrStdout(), "%sPlease provide a valid directory path.%s\n", utils.RED, utils.RESET_COLOR)
				return
			}

			results, totals := countPercentAnalyzer.CountCommentsByDirectory(utils.DirectoryPath)

			if len(results) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "%sNo JavaScript files found in the provided directory.%s\n", utils.RED, utils.RESET_COLOR)
				return
			}

			// Exibir total de comentários no diretório
			totalLinesStr := fmt.Sprintf("Total lines: %d | Comment lines: %d", totals.TotalLines, totals.CommentLines)

			// Se o usuário desejar gerar o relatório em HTML
			if utils.OutputFilePath != "" {
				templates.SaveResultsToHTML(results, totalLinesStr, utils.OutputFilePath, utils.COUNT_COMMENTS, cmd, false, false)
			} else {
				printResults(results, totals, cmd)
			}
		}
	},
}

// Função para exibir os resultados
func printResults(results analyzer.PercentResultMap, totals analyzer.PercentResult, cmd *cobra.Command) {
	// Exibe os resultados dos arquivos individuais
	for filePath, result := range results {
		shortName := filepath.Base(filePath)
		fmt.Fprintf(cmd.OutOrStdout(), "%sComment lines in %s:%s %d (%.2f%%)\n",
			utils.BLUE, shortName, utils.RESET_COLOR,
			result.CommentLines, result.CommentPercentage)
	}

	// Exibe os totais consolidados
	fmt.Fprintf(cmd.OutOrStdout(), "%s\n%sTOTAL RESULTS%s\n",
		strings.Repeat("-", 40), utils.GREEN, utils.RESET_COLOR)
	fmt.Fprintf(cmd.OutOrStdout(), "%sTotal Lines:%s %d\n",
		utils.BLUE, utils.RESET_COLOR, totals.TotalLines)
	fmt.Fprintf(cmd.OutOrStdout(), "%sTotal Comments:%s %d\n",
		utils.BLUE, utils.RESET_COLOR, totals.CommentLines)
	fmt.Fprintf(cmd.OutOrStdout(), "%sComment Percentage:%s %.2f%%\n",
		utils.BLUE, utils.RESET_COLOR, totals.CommentPercentage)
}

func init() {
	// Inicializa a implementação correta
	countPercentAnalyzer = &analyzer.CountPercentAnalyzerImpl{}

	CountPercentCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
	CountPercentCmd.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files")
	CountPercentCmd.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output HTML report file (only for directory analysis)")
}
