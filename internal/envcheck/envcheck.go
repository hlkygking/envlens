package envcheck

import (
	"fmt"
	"regexp"
	"strings"
)

// Status represents the result of a single key check.
type Status string

const (
	StatusOK      Status = "ok"
	StatusMissing Status = "missing"
	StatusInvalid Status = "invalid"
)

// Rule defines an expected environment variable and optional constraints.
type Rule struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Result holds the outcome of checking a single rule.
type Result struct {
	Key     string
	Status  Status
	Message string
	Value   string
}

// Summary holds aggregate counts.
type Summary struct {
	Total   int
	OK      int
	Missing int
	Invalid int
}

// Apply checks env against the provided rules and returns results.
func Apply(env map[string]string, rules []Rule) []Result {
	results := make([]Result, 0, len(rules))
	for _, rule := range rules {
		val, exists := env[rule.Key]
		if !exists || val == "" {
			if rule.Required {
				results = append(results, Result{
					Key:     rule.Key,
					Status:  StatusMissing,
					Message: "required key is missing or empty",
				})
			} else {
				results = append(results, Result{
					Key:     rule.Key,
					Status:  StatusOK,
					Message: "optional key not set",
				})
			}
			continue
		}
		if rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil || !re.MatchString(val) {
				msg := fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern)
				if err != nil {
					msg = fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err)
				}
				results = append(results, Result{Key: rule.Key, Status: StatusInvalid, Message: msg, Value: val})
				continue
			}
		}
		results = append(results, Result{Key: rule.Key, Status: StatusOK, Value: val})
	}
	return results
}

// GetSummary returns aggregate counts from results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusOK:
			s.OK++
		case StatusMissing:
			s.Missing++
		case StatusInvalid:
			s.Invalid++
		}
	}
	return s
}

// ParseRules parses lines of the form "KEY", "KEY:required", "KEY:required:pattern".
func ParseRules(lines []string) []Rule {
	rules := make([]Rule, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		rule := Rule{Key: parts[0]}
		if len(parts) >= 2 {
			rule.Required = strings.EqualFold(parts[1], "required")
		}
		if len(parts) == 3 {
			rule.Pattern = parts[2]
		}
		rules = append(rules, rule)
	}
	return rules
}
