package analyzer

type Analyzer interface {
	Analyze(path string) []Issue
}
