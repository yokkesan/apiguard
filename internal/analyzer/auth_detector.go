package analyzer

import "strings"

func DetectMissingAuth(routes []Route, file string) []Issue {
	var issues []Issue

	for _, route := range routes {
		if hasAuthMiddleware(route.Middlewares) {
			continue
		}

		issues = append(issues, Issue{
			Code:    "SEC-AUTH-001",
			Level:   "WARNING",
			Message: "Authentication middleware is not configured",
			File:    file,
			Line:    route.Line,
		})
	}

	return issues
}

func hasAuthMiddleware(middlewares []string) bool {
	keywords := []string{
		"auth",
		"jwt",
		"session",
		"authenticate",
		"authorization",
	}

	for _, middleware := range middlewares {
		name := strings.ToLower(middleware)

		for _, keyword := range keywords {
			if strings.Contains(name, keyword) {
				return true
			}
		}
	}

	return false
}
