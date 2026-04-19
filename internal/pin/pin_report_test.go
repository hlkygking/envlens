package pin

import (
	"encoding/json"
	"strings"
	"testing"
)

func samplePinEntries() []Entry {
	return []Entry{
		{Key: "APP_VERSION", Expected: "1.2.3", Actual: "1.2.3", Status: StatusMatch},
		{Key: "DB_HOST", Expected: "prod-db", Actual: "staging-db", Status: StatusMismatch},
		{Key: "SECRET_KEY", Expected: "abc", Actual: "", Status: StatusMissing},
	}
}

func TestPinTextReport_ContainsMismatched(t *testing.T) {
	out := PinTextReport(samplePinEntries())
	if !strings.Contains(out, "MISMATCHED") {
		t.Error("expected MISMATCHED section")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in report")
	}
}

func TestPinTextReport_ContainsMissing(t *testing.T) {
	out := PinTextReport(samplePinEntries())
	if !strings.Contains(out, "MISSING") {
		t.Error("expected MISSING section")
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Error("expected SECRET_KEY in report")
	}
}

func TestPinTextReport_AllMatch(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Expected: "bar", Actual: "bar", Status: StatusMatch},
	}
	out := PinTextReport(entries)
	if !strings.Contains(out, "All pinned keys match") {
		t.Error("expected all-match message")
	}
}

func TestPinTextReport_Summary(t *testing.T) {
	out := PinTextReport(samplePinEntries())
	if !strings.Contains(out, "Summary:") {
		t.Error("expected summary line")
	}
	if !strings.Contains(out, "1 matched") {
		t.Error("expected 1 matched")
	}
}

func TestPinJSONReport_ValidJSON(t *testing.T) {
	out := PinJSONReport(samplePinEntries())
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := v["entries"]; !ok {
		t.Error("expected entries key")
	}
	if _, ok := v["summary"]; !ok {
		t.Error("expected summary key")
	}
}
