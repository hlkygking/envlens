package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable.
type Rule struct {
	Key     string
	Pattern string // optional regex pattern for value
	Required bool
}

// Violation represents a failed validation.
type Violation struct {
	Key     string
	Message string
}

var validKeyPattern = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// ValidateKeys checks that all keys in the env map follow naming conventions.
func ValidateKeys(env map[string]string) []Violation {
	var violations []Violation
	for k := range env {
		if !validKeyPattern.MatchString(k) {
			violations = append(violations, Violation{
				Key:     k,
				Message: fmt.Sprintf("key %q does not match naming convention (uppercase, underscores only)", k),
			})
		}
	}
	return violations
}

// ValidateRules checks that required keys exist and values match patterns.
func ValidateRules(env map[string]string, rules []Rule) []Violation {
	var violations []Violation
	for _, rule := range rules {
		val, exists := env[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: fmt.Sprintf("required key %q is missing", rule.Key),
			})
			continue
		}
		if exists && rule.Pattern != "" {
			re, err := regexp.Compile(rule.Pattern)
			if err != nil {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("invalid pattern for key %q: %v", rule.Key, err),
				})
				continue
			}
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value for %q does not match pattern %q", rule.Key, rule.Pattern),
				})
			}
		}
	}
	return violations
}

// Summary returns a human-readable summary of violations.
func Summary(violations []Violation) string {
	if len(violations) == 0 {
		return "validation passed: no violations found"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d violation(s) found:\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  [%s] %s\n", v.Key, v.Message)
	}
	return sb.String()
}
