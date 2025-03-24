package policies

import (
	"fmt"
	"go-cli-tool/internal/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

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


func ValidateFilePath(err bool, cmd *cobra.Command) bool {

    if !IsJSFileExtension(utils.FilePath) {
        fmt.Fprintf(cmd.OutOrStdout(), "%sOnly JavaScript files are accepted.%s", utils.RED, utils.RESET_COLOR)
        err = true
    }

    return err
}


func ValidateUserInput(cmd *cobra.Command) bool {

    err := false

    if utils.FilePath == "" && utils.DirectoryPath == "" {
        cmd.Println("Please provide the path to the JavaScript file using the -f flag or use the -d flag to provide the path to the directory containing the JavaScript files.")
        err = true
    }

    if utils.FilePath != "" {
       err = ValidateFilePath(err,cmd)
    }

    return err
}