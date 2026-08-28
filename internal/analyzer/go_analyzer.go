package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
)

type GoAnalyzer struct {
	detectors []Detector
}

func NewGoAnalyzer(detectors ...Detector) *GoAnalyzer {
	return &GoAnalyzer{
		detectors: detectors,
	}
}

func (g *GoAnalyzer) Analyze(path string) []Issue {
	fs := token.NewFileSet()

	node, err := parser.ParseFile(
		fs,
		path,
		nil,
		parser.ParseComments,
	)
	if err != nil {
		return []Issue{
			{
				Level:   "ERROR",
				Message: fmt.Sprintf("parse error: %v", err),
				File:    path,
			},
		}
	}

	var issues []Issue

	for _, detector := range g.detectors {
		detected := detector.Detect(node, path, fs)
		issues = append(issues, detected...)
	}

	return issues
}
