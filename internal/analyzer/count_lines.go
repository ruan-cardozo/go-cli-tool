package analyzer

import (
	"bufio"
	"os"
)

func CountLines(filePath string) LineResult {
    file, err := os.Open(filePath)

    if err != nil {
        panic(err)
    }

    defer file.Close()

    var result LineResult
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

        line := scanner.Text()

        if isEmptyLine(line) {
            continue
        }

        result.TotalLines++
    }

    return result
}

func isEmptyLine(line string) bool {
    return len(line) == 0
}