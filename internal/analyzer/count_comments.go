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

type CountCommentsAnalyzerImpl struct{}

func (a *CountCommentsAnalyzerImpl) CountCommentsByFilePath(filePath string) CommentResult {
    file, err := os.Open(filePath)

    if err != nil {
        panic(err)
    }

    defer file.Close()

    var result CommentResult
    scanner := bufio.NewScanner(file)
    inBlockComment := false
    for scanner.Scan() {
        line := scanner.Text()

        if isComment(line,&inBlockComment) {
            result.CommentLines++
        }
    }

    return result
}

func (a *CountCommentsAnalyzerImpl) CountCommentsByDirectory(directoryPath string) (CommentsMap, CommentResult) {
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

    linesByArchive := make(CommentsMap)

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
            linesByArchive[fileOrDirectoryName] = a.CountCommentsByFilePath(path)
        }

        return nil
    })

    if err != nil {
        panic(err)
    }

    var totalCommentsByDirectory CommentResult
 
    for result := range linesByArchive {
        file := linesByArchive[result]
        totalCommentsByDirectory.TotalComments += file.CommentLines
    }

    return linesByArchive, totalCommentsByDirectory
}

func isComment(line string, inBlockComment *bool) bool {
	singleLineComment := regexp.MustCompile(`^\s*//`)
	blockCommentStart := regexp.MustCompile(`^\s*/\*`)
	blockCommentEnd := regexp.MustCompile(`\*/`)

	if *inBlockComment {
		if blockCommentEnd.MatchString(line) {
			*inBlockComment = false
		}
		return false
	}

	if singleLineComment.MatchString(line) {
		return true
	}

	if blockCommentStart.MatchString(line) {
		*inBlockComment = true
		return true
	}

	return false
}