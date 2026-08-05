package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
)

func DetectSQLInjection(node ast.Node, file string, fs *token.FileSet) []Issue {
	var issues []Issue

	ast.Inspect(node, func(n ast.Node) bool {

		binary, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}

		// + 以外は対象外
		if binary.Op.String() != "+" {
			return true
		}

		// 左右どちらかがSQL文字列か確認
		if isSQLString(binary.X) || isSQLString(binary.Y) {
			issues = append(issues, Issue{
				Code:    "SEC-SQL-001",
				Level:   "WARNING",
				Message: "SQL string concatenation detected",
				File:    file,
				Line:    fs.Position(binary.Pos()).Line,
			})
		}

		return true
	})

	return issues
}

func isSQLString(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return false
	}

	value := strings.ToLower(lit.Value)

	keywords := []string{
		"select",
		"insert",
		"update",
		"delete",
		"from",
		"where",
	}

	for _, keyword := range keywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}

	return false
}
