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

func (g *GoAnalyzer) AnalyzeFiles(paths []string) []Issue {
	fs := token.NewFileSet()
	context := NewAnalysisContext(fs)

	var issues []Issue

	for _, path := range paths {
		node, err := parser.ParseFile(
			fs,
			path,
			nil,
			parser.ParseComments,
		)

		if err != nil {
			issues = append(issues, Issue{
				Level:   "ERROR",
				Message: fmt.Sprintf("parse error: %v", err),
				File:    path,
			})

			continue
		}

		context.AddFile(path, node)
	}

	for _, file := range context.Files {
		for _, detector := range g.detectors {
			detected := detector.Detect(
				file.Node,
				file.Path,
				context,
			)

			issues = append(issues, detected...)
		}

		routes := DetectGinRoutes(
			file.Node,
			context.FileSet,
		)

		for _, route := range routes {
			message := route.Method + " " + route.Path + " -> " + route.Handler

			if len(route.Middlewares) > 0 {
				message += " middleware: " + strings.Join(route.Middlewares, ", ")
			}

			issues = append(issues, Issue{
				Code:    "INFO-ROUTE-001",
				Level:   "INFO",
				Message: message,
				File:    file.Path,
				Line:    route.Line,
			})
		}
	}

	return issues
}
