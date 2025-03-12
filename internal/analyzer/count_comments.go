package analyzer

import (
	"bufio"
	"os"
	"regexp"
)

func CountComments(filePath string) CommentResult {
    file, err := os.Open(filePath)

    if err != nil {
        panic(err)
    }

    defer file.Close()

    var result CommentResult
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()

        if isComment(line) {
            result.CommentLines++
        }
    }

    return result
}

func isComment(line string) bool {
	singleLineComment := regexp.MustCompile(`^\s*//`)
	blockCommentStart := regexp.MustCompile(`^\s*/\*`)
	blockCommentEnd := regexp.MustCompile(`\*/`)

	return singleLineComment.MatchString(line) || blockCommentStart.MatchString(line) || blockCommentEnd.MatchString(line)
}