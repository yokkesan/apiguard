package analyzer

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func FindGoFiles(target string) ([]string, error) {
	root := normalizeScanTarget(target)

	var files []string

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if shouldSkipDir(entry.Name()) {
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func normalizeScanTarget(target string) string {
	if target == "./..." {
		return "."
	}

	if strings.HasSuffix(target, "/...") {
		return strings.TrimSuffix(target, "/...")
	}

	return target
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "vendor":
		return true
	default:
		return false
	}
}
