package envnorm

import (
	"strings"
	"unicode"
)

// NormMode defines how keys should be normalized.
type NormMode string

const (
	ModeUppercase  NormMode = "uppercase"
	ModeLowercase  NormMode = "lowercase"
	ModeSnakeCase  NormMode = "snake_case"
	ModeStripSpace NormMode = "strip_space"
)

// Result holds the original and normalized key/value.
type Result struct {
	OriginalKey string
	NormalizedKey string
	Value       string
	Changed     bool
}

// Summary holds aggregate stats.
type Summary struct {
	Total   int
	Changed int
	Unchanged int
}

// Apply normalizes env keys according to the given modes.
func Apply(env map[string]string, modes []NormMode) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		nk := normalize(k, modes)
		results = append(results, Result{
			OriginalKey:   k,
			NormalizedKey: nk,
			Value:         v,
			Changed:       nk != k,
		})
	}
	return results
}

// ToMap returns a map of normalized keys to values.
func ToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		m[r.NormalizedKey] = r.Value
	}
	return m
}

// GetSummary returns counts of changed/unchanged results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Changed {
			s.Changed++
		} else {
			s.Unchanged++
		}
	}
	return s
}

func normalize(key string, modes []NormMode) string {
	for _, m := range modes {
		switch m {
		case ModeUppercase:
			key = strings.ToUpper(key)
		case ModeLowercase:
			key = strings.ToLower(key)
		case ModeSnakeCase:
			key = toSnakeCase(key)
		case ModeStripSpace:
			key = strings.ReplaceAll(key, " ", "_")
		}
	}
	return key
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteRune('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}
