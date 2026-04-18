package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/wryfi/envlens/internal/transform"
)

// TransformTextReport returns a human-readable report of transform results.
func TransformTextReport(results []transform.Result) string {
	var sb strings.Builder
	transformed, unchanged := transform.Summary(results)
	fmt.Fprintf(&sb, "Transform Summary: %d changed, %d unchanged\n\n", transformed, unchanged)

	sorted := make([]transform.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Key < sorted[j].Key })

	for _, r := range sorted {
		if r.Value != r.Original {
			fmt.Fprintf(&sb, "  [changed] %s: %q -> %q (ops: %v)\n", r.Key, r.Original, r.Value, r.Applied)
		} else {
			fmt.Fprintf(&sb, "  [unchanged] %s: %q\n", r.Key, r.Value)
		}
	}
	return sb.String()
}

// TransformJSONReport returns a JSON report of transform results.
func TransformJSONReport(results []transform.Result) (string, error) {
	transformed, unchanged := transform.Summary(results)
	payload := map[string]interface{}{
		"summary": map[string]int{
			"transformed": transformed,
			"unchanged":   unchanged,
		},
		"results": results,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
