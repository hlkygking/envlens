package envdiff

import (
	"strings"
	"testing"
)

func TestFormatEntry_Added(t *testing.T) {
	e := Entry{Key: "NEW_KEY", NewValue: "hello", Status: StatusAdded}
	out := FormatEntry(e, DefaultFormatOptions())
	if !strings.HasPrefix(out, "+ ") {
		t.Errorf("expected prefix '+ ', got %q", out)
	}
	if !strings.Contains(out, "NEW_KEY=hello") {
		t.Errorf("expected key=value in output, got %q", out)
	}
}

func TestFormatEntry_Removed(t *testing.T) {
	e := Entry{Key: "OLD_KEY", OldValue: "bye", Status: StatusRemoved}
	out := FormatEntry(e, DefaultFormatOptions())
	if !strings.HasPrefix(out, "- ") {
		t.Errorf("expected prefix '- ', got %q", out)
	}
	if !strings.Contains(out, "OLD_KEY=bye") {
		t.Errorf("expected old value in output, got %q", out)
	}
}

func TestFormatEntry_Modified(t *testing.T) {
	e := Entry{Key: "FOO", OldValue: "a", NewValue: "b", Status: StatusModified}
	out := FormatEntry(e, DefaultFormatOptions())
	if !strings.HasPrefix(out, "~ ") {
		t.Errorf("expected prefix '~ ', got %q", out)
	}
	if !strings.Contains(out, "a -> b") {
		t.Errorf("expected old->new in output, got %q", out)
	}
}

func TestFormatEntry_ModifiedCompact(t *testing.T) {
	e := Entry{Key: "FOO", OldValue: "a", NewValue: "b", Status: StatusModified}
	opts := DefaultFormatOptions()
	opts.Compact = true
	out := FormatEntry(e, opts)
	if strings.Contains(out, "a ->") {
		t.Errorf("compact mode should not show old value, got %q", out)
	}
	if !strings.Contains(out, "FOO=b") {
		t.Errorf("compact mode should show new value, got %q", out)
	}
}

func TestFormatEntry_Unchanged_Hidden(t *testing.T) {
	e := Entry{Key: "SAME", OldValue: "x", NewValue: "x", Status: StatusUnchanged}
	out := FormatEntry(e, DefaultFormatOptions())
	if out != "" {
		t.Errorf("unchanged entry should produce empty string by default, got %q", out)
	}
}

func TestFormatEntry_Unchanged_Shown(t *testing.T) {
	e := Entry{Key: "SAME", OldValue: "x", NewValue: "x", Status: StatusUnchanged}
	opts := DefaultFormatOptions()
	opts.ShowUnchanged = true
	out := FormatEntry(e, opts)
	if out == "" {
		t.Error("expected non-empty output when ShowUnchanged is true")
	}
}

func TestFormatAll_NoChanges(t *testing.T) {
	entries := []Entry{
		{Key: "A", OldValue: "1", NewValue: "1", Status: StatusUnchanged},
	}
	out := FormatAll(entries, DefaultFormatOptions())
	if out != "(no changes)" {
		t.Errorf("expected '(no changes)', got %q", out)
	}
}

func TestFormatAll_MultipleEntries(t *testing.T) {
	entries := []Entry{
		{Key: "A", NewValue: "1", Status: StatusAdded},
		{Key: "B", OldValue: "2", Status: StatusRemoved},
	}
	out := FormatAll(entries, DefaultFormatOptions())
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %q", len(lines), out)
	}
}
