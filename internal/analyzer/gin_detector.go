package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

type Route struct {
	Method      string
	Path        string
	Handler     string
	Middlewares []string
	Line        int
}

func DetectGinRoutes(node ast.Node, fs *token.FileSet) []Route {
	var routes []Route

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		method := selector.Sel.Name

		if !isHTTPMethod(method) {
			return true
		}

		if len(call.Args) < 2 {
			return true
		}

		pathLiteral, ok := call.Args[0].(*ast.BasicLit)
		if !ok || pathLiteral.Kind != token.STRING {
			return true
		}

		path, err := strconv.Unquote(pathLiteral.Value)
		if err != nil {
			return true
		}

		lastIndex := len(call.Args) - 1

		handler := expressionName(call.Args[lastIndex])
		if handler == "" {
			return true
		}

		var middlewares []string

		for _, arg := range call.Args[1:lastIndex] {
			name := expressionName(arg)
			if name == "" {
				continue
			}

			middlewares = append(middlewares, name)
		}

		routes = append(routes, Route{
			Method:      method,
			Path:        path,
			Handler:     handler,
			Middlewares: middlewares,
			Line:        fs.Position(call.Pos()).Line,
		})

		return true
	})

	return routes
}

func isHTTPMethod(method string) bool {
	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

func expressionName(expr ast.Expr) string {
	switch value := expr.(type) {
	case *ast.Ident:
		return value.Name

	case *ast.SelectorExpr:
		prefix := expressionName(value.X)
		if prefix == "" {
			return value.Sel.Name
		}

		return prefix + "." + value.Sel.Name
	}

	return ""
}
