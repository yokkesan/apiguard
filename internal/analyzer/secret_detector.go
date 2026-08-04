package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
)

func DetectSecret(node ast.Node, file string, fs *token.FileSet) []Issue {
	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {

		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		for i, expr := range assign.Lhs {

			ident, ok := expr.(*ast.Ident)
			if !ok {
				continue
			}

			name := strings.ToLower(ident.Name)

			if !isSecretName(name) {
				continue
			}

			if i >= len(assign.Rhs) {
				continue
			}

			_, ok = assign.Rhs[i].(*ast.BasicLit)
			if !ok {
				continue
			}

			issues = append(issues, Issue{
				Code:    "SEC-SECRET-001",
				Level:   "WARNING",
				Message: "Secret information detected",
				File:    file,
				Line:    fs.Position(ident.Pos()).Line,
			})
		}

		return true
	})

	return issues
}

func isSecretName(name string) bool {
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
