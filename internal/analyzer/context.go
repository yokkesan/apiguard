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
	Files   map[string]FileContext

	// packageKey -> functionName -> function
	Functions map[string]map[string]*ast.FuncDecl

	filePackages  map[string]string
	fileFunctions map[string][]string
}

func NewAnalysisContext(fs *token.FileSet) *AnalysisContext {
	return &AnalysisContext{
		FileSet:       fs,
		Files:         make(map[string]FileContext),
		Functions:     make(map[string]map[string]*ast.FuncDecl),
		filePackages:  make(map[string]string),
		fileFunctions: make(map[string][]string),
	}
}

func (c *AnalysisContext) AddFile(path string, node *ast.File) {
	c.removeFileData(path)

	packageKey := buildPackageKey(path, node.Name.Name)

	c.Files[path] = FileContext{
		Path:       path,
		PackageKey: packageKey,
		Node:       node,
	}

	c.filePackages[path] = packageKey

	if _, ok := c.Functions[packageKey]; !ok {
		c.Functions[packageKey] = make(map[string]*ast.FuncDecl)
	}

	var functionNames []string

	for _, declaration := range node.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := function.Name.Name

		c.Functions[packageKey][name] = function
		functionNames = append(functionNames, name)
	}

	c.fileFunctions[path] = functionNames
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

func (c *AnalysisContext) GetFile(path string) (FileContext, bool) {
	file, ok := c.Files[path]
	return file, ok
}

func (c *AnalysisContext) removeFileData(path string) {
	packageKey, ok := c.filePackages[path]
	if !ok {
		return
	}

	functionNames := c.fileFunctions[path]

	if functions, exists := c.Functions[packageKey]; exists {
		for _, name := range functionNames {
			delete(functions, name)
		}

		if len(functions) == 0 {
			delete(c.Functions, packageKey)
		}
	}

	delete(c.Files, path)
	delete(c.filePackages, path)
	delete(c.fileFunctions, path)
}

func buildPackageKey(path string, packageName string) string {
	dir := filepath.Clean(filepath.Dir(path))

	return dir + ":" + packageName
}
