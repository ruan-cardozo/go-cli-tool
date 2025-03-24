package tests

import "go-cli-tool/internal/utils"

func ResetGlobals() {
	utils.FilePath = ""
	utils.DirectoryPath = ""
	utils.OutputFilePath = ""
}