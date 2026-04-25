package envdiff

// StatusCounts holds aggregated counts of diff entry statuses.
type StatusCounts struct {
	Added    int
	Removed  int
	Modified int
	Unchanged int
	Total    int
}

// CountByStatus returns a StatusCounts summary over a slice of Entry.
func CountByStatus(entries []Entry) StatusCounts {
	var s StatusCounts
	for _, e := range entries {
		s.Total++
		switch e.Status {
		case StatusAdded:
			s.Added++
		case StatusRemoved:
			s.Removed++
		case StatusModified:
			s.Modified++
		case StatusUnchanged:
			s.Unchanged++
		}
	}
	return s
}

// HasChanges returns true when any added, removed, or modified entries exist.
func (s StatusCounts) HasChanges() bool {
	return s.Added > 0 || s.Removed > 0 || s.Modified > 0
}

// ChangeCount returns the number of non-unchanged entries.
func (s StatusCounts) ChangeCount() int {
	return s.Added + s.Removed + s.Modified
}

// FilterByStatus returns only the entries whose status matches one of the
// provided statuses. Pass no statuses to receive all entries unchanged.
func FilterByStatus(entries []Entry, statuses ...string) []Entry {
	if len(statuses) == 0 {
		return entries
	}
	set := make(map[string]struct{}, len(statuses))
	for _, s := range statuses {
		set[s] = struct{}{}
	}
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if _, ok := set[e.Status]; ok {
			out = append(out, e)
		}
	}
	return out
}
