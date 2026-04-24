package envdoc

import (
	"bufio"
	"os"
	"strings"
)

// Entry holds a documented environment variable.
type Entry struct {
	Key         string
	Description string
	Example     string
	Required    bool
	Default     string
}

// Result holds the outcome of annotating a key.
type Result struct {
	Key         string
	Description string
	Example     string
	Required    bool
	Default     string
	Found       bool
}

// Summary holds aggregate counts.
type Summary struct {
	Total       int
	Documented  int
	Undocumented int
}

// Apply matches env keys against a documentation map and returns results.
func Apply(env map[string]string, docs []Entry) []Result {
	docMap := make(map[string]Entry, len(docs))
	for _, d := range docs {
		docMap[d.Key] = d
	}

	results := make([]Result, 0, len(env))
	for k, _ := range env {
		r := Result{Key: k}
		if d, ok := docMap[k]; ok {
			r.Description = d.Description
			r.Example = d.Example
			r.Required = d.Required
			r.Default = d.Default
			r.Found = true
		}
		results = append(results, r)
	}
	return results
}

// GetSummary returns counts from results.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		if r.Found {
			s.Documented++
		} else {
			s.Undocumented++
		}
	}
	return s
}

// ParseDocFile reads a simple TSV/colon-delimited doc file:
// KEY:description:example:required:default
func ParseDocFile(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 5)
		if len(parts) < 1 {
			continue
		}
		e := Entry{Key: strings.TrimSpace(parts[0])}
		if len(parts) > 1 {
			e.Description = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 {
			e.Example = strings.TrimSpace(parts[2])
		}
		if len(parts) > 3 {
			e.Required = strings.TrimSpace(parts[3]) == "true"
		}
		if len(parts) > 4 {
			e.Default = strings.TrimSpace(parts[4])
		}
		entries = append(entries, e)
	}
	return entries, scanner.Err()
}
