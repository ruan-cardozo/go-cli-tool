package analyzer

import (
	"bufio"
	"fmt"

	// "fmt"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)
type FilesNameCountLineMap map[string]LineResult
var directoryOrFilesToIgnore = []string{".git", "node_modules"}
var totalLinesByDirectory LineResult

type CountLinesAnalyzerImpl struct{}

func (a *CountLinesAnalyzerImpl) CountLinesByFilePath(filePath string) LineResult {
    file, err := os.Open(filePath)

    if err != nil {
        panic(err)
    }

    defer file.Close()

    var result LineResult
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

        line := scanner.Text()

        if isEmptyLine(line) {
            continue
        }

        result.TotalLines++
    }

    return result
}

func (a *CountLinesAnalyzerImpl) CountLinesByDirectory(directoryPath string) (FilesNameCountLineMap, LineResult) {
    if directoryPath == "." {
        var err error
        directoryPath, err = os.Getwd()
        if err != nil {
            panic(err)
        }
    } else {
        var err error
        directoryPath, err = utils.ExpandPath(directoryPath)
        if err != nil {
            panic(err)
        }
    }

    // Verificar se o diretório existe
    if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
        panic(fmt.Sprintf("directory %s does not exist", directoryPath))
    }

    linesByArchive := make(FilesNameCountLineMap)

    err := filepath.WalkDir(directoryPath, func(path string, directory fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        fileOrDirectoryName := directory.Name()
        fileExtension := filepath.Ext(fileOrDirectoryName)

        if slices.Contains(directoryOrFilesToIgnore, fileOrDirectoryName) {
            if directory.IsDir() {
                return filepath.SkipDir
            }
            return nil
        }

        if policies.IsJSFileExtension(fileExtension) {
            linesByArchive[fileOrDirectoryName] = a.CountLinesByFilePath(path)
        }

        return nil
    })

    if err != nil {
        panic(err)
    }

    for result := range linesByArchive {
        file := linesByArchive[result]
        totalLinesByDirectory.TotalLines += file.TotalLines
    }

    return linesByArchive, totalLinesByDirectory
}

func isEmptyLine(line string) bool {
    return len(line) == 0
}