package templates

import (
	_ "embed"
	"fmt"
	"go-cli-tool/internal/analyzer"
	"go-cli-tool/internal/utils"
	"html/template"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed report.html
var reportTemplate string

type FileResult struct {
    FileName   string
    TotalLines int
}

type ReportData struct {
    Files      []FileResult
    TotalLines int
    CommandType utils.CommandType
}

func SaveResultsToHTML(
    result analyzer.FilesNameCountLineMap, 
    totalLinesByDirectory analyzer.LineResult, 
    filePath string,
    commandType utils.CommandType,
    cmd *cobra.Command) {

    var files []FileResult
    for fileName, res := range result {
        files = append(files, FileResult{FileName: fileName, TotalLines: res.TotalLines})
    }

    data := ReportData{
        Files:      files,
        TotalLines: totalLinesByDirectory.TotalLines,
        CommandType: commandType,
    }

    tmpl, err := template.New("report").Parse(reportTemplate)
    if err != nil {
        fmt.Printf("Error parsing template: %s\n", err)
        return
    }

	// Ensure the directory exists
	outputDir := filepath.Dir(filePath)
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory: %s\n", err)
			return
		}
	}

	filePath = filepath.Join(outputDir, "report.html")
    file, err := os.Create(filePath)
    if err != nil {
        fmt.Printf("Error creating file: %s\n", err)
        return
    }
    defer file.Close()

    err = tmpl.Execute(file, data)
    if err != nil {
        fmt.Printf("Error executing template: %s\n", err)
    }

    fmt.Fprintf(cmd.OutOrStdout(),"\033[1;34mReport generated successfully at %s\033[0m\n", filePath)
}