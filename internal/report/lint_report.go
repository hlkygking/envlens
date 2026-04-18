package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourorg/envlens/internal/lint"
)

// LintTextReport returns a human-readable report of lint issues.
func LintTextReport(issues []lint.Issue) string {
	if len(issues) == 0 {
		return "No lint issues found.\n"
	}

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Severity != issues[j].Severity {
			return issues[i].Severity < issues[j].Severity
		}
		return issues[i].Key < issues[j].Key
	})

	var sb strings.Builder
	sb.WriteString("Lint Issues\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for _, iss := range issues {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", strings.ToUpper(iss.Severity), iss.Key, iss.Message))
	}

	warns, errors := lint.Summary(issues)
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	sb.WriteString(fmt.Sprintf("Summary: %d error(s), %d warning(s)\n", errors, warns))
	return sb.String()
}

// LintJSONReport returns a JSON-encoded report of lint issues.
func LintJSONReport(issues []lint.Issue) (string, error) {
	type payload struct {
		Issues []lint.Issue `json:"issues"`
		Errors int          `json:"errors"`
		Warns  int          `json:"warns"`
	}

	warns, errors := lint.Summary(issues)
	p := payload{
		Issues: issues,
		Errors: errors,
		Warns:  warns,
	}
	if p.Issues == nil {
		p.Issues = []lint.Issue{}
	}

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
