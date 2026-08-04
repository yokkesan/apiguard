package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
)

type GoAnalyzer struct{}

func (g GoAnalyzer) Analyze(path string) []Issue {
	fs := token.NewFileSet()

	_, err := parser.ParseFile(
		fs,
		path,
		nil,
		parser.ParseComments,
	)

	if err != nil {
		return []Issue{
			{
				File:    path,
				Message: fmt.Sprintf("parse error: %v", err),
				Level:   "ERROR",
			},
		}
	}

	return []Issue{
		{
			File:    path,
			Message: "Go analysis completed",
			Level:   "INFO",
		},
	}
}
