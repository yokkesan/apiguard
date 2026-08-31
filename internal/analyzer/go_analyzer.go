package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type GoAnalyzer struct {
	detectors []Detector
	context   *AnalysisContext
}

func NewGoAnalyzer(detectors ...Detector) *GoAnalyzer {
	return &GoAnalyzer{
		detectors: detectors,
	}
}

func (g *GoAnalyzer) LoadFiles(paths []string) []Issue {
	fs := token.NewFileSet()

	g.context = NewAnalysisContext(fs)

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

		g.context.AddFile(path, node)
	}

	return issues
}

func (g *GoAnalyzer) AnalyzeFiles(paths []string) []Issue {
	issues := g.LoadFiles(paths)

	if g.context == nil {
		return issues
	}

	for _, path := range paths {
		fileIssues := g.analyzeLoadedFile(path)
		issues = append(issues, fileIssues...)
	}

	return issues
}

func (g *GoAnalyzer) AnalyzeChangedFile(path string) []Issue {
	if g.context == nil {
		return []Issue{
			{
				Level:   "ERROR",
				Message: "analysis context is not initialized",
				File:    path,
			},
		}
	}

	node, err := parser.ParseFile(
		g.context.FileSet,
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

	g.context.AddFile(path, node)

	return g.analyzeLoadedFile(path)
}

func (g *GoAnalyzer) analyzeLoadedFile(path string) []Issue {
	file, ok := g.context.GetFile(path)
	if !ok {
		return nil
	}

	var issues []Issue

	for _, detector := range g.detectors {
		detected := detector.Detect(
			file.Node,
			file.Path,
			g.context,
		)

		issues = append(issues, detected...)
	}

	routes := DetectGinRoutes(
		file.Node,
		g.context.FileSet,
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

	return issues
}
