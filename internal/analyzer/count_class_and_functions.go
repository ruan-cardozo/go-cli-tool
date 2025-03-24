package analyzer

import (
	"bufio"
	"fmt"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

type CountClassAndFunctionsImpl struct {}

func (a *CountClassAndFunctionsImpl) CountClassesAndFunctionsByFilePath(filePath string) ClassFuncResult {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var result ClassFuncResult
    scanner := bufio.NewScanner(file)

    functionRegex := regexp.MustCompile(`(?m)(function\s+\w+\s*\([^)]*\)|\w+\s*=\s*function\s*\([^)]*\)|\w+\s*=\s*\([^)]*\)\s*=>|\w+\s*\([^)]*\)\s*{)`)
    classRegex := regexp.MustCompile(`(?m)class\s+\w+`)

    for scanner.Scan() {
        line := scanner.Text()

        if functionRegex.MatchString(line) {
            result.Functions++
        }

        if classRegex.MatchString(line) {
            result.Classes++
        }
    }

    return result
}
func (a *CountClassAndFunctionsImpl) CountClassesAndFunctionsByDirectory(directoryPath string) (ClassesAndFunctionsMap, ClassFuncResult) {
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

    if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
        panic(fmt.Sprintf("directory %s does not exist", directoryPath))
    }

    linesByArchive := make(ClassesAndFunctionsMap)

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
            linesByArchive[fileOrDirectoryName] = a.CountClassesAndFunctionsByFilePath(path)
        }

        return nil
    })

    if err != nil {
        panic(err)
    }

    var totalClassesAndFunctions ClassFuncResult
    for _, result := range linesByArchive {
        totalClassesAndFunctions.Classes += result.Classes
        totalClassesAndFunctions.Functions += result.Functions
    }

    return linesByArchive, totalClassesAndFunctions
}