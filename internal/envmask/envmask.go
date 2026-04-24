package envmask

import (
	"regexp"
	"strings"
)

// Strategy defines how masked values are rendered.
type Strategy string

const (
	StrategyRedact  Strategy = "redact"  // replace with [REDACTED]
	StrategyPartial Strategy = "partial" // show first 2 chars + ***
	StrategyHash    Strategy = "hash"    // show length hint: ***N
)

// Result holds the original and masked value for a single key.
type Result struct {
	Key      string
	Original string
	Masked   string
	WasMasked bool
}

var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|secret|token|api_?key|auth|private_?key|credential)`),
}

// IsSensitive returns true if the key matches any sensitive pattern.
func IsSensitive(key string, extra []*regexp.Regexp) bool {
	patterns := append(defaultSensitivePatterns, extra...)
	for _, p := range patterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

// maskValue applies the chosen strategy to a value.
func maskValue(v string, s Strategy) string {
	switch s {
	case StrategyPartial:
		if len(v) <= 2 {
			return "***"
		}
		return v[:2] + "***"
	case StrategyHash:
		return strings.Repeat("*", 3) + fmt.Sprintf("%d", len(v))
	default: // StrategyRedact
		return "[REDACTED]"
	}
}

// Apply masks sensitive keys in the provided map using the given strategy.
// Extra patterns extend the built-in sensitive key detection.
func Apply(env map[string]string, strategy Strategy, extra []*regexp.Regexp) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		r := Result{Key: k, Original: v, Masked: v, WasMasked: false}
		if IsSensitive(k, extra) {
			r.Masked = maskValue(v, strategy)
			r.WasMasked = true
		}
		results = append(results, r)
	}
	return results
}

// ToMap returns only the (possibly masked) key→value pairs.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.Key] = r.Masked
	}
	return m
}

// GetSummary returns total, masked, and plain counts.
func GetSummary(results []Result) (total, masked, plain int) {
	for _, r := range results {
		if r.WasMasked {
			masked++
		} else {
			plain++
		}
	}
	return len(results), masked, plain
}
