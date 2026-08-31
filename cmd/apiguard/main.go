package main

import (
	"fmt"
	"os"

	"github.com/yokkesan/apiguard/internal/analyzer"
	"github.com/yokkesan/apiguard/internal/watcher"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "scan":
		runScan(os.Args[2:])
	case "watch":
		runWatch(os.Args[2:])
	default:
		fmt.Printf("unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runScan(args []string) {
	if len(args) == 0 {
		fmt.Println("scan target is required")
		fmt.Println("example: apiguard scan ./...")
		os.Exit(1)
	}

	target := args[0]

	files, err := analyzer.FindGoFiles(target)
	if err != nil {
		fmt.Printf("failed to find Go files: %v\n", err)
		os.Exit(1)
	}

	goAnalyzer := newGoAnalyzer()

	issues := goAnalyzer.AnalyzeFiles(files)

	analyzer.Report(issues)
}

func runWatch(args []string) {
	if len(args) == 0 {
		fmt.Println("watch target is required")
		fmt.Println("example: apiguard watch .")
		os.Exit(1)
	}

	target := args[0]

	files, err := analyzer.FindGoFiles(target)
	if err != nil {
		fmt.Printf("failed to find Go files: %v\n", err)
		os.Exit(1)
	}

	goAnalyzer := newGoAnalyzer()

	loadIssues := goAnalyzer.LoadFiles(files)
	analyzer.Report(loadIssues)

	fileWatcher, err := watcher.New()
	if err != nil {
		fmt.Printf("failed to create watcher: %v\n", err)
		os.Exit(1)
	}
	defer fileWatcher.Close()

	if err := fileWatcher.AddRecursive(target); err != nil {
		fmt.Printf("failed to watch target: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("watching %s\n", target)

	err = fileWatcher.Run(func(changedFile string) {
		fmt.Printf("changed: %s\n", changedFile)

		issues := goAnalyzer.AnalyzeChangedFile(changedFile)

		analyzer.Report(issues)
	})

	if err != nil {
		fmt.Printf("watch error: %v\n", err)
		os.Exit(1)
	}
}

func newGoAnalyzer() *analyzer.GoAnalyzer {
	return analyzer.NewGoAnalyzer(
		analyzer.SecretDetector{},
		analyzer.SQLDetector{},
		analyzer.AuthDetector{},
		analyzer.ValidationDetector{},
	)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  apiguard scan <target>")
	fmt.Println("  apiguard watch <target>")
}
