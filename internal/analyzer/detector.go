package analyzer

import (
	"go/ast"
	"go/token"
)

type Detector interface {
	Detect(node ast.Node, file string, fs *token.FileSet) []Issue
}
