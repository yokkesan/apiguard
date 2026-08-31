package analyzer

import (
	"go/ast"
	"strings"
)

type ValidationDetector struct{}

func (ValidationDetector) Detect(
	node ast.Node,
	file string,
	context *AnalysisContext,
) []Issue {
	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {
		function, ok := n.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			return true
		}

		hasInput := false
		hasValidation := false

		ast.Inspect(function.Body, func(bodyNode ast.Node) bool {
			call, ok := bodyNode.(*ast.CallExpr)
			if !ok {
				return true
			}

			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			method := selector.Sel.Name

			if isInputMethod(method) {
				hasInput = true
			}

			if isValidationMethod(method) {
				hasValidation = true
			}

			return true
		})

		if hasInput && !hasValidation {
			issues = append(issues, Issue{
				Code:    "SEC-VALIDATION-001",
				Level:   "WARNING",
				Message: "Input validation could not be verified",
				File:    file,
				Line:    context.FileSet.Position(function.Pos()).Line,
			})
		}

		return true
	})

	return issues
}

func isInputMethod(method string) bool {
	switch method {
	case
		"ShouldBind",
		"ShouldBindJSON",
		"ShouldBindXML",
		"Bind",
		"BindJSON",
		"BindXML",
		"PostForm",
		"Query",
		"Param":
		return true
	default:
		return false
	}
}

func isValidationMethod(method string) bool {
	method = strings.ToLower(method)

	switch method {
	case
		"shouldbind",
		"shouldbindjson",
		"shouldbindxml",
		"bind",
		"bindjson",
		"bindxml",
		"validate",
		"struct":
		return true
	default:
		return false
	}
}
