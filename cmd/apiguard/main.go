package main

import (
	"fmt"
	"os"

	"github.com/yokkesan/apiguard/internal/analyzer"
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
		fmt.Println("watch command is not implemented yet")
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

	goAnalyzer := analyzer.NewGoAnalyzer()

	var issues []analyzer.Issue

	for _, file := range files {
		fileIssues := goAnalyzer.Analyze(file)
		issues = append(issues, fileIssues...)
	}

	analyzer.Report(issues)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  apiguard scan <target>")
	fmt.Println("  apiguard watch <target>")
}
