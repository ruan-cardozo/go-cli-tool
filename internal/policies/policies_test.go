package policies_test

import (
	"go-cli-tool/internal/policies"
	"go-cli-tool/internal/utils"
	"go-cli-tool/tests"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestValidateDirectoryPath(t *testing.T) {
	// Test case: Directory exists
	dirPath := t.TempDir()
	assert.True(t, policies.ValidateDirectoryPath(dirPath), "Expected true for existing directory")

	// Test case: Directory does not exist
	nonExistentDirPath := "non_existent_directory"
	assert.False(t, policies.ValidateDirectoryPath(nonExistentDirPath), "Expected false for non-existent directory")

	// Test case: Path is a file, not a directory
	file, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	assert.False(t, policies.ValidateDirectoryPath(file.Name()), "Expected false for a file path")
}

func TestIsJsFileExetension(t *testing.T) {
	// Test case: File has .js extension
	filePath := "file.js"
	assert.True(t, policies.IsJSFileExtension(filePath), "Expected true for .js file extension")

	// Test case: File has .mjs extension
	filePath = "file.mjs"
	assert.True(t, policies.IsJSFileExtension(filePath), "Expected true for .mjs file extension")

	// Test case: File has .txt extension
	filePath = "file.txt"
	assert.False(t, policies.IsJSFileExtension(filePath), "Expected false for .txt file extension")
}

func TestValidateFilePath(t *testing.T) {

	var err bool
	utils.FilePath = "file.js"
	cmd := &cobra.Command{}
	assert.False(t, policies.ValidateFilePath(err, cmd), "Expected false for .js file extension")

	utils.FilePath = "file.mjs"
	cmd = &cobra.Command{}
	assert.False(t, policies.ValidateFilePath(err, cmd), "Expected false for .mjs file extension")

	utils.FilePath = "file.txt"
	cmd = &cobra.Command{}
	assert.True(t, policies.ValidateFilePath(err, cmd), "Expected true for .txt file extension")
}

func TestValidateUserInput(t *testing.T) {
	cmd := &cobra.Command{}

	tests.ResetGlobals()
	assert.True(t, policies.ValidateUserInput(cmd), "Expected true for no file path or directory path provided")

	tests.ResetGlobals()
	utils.FilePath = "file.js"
	assert.False(t, policies.ValidateUserInput(cmd), "Expected false for file path provided")

	tests.ResetGlobals()
	utils.DirectoryPath = t.TempDir()
	assert.False(t, policies.ValidateUserInput(cmd), "Expected false for directory path provided")
}