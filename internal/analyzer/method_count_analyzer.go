package analyzer

import (
	"bufio"
	"go-cli-tool/internal/utils"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type MethodCountResult struct {
	Public  int
	Private int
}

type MethodCountMap map[string]MethodCountResult

type MethodCountAnalyzer interface {
	AnalyzeFile(filePath string) MethodCountResult
	AnalyzeDirectory(dirPath string) (MethodCountMap, MethodCountResult)
}

type MethodCountAnalyzerImpl struct{}

var filesToIgnore = []string{".git", "node_modules"}

func (a *MethodCountAnalyzerImpl) AnalyzeFile(filePath string) MethodCountResult {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var result MethodCountResult
	scanner := bufio.NewScanner(file)

	publicMethodRegex := regexp.MustCompile(`^\s*([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\([^)]*\)\s*\{`)
	privateMethodRegex := regexp.MustCompile(`^\s*([#_][a-zA-Z_$][a-zA-Z0-9_$]*)\s*\([^)]*\)\s*\{`)
	functionRegex := regexp.MustCompile(`^\s*function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\([^)]*\)\s*\{`)
	constFunctionRegex := regexp.MustCompile(`^\s*const\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*function`)
	arrowFunctionRegex := regexp.MustCompile(`^\s*const\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*\([^)]*\)\s*=>`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Pula comentários e linhas vazias
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}

		// Verifica métodos privados primeiro (com # ou _)
		if privateMethodRegex.MatchString(line) {
			result.Private++
		} else if publicMethodRegex.MatchString(line) && !strings.Contains(line, "function") {
			// Métodos de classe públicos (dentro de classes)
			result.Public++
		} else if functionRegex.MatchString(line) {
			// Funções declaradas (públicas por padrão)
			result.Public++
		} else if constFunctionRegex.MatchString(line) {
			// const funcName = function() {}
			funcName := constFunctionRegex.FindStringSubmatch(line)
			if len(funcName) > 1 {
				if strings.HasPrefix(funcName[1], "_") || strings.HasPrefix(funcName[1], "#") {
					result.Private++
				} else {
					result.Public++
				}
			}
		} else if arrowFunctionRegex.MatchString(line) {
			// const funcName = () => {}
			funcName := arrowFunctionRegex.FindStringSubmatch(line)
			if len(funcName) > 1 {
				if strings.HasPrefix(funcName[1], "_") || strings.HasPrefix(funcName[1], "#") {
					result.Private++
				} else {
					result.Public++
				}
			}
		}
	}

	return result
}

func (a *MethodCountAnalyzerImpl) AnalyzeDirectory(dirPath string) (MethodCountMap, MethodCountResult) {
	if dirPath == "." {
		var err error
		dirPath, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		dirPath, err = utils.ExpandPath(dirPath)
		if err != nil {
			panic(err)
		}
	}

	results := make(MethodCountMap)
	var total MethodCountResult

	filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if slices.Contains(filesToIgnore, d.Name()) && d.IsDir() {
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".js" || filepath.Ext(path) == ".mjs" {
			count := a.AnalyzeFile(path)
			results[d.Name()] = count
			total.Public += count.Public
			total.Private += count.Private
		}

		return nil
	})

	return results, total
}
