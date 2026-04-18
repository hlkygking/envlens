package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envlens/internal/compare"
)

// MultiTextReport writes a human-readable multi-file comparison report.
func MultiTextReport(w io.Writer, res *compare.MultiResult) {
	fmt.Fprintf(w, "Base: %s\n", res.Base)
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 50))
	for _, pr := range res.Targets {
		fmt.Fprintf(w, "\nTarget: %s\n", pr.Target)
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", 40))
		if len(pr.Entries) == 0 {
			fmt.Fprintln(w, "  No differences.")
			continue
		}
		for _, e := range pr.Entries {
			switch e.Status {
			case "added":
				fmt.Fprintf(w, "  + %s=%s\n", e.Key, e.NewValue)
			case "removed":
				fmt.Fprintf(w, "  - %s=%s\n", e.Key, e.OldValue)
			case "modified":
				fmt.Fprintf(w, "  ~ %s: %s -> %s\n", e.Key, e.OldValue, e.NewValue)
			}
		}
	}
}

// MultiJSONReport writes a JSON multi-file comparison report.
func MultiJSONReport(w io.Writer, res *compare.MultiResult) error {
	type pairJSON struct {
		Target  string      `json:"target"`
		Entries interface{} `json:"entries"`
	}
	type output struct {
		Base    string      `json:"base"`
		Targets []pairJSON  `json:"targets"`
	}
	out := output{Base: res.Base}
	for _, pr := range res.Targets {
		out.Targets = append(out.Targets, pairJSON{
			Target:  pr.Target,
			Entries: pr.Entries,
		})
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
