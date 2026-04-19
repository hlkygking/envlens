package sort

import (
	"sort"
	"strings"
)

// SortOrder defines how entries should be sorted.
type SortOrder string

const (
	OrderAlpha  SortOrder = "alpha"
	OrderAlphaR SortOrder = "alpha-desc"
	OrderLength SortOrder = "length"
	OrderGroup  SortOrder = "group"
)

// Entry represents a key-value pair with optional group tag.
type Entry struct {
	Key   string
	Value string
	Group string
}

// Result holds the output of a sort operation.
type Result struct {
	Entries []Entry
	Order   SortOrder
	Total   int
}

// Apply sorts the given entries by the specified order.
func Apply(entries []Entry, order SortOrder) Result {
	out := make([]Entry, len(entries))
	copy(out, entries)

	switch order {
	case OrderAlphaR:
		sort.Slice(out, func(i, j int) bool {
			return out[i].Key > out[j].Key
		})
	case OrderLength:
		sort.Slice(out, func(i, j int) bool {
			if len(out[i].Key) != len(out[j].Key) {
				return len(out[i].Key) < len(out[j].Key)
			}
			return out[i].Key < out[j].Key
		})
	case OrderGroup:
		sort.Slice(out, func(i, j int) bool {
			if out[i].Group != out[j].Group {
				return out[i].Group < out[j].Group
			}
			return out[i].Key < out[j].Key
		})
	default: // alpha
		sort.Slice(out, func(i, j int) bool {
			return strings.ToLower(out[i].Key) < strings.ToLower(out[j].Key)
		})
	}

	return Result{Entries: out, Order: order, Total: len(out)}
}

// ToMap converts entries to a key-value map.
func ToMap(entries []Entry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
