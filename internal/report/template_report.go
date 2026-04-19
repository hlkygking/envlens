package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envlens/internal/template"
)

// TemplateTextReport returns a human-readable report of template render results.
func TemplateTextReport(results []template.Result) string {
	var sb strings.Builder
	ok, failed := countTemplate(results)
	fmt.Fprintf(&sb, "Template Render Report\n")
	fmt.Fprintf(&sb, "======================\n")
	fmt.Fprintf(&sb, "Total: %d | OK: %d | Failed: %d\n\n", len(results), ok, failed)

	if ok > 0 {
		sb.WriteString("[OK]\n")
		for _, r := range results {
			if r.OK {
				fmt.Fprintf(&sb, "  %-20s = %s\n", r.Key, r.Rendered)
			}
		}
		sb.WriteString("\n")
	}

	if failed > 0 {
		sb.WriteString("[FAILED]\n")
		for _, r := range results {
			if !r.OK {
				fmt.Fprintf(&sb, "  %-20s template: %q\n", r.Key, r.Template)
				fmt.Fprintf(&sb, "  %-20s missing:  %s\n", "", strings.Join(r.Missing, ", "))
			}
		}
	}
	return sb.String()
}

// TemplateJSONReport returns a JSON-encoded report of template render results.
func TemplateJSONReport(results []template.Result) (string, error) {
	ok, failed := countTemplate(results)
	payload := map[string]interface{}{
		"total":   len(results),
		"ok":      ok,
		"failed":  failed,
		"results": results,
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func countTemplate(results []template.Result) (ok, failed int) {
	for _, r := range results {
		if r.OK {
			ok++
		} else {
			failed++
		}
	}
	return
}
