package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type GoAnalyzer struct{}

func (g GoAnalyzer) Analyze(path string) []Issue {
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
				File:    path,
				Message: fmt.Sprintf("parse error: %v", err),
				Level:   "ERROR",
			},
		}
	}

	issues := DetectSecret(node, path, fs)

	sqlIssues := DetectSQLInjection(node, path, fs)
	issues = append(issues, sqlIssues...)

	routes := DetectGinRoutes(node, fs)

	authIssues := DetectMissingAuth(routes, path)
	issues = append(issues, authIssues...)

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
