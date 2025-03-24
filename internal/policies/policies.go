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

func ValidateDirectoryPath(directoryPath string) bool {

    info, err := os.Stat(directoryPath)
    if os.IsNotExist(err) {
        return false
    }
    return info.IsDir()
}

func IsJSFileExtension(filePath string) bool {
	ext := filepath.Ext(filePath)

	return ext == ".js" || ext == ".mjs"
}