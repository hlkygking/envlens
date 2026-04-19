package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yourorg/envlens/internal/envcast"
)

// EnvCastTextReport renders cast results as human-readable text.
func EnvCastTextReport(results []envcast.Result) string {
	var sb strings.Builder
	ok, failed := envcast.Summary(results)
	fmt.Fprintf(&sb, "=== Cast Results ===\n")
	fmt.Fprintf(&sb, "Total: %d | OK: %d | Failed: %d\n\n", len(results), ok, failed)

	var failures []envcast.Result
	for _, r := range results {
		if !r.OK {
			failures = append(failures, r)
		}
	}

	if len(failures) > 0 {
		fmt.Fprintf(&sb, "[FAILED]\n")
		for _, r := range failures {
			fmt.Fprintf(&sb, "  %-24s  type=%-8s  raw=%q  error=%s\n",
				r.Key, r.CastType, r.RawValue, r.Error)
		}
		sb.WriteString("\n")
	}

	fmt.Fprintf(&sb, "[OK]\n")
	for _, r := range results {
		if r.OK {
			fmt.Fprintf(&sb, "  %-24s  type=%-8s  value=%v\n",
				r.Key, r.CastType, r.CastValue)
		}
	}
	return sb.String()
}

// EnvCastJSONReport renders cast results as JSON.
func EnvCastJSONReport(results []envcast.Result) (string, error) {
	type entry struct {
		Key       string      `json:"key"`
		RawValue  string      `json:"raw_value"`
		CastType  string      `json:"cast_type"`
		CastValue interface{} `json:"cast_value,omitempty"`
		OK        bool        `json:"ok"`
		Error     string      `json:"error,omitempty"`
	}
	ok, failed := envcast.Summary(results)
	payload := struct {
		Total   int     `json:"total"`
		OK      int     `json:"ok"`
		Failed  int     `json:"failed"`
		Entries []entry `json:"entries"`
	}{
		Total:  len(results),
		OK:     ok,
		Failed: failed,
	}
	for _, r := range results {
		payload.Entries = append(payload.Entries, entry{
			Key:       r.Key,
			RawValue:  r.RawValue,
			CastType:  string(r.CastType),
			CastValue: r.CastValue,
			OK:        r.OK,
			Error:     r.Error,
		})
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
