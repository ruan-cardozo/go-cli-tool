package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var ignoreList = []string{"node_modules", ".git", "dist", "build", ".DS_Store"}

type CountPercentAnalyzer interface {
	CountPercentByFilePath(filePath string) PercentResult
	CountCommentsByDirectory(directoryPath string) (PercentResultMap, PercentResult)
}

type CountPercentAnalyzerImpl struct{}

type PercentResult struct {
	TotalLines        int
	CommentLines      int
	CommentPercentage float64
}

type PercentResultMap map[string]PercentResult

func shouldIgnore(name string) bool {
	for _, item := range ignoreList {
		if item == name {
			return true
		}
	}
	return false
}

func (a *CountPercentAnalyzerImpl) CountPercentByFilePath(filePath string) PercentResult {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var result PercentResult
	scanner := bufio.NewScanner(file)
	inBlockComment := false

	for scanner.Scan() {
		line := scanner.Text()
		result.TotalLines++

		if isComment(line, &inBlockComment) { // Usando a função que já existe no pacote
			result.CommentLines++
		}
	}

	if result.TotalLines > 0 {
		result.CommentPercentage = float64(result.CommentLines) / float64(result.TotalLines) * 100
	}

	return result
}

func (a *CountPercentAnalyzerImpl) CountCommentsByDirectory(directoryPath string) (PercentResultMap, PercentResult) {
	if directoryPath == "." {
		var err error
		directoryPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	absPath, err := filepath.Abs(directoryPath)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("directory %s does not exist", absPath))
	}

	linesByArchive := make(PercentResultMap)

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldIgnore(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) == ".js" {
			result := a.CountPercentByFilePath(path)
			linesByArchive[path] = result
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	var total PercentResult
	for _, fileResult := range linesByArchive {
		total.CommentLines += fileResult.CommentLines
		total.TotalLines += fileResult.TotalLines
	}

	if total.TotalLines > 0 {
		total.CommentPercentage = float64(total.CommentLines) / float64(total.TotalLines) * 100
	}

	return linesByArchive, total
}
