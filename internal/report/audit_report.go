package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/envlens/envlens/internal/audit"
)

// AuditTextReport writes a human-readable audit findings report to w.
func AuditTextReport(w io.Writer, findings []audit.Finding) {
	if len(findings) == 0 {
		fmt.Fprintln(w, "No audit findings.")
		return
	}

	counts := map[audit.Severity]int{}
	for _, f := range findings {
		counts[f.Severity]++
	}

	fmt.Fprintf(w, "Audit Summary: %d finding(s) — HIGH:%d MEDIUM:%d LOW:%d\n",
		len(findings), counts[audit.SeverityHigh], counts[audit.SeverityMedium], counts[audit.SeverityLow])
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, sev := range []audit.Severity{audit.SeverityHigh, audit.SeverityMedium, audit.SeverityLow} {
		for _, f := range findings {
			if f.Severity != sev {
				continue
			}
			fmt.Fprintf(w, "[%s] %s: %s\n", f.Severity, f.Key, f.Message)
		}
	}
}

// AuditJSONReport writes audit findings as JSON to w.
func AuditJSONReport(w io.Writer, findings []audit.Finding) error {
	type jsonFinding struct {
		Key      string `json:"key"`
		Severity string `json:"severity"`
		Message  string `json:"message"`
	}
	out := make([]jsonFinding, len(findings))
	for i, f := range findings {
		out[i] = jsonFinding{Key: f.Key, Severity: string(f.Severity), Message: f.Message}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
