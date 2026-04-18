package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/your-org/envlens/internal/resolve"
)

// ResolveTextReport returns a human-readable report of resolution results.
func ResolveTextReport(results []resolve.Result) string {
	var sb strings.Builder

	sorted := make([]resolve.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	sb.WriteString("=== Resolve Report ===\n")
	for _, r := range sorted {
		status := "OK"
		if !r.Resolved {
			status = "MISSING"
		}
		fmt.Fprintf(&sb, "  [%-7s] %-30s = %s  (source: %s)\n",
			status, r.Key, r.Value, r.Source)
	}

	sb.WriteString("\n")
	sb.WriteString(resolve.Summary(results))
	sb.WriteString("\n")
	return sb.String()
}

// resolveJSON is the JSON shape for a single result.
type resolveJSON struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Source   string `json:"source"`
	Resolved bool   `json:"resolved"`
}

// ResolveJSONReport returns a JSON-encoded report of resolution results.
func ResolveJSONReport(results []resolve.Result) (string, error) {
	items := make([]resolveJSON, 0, len(results))
	for _, r := range results {
		items = append(items, resolveJSON{
			Key:      r.Key,
			Value:    r.Value,
			Source:   r.Source,
			Resolved: r.Resolved,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	out := map[string]interface{}{
		"summary": resolve.Summary(results),
		"results": items,
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
