package utils_test

import (
	"go-cli-tool/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	// Test case: Path is already expanded
	path := "/home/user/file.js"
	expandedPath, err := utils.ExpandPath(path)
	assert.NoError(t, err)
	assert.Equal(t, path, expandedPath)

	// Test case: Path is not expanded
	path = "~/file.js"
	expandedPath, err = utils.ExpandPath(path)
	assert.NoError(t, err)
	assert.NotEqual(t, path, expandedPath)
}