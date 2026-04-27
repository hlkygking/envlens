package envstats

import (
	"sort"
	"strings"
	"unicode"
)

// Result holds statistics for a single env entry.
type Result struct {
	Key        string
	Value      string
	Length     int
	WordCount  int
	HasDigits  bool
	HasUpper   bool
	HasLower   bool
	HasSpecial bool
	IsEmpty    bool
}

// Summary holds aggregate statistics across all entries.
type Summary struct {
	Total       int
	Empty       int
	AvgLength   float64
	MaxLength   int
	MinLength   int
	LongestKey  string
	ShortestKey string
}

// Apply computes per-entry statistics for the given env map.
func Apply(env map[string]string) []Result {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	results := make([]Result, 0, len(keys))
	for _, k := range keys {
		v := env[k]
		r := Result{
			Key:       k,
			Value:     v,
			Length:    len(v),
			WordCount: len(strings.Fields(v)),
			IsEmpty:   len(v) == 0,
		}
		for _, ch := range v {
			switch {
			case unicode.IsDigit(ch):
				r.HasDigits = true
			case unicode.IsUpper(ch):
				r.HasUpper = true
			case unicode.IsLower(ch):
				r.HasLower = true
			case !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && !unicode.IsSpace(ch):
				r.HasSpecial = true
			}
		}
		results = append(results, r)
	}
	return results
}

// GetSummary computes aggregate statistics from a slice of Results.
func GetSummary(results []Result) Summary {
	if len(results) == 0 {
		return Summary{}
	}
	s := Summary{
		Total:     len(results),
		MinLength: results[0].Length,
	}
	var totalLen int
	for _, r := range results {
		totalLen += r.Length
		if r.IsEmpty {
			s.Empty++
		}
		if r.Length > s.MaxLength {
			s.MaxLength = r.Length
			s.LongestKey = r.Key
		}
		if r.Length < s.MinLength {
			s.MinLength = r.Length
			s.ShortestKey = r.Key
		}
	}
	s.AvgLength = float64(totalLen) / float64(len(results))
	return s
}
