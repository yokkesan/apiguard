package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type SecretDetector struct{}

func (SecretDetector) Detect(
	node ast.Node,
	file string,
	context *AnalysisContext,
) []Issue {
	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, lhs := range assign.Lhs {
			if i >= len(assign.Rhs) {
				continue
			}

			ident, ok := lhs.(*ast.Ident)
			if !ok {
				continue
			}

			if !isSecretName(ident.Name) {
				continue
			}

			literal, ok := assign.Rhs[i].(*ast.BasicLit)
			if !ok || literal.Kind != token.STRING {
				continue
			}

			value, err := strconv.Unquote(literal.Value)
			if err != nil || value == "" {
				continue
			}

			issues = append(issues, Issue{
				Code:    "SEC-SECRET-001",
				Level:   "WARNING",
				Message: "Hard-coded secret detected",
				File:    file,
				Line:    context.FileSet.Position(assign.Pos()).Line,
			})
		}

		return true
	})

	return issues
}

func isSecretName(name string) bool {
	name = strings.ToLower(name)

	keywords := []string{
		"password",
		"passwd",
		"secret",
		"apikey",
		"api_key",
		"token",
	}

	for _, keyword := range keywords {
		if strings.Contains(name, keyword) {
			return true
		}
	}

	return false
}
