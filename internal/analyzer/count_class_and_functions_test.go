package analyzer_test

import (
	"go-cli-tool/internal/analyzer"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountClassesAndFunctionsByFilePath(t *testing.T) {
    // Create a temporary JavaScript file for testing
	content := `
	class Teste {

		constructor() {
		}
	}

	function testFunction() { };

	const arrowFunction = () => {
	}

	const myFunction = function() {
	}

	() => {
	}

	var oldVarFybc = () => {}

	var teste = setTimeout(() => {}, teste)

	const filter = array.map(() => {});
	`

    tmpFile, err := os.CreateTemp("", "test.js")
    assert.NoError(t, err)
    defer os.Remove(tmpFile.Name())

    _, err = tmpFile.WriteString(content)
    assert.NoError(t, err)
    tmpFile.Close()

    // Create an instance of CountClassAndFunctionsImpl
    analyzer := &analyzer.CountClassAndFunctionsImpl{}

    // Call the function to test
    result := analyzer.CountClassesAndFunctionsByFilePath(tmpFile.Name())

    // Assert the results
    assert.Equal(t, 1, result.Classes, "Expected 1 class")
    assert.Equal(t, 8, result.Functions, "Expected 8 functions")
}

func TestCountClassesAndFunctionsByDirectory(t *testing.T) {
    // Create a temporary directory for testing
    tmpDir, err := os.MkdirTemp("", "testdir")
    assert.NoError(t, err)
    defer os.RemoveAll(tmpDir)

    // Create test JavaScript files
    files := map[string]string{
        "file1.js": `
            class TestClass1 {
                constructor() {}
            }

            function testFunction1() {}

            const arrowFunction1 = () => {};
        `,
        "file2.js": `
            class TestClass2 {
                constructor() {}
            }

            function testFunction2() {}

            const arrowFunction2 = () => {};
        `,
    }

    for fileName, content := range files {
        filePath := filepath.Join(tmpDir, fileName)
        err := os.WriteFile(filePath, []byte(content), 0644)
		assert.NoError(t, err)
    }

    // Create an instance of CountClassAndFunctionsImpl
    analyzer := &analyzer.CountClassAndFunctionsImpl{}

    // Call the function to test
    linesByArchive, totalClassesAndFunctions := analyzer.CountClassesAndFunctionsByDirectory(tmpDir)

    // Assert the results for each file
    assert.Equal(t, 1, linesByArchive["file1.js"].Classes, "Expected 1 class in file1.js")
    assert.Equal(t, 3, linesByArchive["file1.js"].Functions, "Expected 2 functions in file1.js")

    assert.Equal(t, 1, linesByArchive["file2.js"].Classes, "Expected 1 class in file2.js")
    assert.Equal(t, 3, linesByArchive["file2.js"].Functions, "Expected 2 functions in file2.js")

    // Assert the total results
    assert.Equal(t, 2, totalClassesAndFunctions.Classes, "Expected 2 classes in total")
    assert.Equal(t, 6, totalClassesAndFunctions.Functions, "Expected 4 functions in total")
}