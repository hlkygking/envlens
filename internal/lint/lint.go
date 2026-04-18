package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Issue represents a linting problem found in an env map.
type Issue struct {
	Key      string
	Message  string
	Severity string // "warn" or "error"
}

var (
	dupSpaceRe  = regexp.MustCompile(`\s{2,}`)
	lowerKeyRe  = regexp.MustCompile(`[a-z]`)
	invalidKeyRe = regexp.MustCompile(`[^A-Z0-9_]`)
)

// Lint checks an env map for common style and correctness issues.
func Lint(env map[string]string) []Issue {
	var issues []Issue

	for k, v := range env {
		// Key must be uppercase
		if lowerKeyRe.MatchString(k) {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "key contains lowercase letters; prefer UPPER_SNAKE_CASE",
				Severity: "warn",
			})
		}

		// Key must not contain invalid characters
		if invalidKeyRe.MatchString(k) {
			issues = append(issues, Issue{
				Key:      k,
				Message:  fmt.Sprintf("key %q contains invalid characters", k),
				Severity: "error",
			})
		}

		// Value should not have leading/trailing whitespace
		if strings.TrimSpace(v) != v {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value has leading or trailing whitespace",
				Severity: "warn",
			})
		}

		// Value should not contain double spaces
		if dupSpaceRe.MatchString(v) {
			issues = append(issues, Issue{
				Key:      k,
				Message:  "value contains consecutive whitespace",
				Severity: "warn",
			})
		}
	}

	return issues
}

// Summary returns counts of warn and error issues.
func Summary(issues []Issue) (warns, errors int) {
	for _, i := range issues {
		switch i.Severity {
		case "warn":
			warns++
		case "error":
			errors++
		}
	}
	return
}
