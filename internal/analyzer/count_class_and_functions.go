package analyzer

import (
	"bufio"
	"os"
	"regexp"
)

func CountClassesAndFunctions(filePath string) ClassFuncResult {
    file, err := os.Open(filePath)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    var result ClassFuncResult
    scanner := bufio.NewScanner(file)

    functionRegex := regexp.MustCompile(`function\s+\w+|\w+\s*=\s*function|\w+\s*=>`)
    classRegex := regexp.MustCompile(`class\s+\w+`)

    for scanner.Scan() {
        line := scanner.Text()

        if functionRegex.MatchString(line) {
            result.Functions++
        }

        if classRegex.MatchString(line) {
            result.Classes++
        }
    }

    return result
}