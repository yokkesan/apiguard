package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
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

	routes := DetectGinRoutes(node, fs)

	for _, route := range routes {
		message := route.Method + " " + route.Path + " -> " + route.Handler

		if len(route.Middlewares) > 0 {
			message += " middleware: " + strings.Join(route.Middlewares, ", ")
		}

		issues = append(issues, Issue{
			Code:    "INFO-ROUTE-001",
			Level:   "INFO",
			Message: message,
			File:    path,
			Line:    route.Line,
		})
	}

	return issues
}
