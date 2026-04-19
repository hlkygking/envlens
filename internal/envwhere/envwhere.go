package envwhere

import (
	"fmt"
	"regexp"
	"strings"
)

// MatchMode controls how key matching is performed.
type MatchMode string

const (
	MatchExact  MatchMode = "exact"
	MatchPrefix MatchMode = "prefix"
	MatchSuffix MatchMode = "suffix"
	MatchRegex  MatchMode = "regex"
)

// Result represents a single key lookup result.
type Result struct {
	Key     string
	Value   string
	Matched bool
	Mode    MatchMode
	Error   string
}

// Options configures the where query.
type Options struct {
	Mode       MatchMode
	CaseFold   bool
}

// Apply searches env for keys matching the given query using opts.
func Apply(env map[string]string, query string, opts Options) []Result {
	var results []Result

	if opts.Mode == "" {
		opts.Mode = MatchExact
	}

	var re *regexp.Regexp
	if opts.Mode == MatchRegex {
		flags := ""
		if opts.CaseFold {
			flags = "(?i)"
		}
		var err error
		re, err = regexp.Compile(flags + query)
		if err != nil {
			return []Result{{Key: query, Matched: false, Mode: opts.Mode, Error: fmt.Sprintf("invalid regex: %v", err)}}
		}
	}

	for k, v := range env {
		matched := false
		compareKey := k
		compareQuery := query
		if opts.CaseFold && opts.Mode != MatchRegex {
			compareKey = strings.ToLower(k)
			compareQuery = strings.ToLower(query)
		}
		switch opts.Mode {
		case MatchExact:
			matched = compareKey == compareQuery
		case MatchPrefix:
			matched = strings.HasPrefix(compareKey, compareQuery)
		case MatchSuffix:
			matched = strings.HasSuffix(compareKey, compareQuery)
		case MatchRegex:
			matched = re.MatchString(k)
		}
		if matched {
			results = append(results, Result{Key: k, Value: v, Matched: true, Mode: opts.Mode})
		}
	}
	return results
}

// GetSummary returns counts of matched and total keys.
func GetSummary(results []Result) (matched, total int) {
	total = len(results)
	for _, r := range results {
		if r.Matched {
			matched++
		}
	}
	return
}
