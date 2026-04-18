package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/diff"
)

func sampleEntries() []diff.Entry {
	return []diff.Entry{
		{Key: "NEW_KEY", Status: diff.Added, OldValue: "", NewValue: "hello"},
		{Key: "OLD_KEY", Status: diff.Removed, OldValue: "bye", NewValue: ""},
		{Key: "CHANGED", Status: diff.Modified, OldValue: "v1", NewValue: "v2"},
	}
}

func TestTextReport_ContainsSections(t *testing.T) {
	var buf bytes.Buffer
	TextReport(&buf, sampleEntries())
	out := buf.String()

	for _, want := range []string{"Added:", "Removed:", "Modified:", "Summary:"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output", want)
		}
	}
}

func TestTextReport_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	TextReport(&buf, nil)
	if !strings.Contains(buf.String(), "No differences found.") {
		t.Error("expected no-diff message")
	}
}

func TestTextReport_Summary(t *testing.T) {
	var buf bytes.Buffer
	TextReport(&buf, sampleEntries())
	if !strings.Contains(buf.String(), "1 added, 1 removed, 1 modified") {
		t.Error("expected summary counts")
	}
}

func TestJSONReport_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	JSONReport(&buf, sampleEntries())
	out := buf.String()

	for _, want := range []string{"NEW_KEY", "OLD_KEY", "CHANGED", "added", "removed", "modified"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in JSON output", want)
		}
	}
}

func TestJSONReport_Empty(t *testing.T) {
	var buf bytes.Buffer
	JSONReport(&buf, []diff.Entry{})
	out := buf.String()
	if !strings.Contains(out, "[") || !strings.Contains(out, "]") {
		t.Error("expected valid JSON array brackets")
	}
}
