package analyzer

import (
	"bufio"
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type AverageFunctionAnalyzer interface {
	CalculateAverageFunctionSize(filePath string) float64
	CalculateAverageFunctionSizeByDirectory(directoryPath string) (map[string]float64, float64)
}

type AverageFunctionAnalyzerImpl struct{}

func (a *AverageFunctionAnalyzerImpl) CalculateAverageFunctionSize(filePath string) float64 {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inFunction := false
	bracesCount := 0
	currentFunctionLines := 0
	totalFunctionLines := 0
	functionCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Verifica início de função
		if !inFunction && isFunctionStart(line) {
			inFunction = true
			bracesCount = strings.Count(line, "{") - strings.Count(line, "}")
			currentFunctionLines = 1
			if bracesCount <= 0 {
				inFunction = false
				totalFunctionLines += currentFunctionLines
				functionCount++
			}
			continue
		}

		if inFunction {
			currentFunctionLines++
			bracesCount += strings.Count(line, "{") - strings.Count(line, "}")

			if bracesCount <= 0 {
				inFunction = false
				totalFunctionLines += currentFunctionLines
				functionCount++
			}
		}
	}

	if functionCount == 0 {
		return 0.0
	}
	return float64(totalFunctionLines) / float64(functionCount)
}

func (a *AverageFunctionAnalyzerImpl) CalculateAverageFunctionSizeByDirectory(directoryPath string) (map[string]float64, float64) {
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

	results := make(map[string]float64)
	var totalSum float64
	var fileCount int

	// Evita repetição e conflitos entre arquivos
	var directoryOrFilesToIgnore = []string{
		"node_modules", ".git", "dist", "build", "coverage", "vendor",
	}

	err := filepath.WalkDir(directoryPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileName := d.Name()
		if slices.Contains(directoryOrFilesToIgnore, fileName) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() && policies.IsJSFileExtension(filepath.Ext(fileName)) {
			average := a.CalculateAverageFunctionSize(path)
			results[fileName] = average
			if average > 0 {
				totalSum += average
				fileCount++
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	var overallAverage float64
	if fileCount > 0 {
		overallAverage = totalSum / float64(fileCount)
	}
	return results, overallAverage
}

func isFunctionStart(line string) bool {
	// suporta function normal, arrow e métodos
	return strings.HasPrefix(line, "function") ||
		strings.Contains(line, "=>") ||
		strings.HasPrefix(line, "async function") ||
		strings.HasSuffix(line, "{") &&
			(strings.Contains(line, "(") && strings.Contains(line, ")"))
}
