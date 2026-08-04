package analyzer

import "path/filepath"

type Language string

const (
	LanguageGo         Language = "go"
	LanguagePHP        Language = "php"
	LanguageTypeScript Language = "typescript"
	LanguageDart       Language = "dart"
	LanguageUnknown    Language = "unknown"
)

func DetectLanguage(path string) Language {
	ext := filepath.Ext(path)

	switch ext {
	case ".go":
		return LanguageGo
	case ".php":
		return LanguagePHP
	case ".ts", ".tsx":
		return LanguageTypeScript
	case ".dart":
		return LanguageDart
	default:
		return LanguageUnknown
	}
}
