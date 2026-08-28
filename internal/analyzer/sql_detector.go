package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
)

type SQLDetector struct{}

func (SQLDetector) Detect(
	node ast.Node,
	file string,
	context *AnalysisContext,
) []Issue {
	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {
		binaryExpr, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		if binaryExpr.Op != token.ADD {
			return true
		}

		if !containsSQLString(binaryExpr) {
			return true
		}

		issues = append(issues, Issue{
			Code:    "SEC-SQL-001",
			Level:   "WARNING",
			Message: "SQL string concatenation detected",
			File:    file,
			Line:    context.FileSet.Position(binaryExpr.Pos()).Line,
		})

		return true
	})

	return issues
}

func containsSQLString(expr ast.Expr) bool {
	switch value := expr.(type) {
	case *ast.BasicLit:
		if value.Kind != token.STRING {
			return false
		}

		sql := strings.ToUpper(value.Value)

		keywords := []string{
			"SELECT ",
			"INSERT ",
			"UPDATE ",
			"DELETE ",
		}

		for _, keyword := range keywords {
			if strings.Contains(sql, keyword) {
				return true
			}
		}

	case *ast.BinaryExpr:
		return containsSQLString(value.X) || containsSQLString(value.Y)
	}

	return false
}
