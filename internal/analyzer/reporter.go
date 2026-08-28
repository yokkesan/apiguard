package analyzer

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func Report(issues []Issue) {
	for _, issue := range issues {
		color := colorReset

		switch issue.Level {
		case "ERROR":
			color = colorRed
		case "WARNING":
			color = colorYellow
		case "INFO":
			color = colorBlue
		}

		fmt.Printf(
			"%s[%s] %s%s\n",
			color,
			issue.Level,
			issue.Code,
			colorReset,
		)

		fmt.Println(issue.Message)

		if issue.File != "" {
			if issue.Line > 0 {
				fmt.Printf("%s:%d\n", issue.File, issue.Line)
			} else {
				fmt.Println(issue.File)
			}
		}

		fmt.Println()
	}
}
