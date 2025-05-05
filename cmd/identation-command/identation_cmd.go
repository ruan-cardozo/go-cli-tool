package identation

import (
	"encoding/json"
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/utils"
	"os"

	"github.com/spf13/cobra"
)

var IdentationAnalyzerCmd = &cobra.Command{
    Use:   "identation",
    Short: "Check identation in JavaScript files",
    Run: func(cmd *cobra.Command, args []string) {
        if utils.FilePath == "" && utils.DirectoryPath == "" {
            fmt.Println("Error: You must specify either a file path or a directory path")
            cmd.Help()
            return
        }
        
        identationAnalyzer := &analyzer.IdentationAnalyzerImpl{}
        results, err := identationAnalyzer.IdentationByFilePath()
        if err != nil {
            fmt.Printf("Error analyzing indentation: %v\n", err)
            return
        }
        
        if utils.OutputFilePath != "" {

            file, err := os.Create(utils.OutputFilePath)
            if err != nil {
                fmt.Printf("Error creating output file: %v\n", err)
                return
            }
            defer file.Close()

            encoder := json.NewEncoder(file)
            encoder.SetIndent("", "  ")
            if err := encoder.Encode(results); err != nil {
                fmt.Printf("Error writing JSON to file: %v\n", err)
                return
            }

            fmt.Printf("Results written to %s\n", utils.OutputFilePath)
        } else {

            jsonData, err := json.MarshalIndent(results, "", "  ")
            if err != nil {
                fmt.Printf("Error formatting results: %v\n", err)
                return
            }
            
            fmt.Println(string(jsonData))
        }
    },
}

func init() {
    IdentationAnalyzerCmd.Flags().StringVarP(&utils.FilePath, "file", "f", "", "Path to the JavaScript file (must be a single file, not a directory)")
    IdentationAnalyzerCmd.Flags().StringVarP(&utils.DirectoryPath, "directory", "d", "", "Path to the directory containing JavaScript files. The tool will automatically expand the provided path.")
    IdentationAnalyzerCmd.Flags().StringVarP(&utils.OutputFilePath, "output", "o", "", "Path to the output file. The expected output is a json file.")
}