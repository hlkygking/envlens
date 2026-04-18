package interpolate

import (
	"fmt"
	"regexp"
	"strings"
)

// Result holds the interpolation result for a single key.
type Result struct {
	Key      string
	Original string
	Resolved string
	Refs     []string
	Missing  []string
}

// Summary holds aggregate stats.
type Summary struct {
	Total   int
	Resolved int
	Missing  int
}

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Interpolate expands ${VAR} references in env values using the provided env map.
// Self-references are not expanded.
func Interpolate(env map[string]string) []Result {
	results := make([]Result, 0, len(env))
	for key, val := range env {
		refs := extractRefs(val)
		missing := []string{}
		resolved := refPattern.ReplaceAllStringFunc(val, func(match string) string {
			ref := match[2 : len(match)-1]
			if ref == key {
				return match
			}
			if v, ok := env[ref]; ok {
				return v
			}
			missing = append(missing, ref)
			return match
		})
		results = append(results, Result{
			Key:      key,
			Original: val,
			Resolved: resolved,
			Refs:     refs,
			Missing:  missing,
		})
	}
	return results
}

// GetSummary returns aggregate stats for a set of results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if len(r.Missing) > 0 {
			s.Missing++
		} else if strings.Contains(r.Original, "${") {
			s.Resolved++
		}
	}
	return s
}

func extractRefs(val string) []string {
	matches := refPattern.FindAllStringSubmatch(val, -1)
	refs := make([]string, 0, len(matches))
	for _, m := range matches {
		refs = append(refs, m[1])
	}
	return refs
}

// Validate checks for circular references (A->B->A).
func Validate(env map[string]string) []string {
	var cycles []string
	for key, val := range env {
		refs := extractRefs(val)
		for _, ref := range refs {
			if refVal, ok := env[ref]; ok {
				if strings.Contains(refVal, fmt.Sprintf("${%s}", key)) {
					cycles = append(cycles, fmt.Sprintf("%s <-> %s", key, ref))
				}
			}
		}
	}
	return cycles
}
