package analyzer

type Issue struct {
	File    string
	Message string
	Level   string
}

type Analyzer interface {
	Analyze(path string) []Issue
}

func Analyze(path string) []Issue {
	language := DetectLanguage(path)

	switch language {
	case LanguageGo:
		return GoAnalyzer{}.Analyze(path)

	default:
		return []Issue{
			{
				File:    path,
				Message: "unsupported language",
				Level:   "INFO",
			},
		}
	}
}
