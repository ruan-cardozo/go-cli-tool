package analyzer

import (
	"bufio"
	"fmt"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

// Definindo a interface CountPercentAnalyzer
type CountPercentAnalyzer interface {
	CountPercentByFilePath(filePath string) PercentResult
	CountCommentsByDirectory(directoryPath string) (PercentResultMap, PercentResult)
}

// Implementando a interface CountPercentAnalyzer
type CountPercentAnalyzerImpl struct{}

// Estrutura para armazenar o resultado do cálculo
type PercentResult struct {
	TotalLines        int
	CommentLines      int
	CommentPercentage float64
	TotalComments     int
}

// Função para contar o percentual de comentários em um arquivo
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

		if isComment(line, &inBlockComment) {
			result.CommentLines++
		}
	}

	if result.TotalLines > 0 {
		result.CommentPercentage = float64(result.CommentLines) / float64(result.TotalLines) * 100
	}

	return result
}

// Função para contar os comentários em todos os arquivos de um diretório
func (a *CountPercentAnalyzerImpl) CountCommentsByDirectory(directoryPath string) (PercentResultMap, PercentResult) {
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

	linesByArchive := make(PercentResultMap)

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
			linesByArchive[fileOrDirectoryName] = a.CountPercentByFilePath(path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	var totalCommentsByDirectory PercentResult

	for result := range linesByArchive {
		file := linesByArchive[result]
		totalCommentsByDirectory.CommentLines += file.CommentLines
		totalCommentsByDirectory.TotalLines += file.TotalLines
		totalCommentsByDirectory.TotalComments += file.CommentLines
	}

	if totalCommentsByDirectory.TotalLines > 0 {
		totalCommentsByDirectory.CommentPercentage = float64(totalCommentsByDirectory.CommentLines) / float64(totalCommentsByDirectory.TotalLines) * 100
	}

	return linesByArchive, totalCommentsByDirectory
}

type PercentResultMap map[string]PercentResult
