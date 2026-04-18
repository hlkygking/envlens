package audit

import (
	"strings"

	"github.com/envlens/envlens/internal/diff"
)

// Severity represents the risk level of an audit finding.
type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

// Finding represents a single audit observation.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

// sensitivePatterns are substrings that flag a key as sensitive.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
}

// isSensitiveKey returns true when the key name suggests sensitive data.
func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Audit inspects a slice of diff entries and returns audit findings.
func Audit(entries []diff.Entry) []Finding {
	var findings []Finding

	for _, e := range entries {
		switch e.Status {
		case diff.Added:
			if isSensitiveKey(e.Key) {
				findings = append(findings, Finding{
					Key:      e.Key,
					Message:  "Sensitive key added in target environment",
					Severity: SeverityHigh,
				})
			}
		case diff.Removed:
			if isSensitiveKey(e.Key) {
				findings = append(findings, Finding{
					Key:      e.Key,
					Message:  "Sensitive key removed from target environment",
					Severity: SeverityHigh,
				})
			}
		case diff.Modified:
			if isSensitiveKey(e.Key) {
				findings = append(findings, Finding{
					Key:      e.Key,
					Message:  "Sensitive key value changed between environments",
					Severity: SeverityMedium,
				})
			}
			if e.OldValue == "" {
				findings = append(findings, Finding{
					Key:      e.Key,
					Message:  "Key transitioned from empty value to non-empty",
					Severity: SeverityLow,
				})
			}
		}
	}

	return findings
}
