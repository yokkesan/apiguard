package analyzer

import "go/ast"

type Detector interface {
	Detect(
		node ast.Node,
		file string,
		context *AnalysisContext,
	) []Issue
}
