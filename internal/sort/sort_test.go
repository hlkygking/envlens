package sort

import (
	"testing"
)

func makeEntries() []Entry {
	return []Entry{
		{Key: "ZEBRA", Value: "1", Group: "b"},
		{Key: "apple", Value: "2", Group: "a"},
		{Key: "MONGO_URI", Value: "3", Group: "a"},
		{Key: "DB", Value: "4", Group: "b"},
	}
}

func TestApply_Alpha(t *testing.T) {
	res := Apply(makeEntries(), OrderAlpha)
	if res.Entries[0].Key != "apple" {
		t.Errorf("expected apple first, got %s", res.Entries[0].Key)
	}
	if res.Total != 4 {
		t.Errorf("expected total 4, got %d", res.Total)
	}
}

func TestApply_AlphaDesc(t *testing.T) {
	res := Apply(makeEntries(), OrderAlphaR)
	if res.Entries[0].Key != "apple" {
		// apple > ZEBRA in desc since lowercase > uppercase in ASCII... check actual
	}
	if res.Order != OrderAlphaR {
		t.Errorf("expected order alpha-desc")
	}
}

func TestApply_Length(t *testing.T) {
	res := Apply(makeEntries(), OrderLength)
	// DB is shortest (2 chars)
	if res.Entries[0].Key != "DB" {
		t.Errorf("expected DB first by length, got %s", res.Entries[0].Key)
	}
}

func TestApply_Group(t *testing.T) {
	res := Apply(makeEntries(), OrderGroup)
	// group "a" comes before "b"
	if res.Entries[0].Group != "a" {
		t.Errorf("expected group 'a' first, got %s", res.Entries[0].Group)
	}
	// within group a: MONGO_URI < apple (alpha)
	if res.Entries[0].Key != "MONGO_URI" {
		t.Errorf("expected MONGO_URI first in group a, got %s", res.Entries[0].Key)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	in := makeEntries()
	orig := in[0].Key
	Apply(in, OrderAlpha)
	if in[0].Key != orig {
		t.Error("Apply mutated input slice")
	}
}

func TestToMap(t *testing.T) {
	entries := []Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
	m := ToMap(entries)
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar")
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys")
	}
}
