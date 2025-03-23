package utils

type CommandType string

const (
    COUNT_LINES              CommandType = "Count Lines"
    COUNT_CLASS_AND_FUNCTIONS CommandType = "Count Class And Functions"
    COUNT_COMMENTS           CommandType = "Count Comments"
)