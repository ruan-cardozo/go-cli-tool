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
    Classes    int
    Func       int
    Comments   int
}

type ReportData struct {
    Files        []FileResult
    TotalByDirectory  string
    CommandType  utils.CommandType
    HasClasses   bool
    HasFunctions bool
}

type GenericsType interface {
   analyzer.FilesNameCountLineMap | analyzer.ClassesAndFunctionsMap | analyzer.CommentsMap
}

func SaveResultsToHTML[T GenericsType](
    result T,
    totalByDirectory string,
    filePath string,
    commandType utils.CommandType,
    cmd *cobra.Command,
    hasClasses bool,
    hasFunctions bool) {

    var files []FileResult
    switch v := any(result).(type) {
    case analyzer.FilesNameCountLineMap:
        for fileName, res := range v {
            files = append(files, FileResult{
                FileName:   fileName,
                TotalLines: res.TotalLines,
            })
        }
    case analyzer.ClassesAndFunctionsMap:
        for fileName, res := range v {
            files = append(files, FileResult{
                FileName:   fileName,
                Classes:    res.Classes,
                Func:       res.Functions,
                TotalLines: res.Classes + res.Functions,
            })
        }
    case analyzer.CommentsMap:
        for fileName, res := range v {
            files = append(files, FileResult{
                FileName:   fileName,
                TotalLines:   res.CommentLines,
            })
        }
    }

    data := ReportData{
        Files:            files,
        TotalByDirectory: totalByDirectory,
        CommandType:      commandType,
        HasClasses:       hasClasses,
        HasFunctions:     hasFunctions,
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

    fmt.Fprintf(cmd.OutOrStdout(), "\033[1;34mReport generated successfully at %s\033[0m\n", filePath)
}

// func SaveResultsToHTMLForClassesAndFunc(
//     result analyzer.ClassesAndFunctionsMap, 
//     totalLinesByDirectory analyzer.ClassFuncResult, 
//     filePath string,
//     commandType utils.CommandType,
//     cmd *cobra.Command) {

//     var files []FileResultClassesFunc
//     for fileName, res := range result {
//         files = append(files, FileResultClassesFunc{FileName: fileName, Classes: res.Classes, Func: res.Functions})
//     }

//     data := ReportDataClassesFunc{
//         Files:      files,
//         Classes: totalLinesByDirectory.Classes,
//         Func: totalLinesByDirectory.Functions,
//         CommandType: commandType,
//     }

//     tmpl, err := template.New("report").Parse(reportTemplate)
//     if err != nil {
//         fmt.Printf("Error parsing template: %s\n", err)
//         return
//     }

// 	// Ensure the directory exists
// 	outputDir := filepath.Dir(filePath)
// 	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
// 		err = os.MkdirAll(outputDir, os.ModePerm)
// 		if err != nil {
// 			fmt.Printf("Error creating directory: %s\n", err)
// 			return
// 		}
// 	}

// 	filePath = filepath.Join(outputDir, "report.html")
//     file, err := os.Create(filePath)
//     if err != nil {
//         fmt.Printf("Error creating file: %s\n", err)
//         return
//     }
//     defer file.Close()

//     err = tmpl.Execute(file, data)
//     if err != nil {
//         fmt.Printf("Error executing template: %s\n", err)
//     }

//     fmt.Fprintf(cmd.OutOrStdout(),"\033[1;34mReport generated successfully at %s\033[0m\n", filePath)
// }