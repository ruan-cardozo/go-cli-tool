package analyzer

type LineResult struct {
	TotalLines int
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

type ClassesAndFunctionsMap map[string]ClassFuncResult

type CountClassesAndFunctionsAnalyzer interface {
	CountClassesAndFunctionsByFilePath(filePath string) ClassFuncResult
	CountClassesAndFunctionsByDirectory(directoryPath string) (ClassesAndFunctionsMap, ClassFuncResult)
}

type CommentsMap map[string]CommentResult
