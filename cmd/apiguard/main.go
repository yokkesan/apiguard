package main

import (
	"fmt"
	"os"

	"apiguard/internal/analyzer"
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

	goAnalyzer := analyzer.NewGoAnalyzer()

	issues := goAnalyzer.Analyze(target)

	analyzer.Report(issues)
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  apiguard scan <target>")
	fmt.Println("  apiguard watch <target>")
}
