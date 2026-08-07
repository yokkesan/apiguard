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
		if !ok {
			return true
		}

		path, err := strconv.Unquote(pathLiteral.Value)
		if err != nil {
			return true
		}

		lastIndex := len(call.Args) - 1

		handler, ok := call.Args[lastIndex].(*ast.Ident)
		if !ok {
			return true
		}

		var middlewares []string

		for _, arg := range call.Args[1:lastIndex] {
			middleware, ok := arg.(*ast.Ident)
			if !ok {
				continue
			}

			middlewares = append(middlewares, middleware.Name)
		}

		routes = append(routes, Route{
			Method:      method,
			Path:        path,
			Handler:     handler.Name,
			Middlewares: middlewares,
			Line:        fs.Position(call.Pos()).Line,
		})

		return true
	})

	return routes
}

func isHTTPMethod(method string) bool {
	methods := []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
	}

	for _, m := range methods {
		if method == m {
			return true
		}
	}

	return false
}
