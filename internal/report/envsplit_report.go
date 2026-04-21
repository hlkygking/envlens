package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yourusername/envlens/internal/envsplit"
)

// EnvSplitTextReport renders a human-readable report of split results grouped by bucket.
func EnvSplitTextReport(results []envsplit.Result, summary envsplit.Summary) string {
	var sb strings.Builder

	// Collect buckets in sorted order for deterministic output.
	buckets := make([]string, 0, len(summary.Buckets))
	for b := range summary.Buckets {
		buckets = append(buckets, b)
	}
	sort.Strings(buckets)

	for _, bucket := range buckets {
		fmt.Fprintf(&sb, "[%s] (%d)\n", strings.ToUpper(bucket), summary.Buckets[bucket])
		matched := envsplit.FilterByBucket(results, bucket)
		sort.Slice(matched, func(i, j int) bool { return matched[i].Key < matched[j].Key })
		for _, r := range matched {
			fmt.Fprintf(&sb, "  %s=%s\n", r.Key, r.Value)
		}
	}

	fmt.Fprintf(&sb, "\nSummary: %d keys across %d bucket(s)\n", summary.Total, len(summary.Buckets))
	return sb.String()
}

// EnvSplitJSONReport renders a JSON report of split results.
func EnvSplitJSONReport(results []envsplit.Result, summary envsplit.Summary) (string, error) {
	type entry struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Bucket string `json:"bucket"`
	}
	type payload struct {
		Entries []entry        `json:"entries"`
		Buckets map[string]int `json:"buckets"`
		Total   int            `json:"total"`
	}

	entries := make([]entry, 0, len(results))
	for _, r := range results {
		entries = append(entries, entry{Key: r.Key, Value: r.Value, Bucket: r.Bucket})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Bucket != entries[j].Bucket {
			return entries[i].Bucket < entries[j].Bucket
		}
		return entries[i].Key < entries[j].Key
	})

	p := payload{Entries: entries, Buckets: summary.Buckets, Total: summary.Total}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
