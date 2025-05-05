package analyzer

import (
	"fmt"
	"go-cli-tool/internal/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IndentResult represents the indentation statistics for a file
type IndentResult struct {
    MaxIndentLevel     int         `json:"maxIndentLevel"`
    AverageIndentLevel float64     `json:"averageIndentLevel"`
    IndentDistribution []IndentFreq `json:"indentDistribution"`
    UsesSpaces         bool        `json:"usesSpaces"`
    UsesTabs           bool        `json:"usesTabs"`
    MixedIndentation   bool        `json:"mixedIndentation"`
}

// IndentFreq represents frequency of a particular indentation level
type IndentFreq struct {
    Level int `json:"level"`
    Count int `json:"count"`
}

// FileIndentMap maps filenames to their indentation results
type FileIndentMap map[string]IndentResult

// IdentationAnalyzerImpl implements indentation analysis functionality
type IdentationAnalyzerImpl struct{}

// IdentationByFilePath analyzes indentation based on provided file or directory path
func (a *IdentationAnalyzerImpl) IdentationByFilePath() (map[string]interface{}, error) {
    var path string
    var results map[string]interface{}
    var err error

    // Determine which path to use
    if utils.FilePath != "" {
        path = utils.FilePath
        results, err = a.analyzeFileIndentation(path)
    } else if utils.DirectoryPath != "" {
        path = utils.DirectoryPath
        results, err = a.analyzeDirectoryIndentation(path)
    } else {
        return nil, fmt.Errorf("no file or directory path provided")
    }

    if err != nil {
        return nil, err
    }

    return results, nil
}

// analyzeFileIndentation analyzes indentation for a single JavaScript file
func (a *IdentationAnalyzerImpl) analyzeFileIndentation(filePath string) (map[string]interface{}, error) {
    // Check if file is JavaScript
    if !strings.HasSuffix(filePath, ".js") {
        return nil, fmt.Errorf("file %s is not a JavaScript file", filePath)
    }
    
    // Read file
    content, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
    }
    
    lines := strings.Split(string(content), "\n")
    indentationStats := a.calculateIndentationStats(lines)
    
    // Create results
    results := map[string]interface{}{
        "filename": filepath.Base(filePath),
        "path":     filePath,
        "stats":    indentationStats,
    }
    
    return results, nil
}

// analyzeDirectoryIndentation analyzes indentation for all JavaScript files in a directory
func (a *IdentationAnalyzerImpl) analyzeDirectoryIndentation(dirPath string) (map[string]interface{}, error) {
    var allFiles []string
    
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && strings.HasSuffix(path, ".js") {
            allFiles = append(allFiles, path)
        }
        return nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("error walking directory %s: %w", dirPath, err)
    }
    
    filesResults := make([]map[string]interface{}, 0, len(allFiles))
    
    for _, file := range allFiles {
        fileResults, err := a.analyzeFileIndentation(file)
        if err != nil {
            return nil, err
        }
        filesResults = append(filesResults, fileResults)
    }
    
    results := map[string]interface{}{
        "directory": dirPath,
        "files":     filesResults,
    }
    
    return results, nil
}

// calculateIndentationStats calculates indentation statistics for a slice of lines
func (a *IdentationAnalyzerImpl) calculateIndentationStats(lines []string) IndentResult {
    maxIndent := 0
    totalIndent := 0
    indentCount := 0
    indentDistribution := make(map[int]int)
    usesSpaces := false
    usesTabs := false
    
    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)
        if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
            // Skip empty lines and comments
            continue
        }
        
        indentLevel := 0
        for _, char := range line {
            if char == ' ' {
                indentLevel++
                usesSpaces = true
            } else if char == '\t' {
                indentLevel += 4 // Count tabs as 4 spaces (common convention)
                usesTabs = true
            } else {
                break
            }
        }
        
        indentDistribution[indentLevel]++
        
        if indentLevel > maxIndent {
            maxIndent = indentLevel
        }
        
        totalIndent += indentLevel
        indentCount++
    }
    
    avgIndent := 0.0
    if indentCount > 0 {
        avgIndent = float64(totalIndent) / float64(indentCount)
    }
    
    // Convert distribution to sorted slice for easier analysis
    distribution := make([]IndentFreq, 0, len(indentDistribution))
    for level, count := range indentDistribution {
        distribution = append(distribution, IndentFreq{Level: level, Count: count})
    }
    
    // Sort by indentation level
    sort.Slice(distribution, func(i, j int) bool {
        return distribution[i].Level < distribution[j].Level
    })
    
    return IndentResult{
        MaxIndentLevel:     maxIndent,
        AverageIndentLevel: avgIndent,
        IndentDistribution: distribution,
        UsesSpaces:         usesSpaces,
        UsesTabs:           usesTabs,
        MixedIndentation:   usesSpaces && usesTabs,
    }
}