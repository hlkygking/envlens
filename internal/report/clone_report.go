package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/your/envlens/internal/clone"
)

// CloneTextReport returns a human-readable summary of a clone operation.
func CloneTextReport(r clone.Result) string {
	var sb strings.Builder
	sb.WriteString("=== Clone Report ===\n")
	fmt.Fprintf(&sb, "  Source : %s\n", r.SourceFile)
	fmt.Fprintf(&sb, "  Dest   : %s\n", r.DestFile)
	fmt.Fprintf(&sb, "  Keys   : %d\n", r.Keys)
	if r.Redacted {
		sb.WriteString("  Redact : sensitive values cleared\n")
	} else {
		sb.WriteString("  Redact : none\n")
	}
	return sb.String()
}

// CloneJSONReport returns a JSON representation of a clone result.
func CloneJSONReport(r clone.Result) (string, error) {
	type payload struct {
		Source   string `json:"source"`
		Dest     string `json:"dest"`
		Keys     int    `json:"keys"`
		Redacted bool   `json:"redacted"`
	}
	p := payload{
		Source:   r.SourceFile,
		Dest:     r.DestFile,
		Keys:     r.Keys,
		Redacted: r.Redacted,
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
