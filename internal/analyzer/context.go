package analyzer

import (
	"go/ast"
	"go/token"
	"path/filepath"
)

type FileContext struct {
	Path       string
	PackageKey string
	Node       *ast.File
}

type AnalysisContext struct {
	FileSet *token.FileSet
	Files   []FileContext

	// packageKey -> functionName -> function
	Functions map[string]map[string]*ast.FuncDecl

	filePackages map[string]string
}

func NewAnalysisContext(fs *token.FileSet) *AnalysisContext {
	return &AnalysisContext{
		FileSet:      fs,
		Functions:    make(map[string]map[string]*ast.FuncDecl),
		filePackages: make(map[string]string),
	}
}

func (c *AnalysisContext) AddFile(path string, node *ast.File) {
	packageKey := buildPackageKey(path, node.Name.Name)

	c.Files = append(c.Files, FileContext{
		Path:       path,
		PackageKey: packageKey,
		Node:       node,
	})

	c.filePackages[path] = packageKey

	if _, ok := c.Functions[packageKey]; !ok {
		c.Functions[packageKey] = make(map[string]*ast.FuncDecl)
	}

	for _, declaration := range node.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}

		c.Functions[packageKey][function.Name.Name] = function
	}
}

func (c *AnalysisContext) FindFunction(
	file string,
	name string,
) (*ast.FuncDecl, bool) {
	packageKey, ok := c.filePackages[file]
	if !ok {
		return nil, false
	}

	functions, ok := c.Functions[packageKey]
	if !ok {
		return nil, false
	}

	function, ok := functions[name]

	return function, ok
}

func buildPackageKey(path string, packageName string) string {
	dir := filepath.Clean(filepath.Dir(path))

	return dir + ":" + packageName
}
