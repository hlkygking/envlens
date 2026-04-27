package envexpand

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Result holds the expansion result for a single env var.
type Result struct {
	Key      string
	Original string
	Expanded string
	Refs     []string
	Status   string // ok, unresolved, unchanged
}

// Summary holds aggregate stats for an expansion run.
type Summary struct {
	Total      int
	Expanded   int
	Unresolved int
	Unchanged  int
}

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Apply expands variable references in env values using the provided env map
// and optionally falls back to OS environment variables.
func Apply(env map[string]string, fallbackToOS bool) []Result {
	results := make([]Result, 0, len(env))

	for key, val := range env {
		refs := extractRefs(val)
		if len(refs) == 0 {
			results = append(results, Result{
				Key:      key,
				Original: val,
				Expanded: val,
				Refs:     nil,
				Status:   "unchanged",
			})
			continue
		}

		expanded, unresolved := expandValue(val, env, fallbackToOS)
		status := "ok"
		if unresolved > 0 {
			status = "unresolved"
		}

		results = append(results, Result{
			Key:      key,
			Original: val,
			Expanded: expanded,
			Refs:     refs,
			Status:   status,
		})
	}

	return results
}

// ToMap returns a key→expanded-value map from results with status "ok" or "unchanged".
func ToMap(results []Result) map[string]string {
	out := make(map[string]string, len(results))
	for _, r := range results {
		if r.Status == "ok" || r.Status == "unchanged" {
			out[r.Key] = r.Expanded
		}
	}
	return out
}

// GetSummary returns aggregate counts for the expansion results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case "ok":
			s.Expanded++
		case "unresolved":
			s.Unresolved++
		case "unchanged":
			s.Unchanged++
		}
	}
	return s
}

func expandValue(val string, env map[string]string, fallbackToOS bool) (string, int) {
	unresolved := 0
	result := refPattern.ReplaceAllStringFunc(val, func(match string) string {
		name := match
		if strings.HasPrefix(match, "${") {
			name = match[2 : len(match)-1]
		} else {
			name = match[1:]
		}
		if v, ok := env[name]; ok {
			return v
		}
		if fallbackToOS {
			if v := os.Getenv(name); v != "" {
				return v
			}
		}
		unresolved++
		return fmt.Sprintf("${%s}", name)
	})
	return result, unresolved
}

func extractRefs(val string) []string {
	matches := refPattern.FindAllStringSubmatch(val, -1)
	seen := map[string]bool{}
	var refs []string
	for _, m := range matches {
		name := m[1]
		if name == "" {
			name = m[2]
		}
		if !seen[name] {
			seen[name] = true
			refs = append(refs, name)
		}
	}
	return refs
}
