package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

type AuthDetector struct{}

func (AuthDetector) Detect(
	node ast.Node,
	file string,
	context *AnalysisContext,
) []Issue {
	var issues []Issue

	routes := DetectGinRoutes(node, context.FileSet)

	for _, route := range routes {
		if len(route.Middlewares) == 0 {
			issues = append(issues, Issue{
				Code:    "SEC-AUTH-001",
				Level:   "WARNING",
				Message: "Authentication middleware is not configured",
				File:    file,
				Line:    route.Line,
			})

			continue
		}

		if hasAuthenticationMiddleware(file, route.Middlewares, context) {
			continue
		}

		issues = append(issues, Issue{
			Code:    "SEC-AUTH-001",
			Level:   "WARNING",
			Message: "Authentication middleware could not be verified",
			File:    file,
			Line:    route.Line,
		})
	}

	return issues
}

func hasAuthenticationMiddleware(
	file string,
	middlewares []string,
	context *AnalysisContext,
) bool {
	for _, middlewareName := range middlewares {
		function, ok := context.FindFunction(file, middlewareName)
		if !ok {
			continue
		}

		if containsAuthenticationLogic(function) {
			return true
		}
	}

	return false
}

func containsAuthenticationLogic(function *ast.FuncDecl) bool {
	if function.Body == nil {
		return false
	}

	foundAuthorizationHeader := false

	ast.Inspect(function.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if selector.Sel.Name != "GetHeader" {
			return true
		}

		if len(call.Args) != 1 {
			return true
		}

		literal, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}

		if literal.Kind != token.STRING {
			return true
		}

		value, err := strconv.Unquote(literal.Value)
		if err != nil {
			return true
		}

		if value == "Authorization" {
			foundAuthorizationHeader = true
			return false
		}

		return true
	})

	return foundAuthorizationHeader
}
