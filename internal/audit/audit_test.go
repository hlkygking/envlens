package audit

import (
	"testing"

	"github.com/envlens/envlens/internal/diff"
)

func makeEntry(status diff.Status, key, old, new_ string) diff.Entry {
	return diff.Entry{Key: key, Status: status, OldValue: old, NewValue: new_}
}

func findFinding(findings []Finding, key string) *Finding {
	for i := range findings {
		if findings[i].Key == key {
			return &findings[i]
		}
	}
	return nil
}

func TestAudit_SensitiveAdded(t *testing.T) {
	entries := []diff.Entry{makeEntry(diff.Added, "DB_PASSWORD", "", "s3cr3t")}
	findings := Audit(entries)
	f := findFinding(findings, "DB_PASSWORD")
	if f == nil {
		t.Fatal("expected finding for DB_PASSWORD")
	}
	if f.Severity != SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", f.Severity)
	}
}

func TestAudit_SensitiveRemoved(t *testing.T) {
	entries := []diff.Entry{makeEntry(diff.Removed, "API_KEY", "old", "")}
	findings := Audit(entries)
	f := findFinding(findings, "API_KEY")
	if f == nil {
		t.Fatal("expected finding for API_KEY")
	}
	if f.Severity != SeverityHigh {
		t.Errorf("expected HIGH severity, got %s", f.Severity)
	}
}

func TestAudit_SensitiveModified(t *testing.T) {
	entries := []diff.Entry{makeEntry(diff.Modified, "AUTH_TOKEN", "old", "new")}
	findings := Audit(entries)
	f := findFinding(findings, "AUTH_TOKEN")
	if f == nil {
		t.Fatal("expected finding for AUTH_TOKEN")
	}
	if f.Severity != SeverityMedium {
		t.Errorf("expected MEDIUM severity, got %s", f.Severity)
	}
}

func TestAudit_NonSensitiveNoFinding(t *testing.T) {
	entries := []diff.Entry{makeEntry(diff.Added, "APP_PORT", "", "8080")}
	findings := Audit(entries)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d", len(findings))
	}
}

func TestAudit_EmptyToNonEmptyLow(t *testing.T) {
	entries := []diff.Entry{makeEntry(diff.Modified, "LOG_LEVEL", "", "debug")}
	findings := Audit(entries)
	f := findFinding(findings, "LOG_LEVEL")
	if f == nil {
		t.Fatal("expected finding for LOG_LEVEL")
	}
	if f.Severity != SeverityLow {
		t.Errorf("expected LOW severity, got %s", f.Severity)
	}
}
