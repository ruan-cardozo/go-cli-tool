package models

import (
	"time"
)

type DependencyDto struct {
	Dependencies      []string `json:"dependencies" validate:"required"`
	NativeModules     []string `json:"native_modules" validate:"required"`
	TotalDependencies int      `json:"total_dependencies" validate:"required,min=0"`
}

type IndentDistributionDto struct {
	Level int `json:"level" validate:"required,min=0"`
	Count int `json:"count" validate:"required,min=0"`
}

type IndentationFileStatsDto struct {
	MaxIndentLevel       int                     `json:"maxIndentLevel" validate:"required,min=0"`
	AverageIndentLevel   float64                 `json:"averageIndentLevel" validate:"required"`
	IndentDistribution   []IndentDistributionDto `json:"indentDistribution" validate:"required"`
	UsesSpaces           bool                    `json:"usesSpaces"`
	UsesTabs             bool                    `json:"usesTabs"`
	MixedIndentation     bool                    `json:"mixedIndentation"`
}

type IndentationFileDto struct {
	Filename string                  `json:"filename" validate:"required"`
	Path     string                  `json:"path" validate:"required"`
	Stats    IndentationFileStatsDto `json:"stats" validate:"required"`
}

type IdentationDto struct {
	Directory string                `json:"directory" validate:"required"`
	Files     []IndentationFileDto  `json:"files" validate:"required"`
}

type CreateMetricDto struct {
	RecordedAt        time.Time     `json:"recorded_at" validate:"required"`
	Lines             int           `json:"lines" validate:"required,min=0"`
	Functions         int           `json:"functions" validate:"required,min=0"`
	Classes           int           `json:"classes" validate:"required,min=0"`
	Comments          int           `json:"comments" validate:"required,min=0"`
	CommentPercentage string        `json:"comment_percentage" validate:"required,comment_percentage"`
	Dependencies      DependencyDto `json:"dependencies" validate:"required"`
	Indentation       IdentationDto `json:"indentation" validate:"required"`
}

type CreateMetricDtoWithStringDate struct {
	RecordedAt        string        `json:"recorded_at" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	Lines             int           `json:"lines" validate:"required,min=0"`
	Functions         int           `json:"functions" validate:"required,min=0"`
	Classes           int           `json:"classes" validate:"required,min=0"`
	Comments          int           `json:"comments" validate:"required,min=0"`
	CommentPercentage string        `json:"comment_percentage" validate:"required,comment_percentage"`
	Dependencies      DependencyDto `json:"dependencies" validate:"required"`
	Indentation       IdentationDto `json:"indentation" validate:"required"`
}