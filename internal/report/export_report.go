package report

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ExportSummary holds metadata about an export operation.
type ExportSummary struct {
	Format string `json:"format"`
	Path   string `json:"path"`
	Keys   int    `json:"keys"`
}

// ExportTextReport returns a human-readable export summary.
func ExportTextReport(s ExportSummary) string {
	var sb strings.Builder
	sb.WriteString("=== Export Summary ===\n")
	fmt.Fprintf(&sb, "Format : %s\n", s.Format)
	fmt.Fprintf(&sb, "Output : %s\n", s.Path)
	fmt.Fprintf(&sb, "Keys   : %d\n", s.Keys)
	return sb.String()
}

// ExportJSONReport returns a JSON export summary.
func ExportJSONReport(s ExportSummary) (string, error) {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}
