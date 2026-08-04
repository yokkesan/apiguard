package analyzer

import "fmt"

const (
	red    = "\033[31m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	reset  = "\033[0m"
)

func Report(issue Issue) {
	color := ""

	switch issue.Level {
	case "ERROR":
		color = red
	case "WARNING":
		color = yellow
	case "INFO":
		color = blue
	}

	fmt.Printf("%s[%s] %s%s\n",
		color,
		issue.Level,
		issue.Code,
		reset,
	)

	fmt.Println()
	fmt.Println(issue.Message)
	fmt.Println()
	fmt.Printf("%s:%d\n", issue.File, issue.Line)
}
