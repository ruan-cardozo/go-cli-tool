package policies

import (
	"os"
	"path/filepath"
)

func ValidateFilePath(filePath string) bool {
    if filePath == "" {
        return false
    }
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return false
    }
    return true
}

func IsJSFileExtension(filePath string) bool {
	ext := filepath.Ext(filePath)

	return ext == ".js"
}