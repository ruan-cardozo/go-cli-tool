package analyzer

type LineResult struct {
    TotalLines   int
}

type CommentResult struct {
    CommentLines int
}

type ClassFuncResult struct {
    Functions int
    Classes   int
}

type CountLinesAnalyzer interface {
    CountLinesByFilePath(filePath string) LineResult
    CountLinesByDirectory(directoryPath string) (FilesNameCountLineMap, LineResult)
}