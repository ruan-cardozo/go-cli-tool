package analyzer

import (
	"bufio"

	// "fmt"
	"go-cli-tool/internal/policies"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"slices"
)
type FilesNameCountLineMap map[string]LineResult
var directoryOrFilesToIgnore = []string{".git", "node_modules"}
var totalLinesByDirectory LineResult

func CountLinesByFilePath(filePath string) LineResult {
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

func CountLinesByDirectory(directoryPath string) (FilesNameCountLineMap, LineResult) {

    if directoryPath == "." {
        var err error
        directoryPath, err = os.Getwd()
        if err != nil {
            panic(err)
        }
    } else {
        var err error
        directoryPath, err = expandPath(directoryPath)
        if err != nil {
            panic(err)
        }
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

            linesByArchive[fileOrDirectoryName] = CountLinesByFilePath(path)
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

func expandPath(path string) (string, error) {
    if path[:2] == "~/" {
        usr, err := user.Current()
        if err != nil {
            return "", err
        }
        path = filepath.Join(usr.HomeDir, path[2:])
    }
    return path, nil
}