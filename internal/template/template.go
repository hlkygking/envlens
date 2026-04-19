package template

import (
	"fmt"
	"regexp"
	"strings"
)

// Result holds the outcome of applying a template to an env map.
type Result struct {
	Key      string
	Template string
	Rendered string
	Missing  []string
	OK       bool
}

var placeholderRe = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)

// Apply renders each template string against the provided env map.
// Templates use {{KEY}} syntax.
func Apply(templates map[string]string, env map[string]string) []Result {
	results := make([]Result, 0, len(templates))
	for key, tmpl := range templates {
		r := render(key, tmpl, env)
		results = append(results, r)
	}
	return results
}

func render(key, tmpl string, env map[string]string) Result {
	var missing []string
	output := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		sub := placeholderRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		ref := strings.TrimSpace(sub[1])
		if val, ok := env[ref]; ok {
			return val
		}
		missing = append(missing, ref)
		return fmt.Sprintf("<missing:%s>", ref)
	})
	return Result{
		Key:      key,
		Template: tmpl,
		Rendered: output,
		Missing:  missing,
		OK:       len(missing) == 0,
	}
}

// ToMap returns a map of key -> rendered value for successful results.
func ToMap(results []Result) map[string]string {
	out := make(map[string]string)
	for _, r := range results {
		if r.OK {
			out[r.Key] = r.Rendered
		}
	}
	return out
}

// GetSummary returns counts of ok and failed renders.
func GetSummary(results []Result) (ok, failed int) {
	for _, r := range results {
		if r.OK {
			ok++
		} else {
			failed++
		}
	}
	return
}
