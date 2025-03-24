package analyzer_test

import (
	"go-cli-tool/internal/analyzer"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountLinesByFilePath(t *testing.T) {

	content := `
		class TestClass2 {

			constructor() {
			}

			testFunction() { };	
		}

		function testFunction2() {}

        const arrowFunction2 = () => {};
`

    tmpFile, err := os.CreateTemp("", "test.js")
    assert.NoError(t, err)
    defer os.Remove(tmpFile.Name())

    _, err = tmpFile.WriteString(content)
    assert.NoError(t, err)
    tmpFile.Close()

    analyzer := &analyzer.CountLinesAnalyzerImpl{}

    result := analyzer.CountLinesByFilePath(tmpFile.Name())

    assert.Equal(t, 7, result.TotalLines, "Expected 1 class")

	result.TotalLines = 0
}

func TestCountLinesByDirectory(t *testing.T) {

    tmpDir, err := os.MkdirTemp("", "testdir")
    assert.NoError(t, err)
    defer os.RemoveAll(tmpDir)

    files := map[string]string{
        "file1.js": `
            class TestClass1 {
                constructor() {}
            }

            function testFunction1() {}

            const arrowFunction1 = () => {}; `,
        "file2.js": `
            class TestClass2 {
                constructor() {}
            }

            function testFunction2() {}

            const arrowFunction2 = () => {};`,
    }

    for fileName, content := range files {
        filePath := filepath.Join(tmpDir, fileName)
        err := os.WriteFile(filePath, []byte(content), 0644)
        assert.NoError(t, err)
    }

    analyzer := &analyzer.CountLinesAnalyzerImpl{}


    linesByArchive, totalLines := analyzer.CountLinesByDirectory(tmpDir)

    assert.Equal(t, 5, linesByArchive["file1.js"].TotalLines, "Expected 7 lines in file1.js")
    assert.Equal(t, 5, linesByArchive["file2.js"].TotalLines, "Expected 7 lines in file2.js")

    assert.Equal(t, 10, totalLines.TotalLines, "Expected 14 lines in total")

	totalLines.TotalLines = 0
}

func TestCountLinesByDirectoryToSkipFileNamedWithDirectoryToIgnore(t *testing.T) {
    tmpDir, err := os.MkdirTemp("", "testdir")
    assert.NoError(t, err)
    defer os.RemoveAll(tmpDir)

    files := map[string]string{
        ".git":         ` `,
        "node_modules": ``,
    }

    for fileName, content := range files {
        filePath := filepath.Join(tmpDir, fileName)
        err := os.WriteFile(filePath, []byte(content), 0644)
        assert.NoError(t, err)
    }

    analyzer := &analyzer.CountLinesAnalyzerImpl{}

    linesByArchive, totalLines := analyzer.CountLinesByDirectory(tmpDir)

    assert.Equal(t, 0, linesByArchive[".git"].TotalLines, "Expected 0 lines in .git")
    assert.Equal(t, 0, linesByArchive["node_modules"].TotalLines, "Expected 0 lines in node_modules")
    assert.Equal(t, 0, totalLines.TotalLines, "Expected to skip all files")
}

func TestCountLinesByDirectoryToSkipDirectoryToIgnore(t *testing.T) {
    parentDir, err := os.MkdirTemp("", "parentdir")
    assert.NoError(t, err)
    defer os.RemoveAll(parentDir)

    ignoredDirs := []string{".git", "node_modules"}

    for _, dirName := range ignoredDirs {
        dirPath := filepath.Join(parentDir, dirName)
        err := os.Mkdir(dirPath, 0755)
        assert.NoError(t, err)

        filePath := filepath.Join(dirPath, "file.js")
        err = os.WriteFile(filePath, []byte(`class TestClass {}`), 0644)
        assert.NoError(t, err)
    }

    analyzer := &analyzer.CountLinesAnalyzerImpl{}

    linesByArchive, totalLines := analyzer.CountLinesByDirectory(parentDir)

    for _, dirName := range ignoredDirs {
        assert.Equal(t, 0, linesByArchive[dirName].TotalLines, "Expected 0 lines in "+dirName)
    }
    assert.Equal(t, 0, totalLines.TotalLines, "Expected to skip all files")
}