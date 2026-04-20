package envdrift

// ScopedApply runs Apply against a named set of environment maps,
// treating the first entry in envMaps as the baseline.
// The key of each entry is the environment name (e.g. "prod", "staging").
// The first key in order is used as the baseline; remaining entries are targets.
func ScopedApply(envMaps map[string]map[string]string, baselineName string) ([]Result, error) {
	baseline, ok := envMaps[baselineName]
	if !ok {
		return nil, &ErrMissingBaseline{Name: baselineName}
	}

	targets := make(map[string]map[string]string, len(envMaps)-1)
	for name, m := range envMaps {
		if name == baselineName {
			continue
		}
		targets[name] = m
	}

	return Apply(baseline, targets), nil
}

// FilterByStatus returns only results with the given status.
func FilterByStatus(results []Result, status Status) []Result {
	out := make([]Result, 0)
	for _, r := range results {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}

// ErrMissingBaseline is returned when the specified baseline env name is not found.
type ErrMissingBaseline struct {
	Name string
}

func (e *ErrMissingBaseline) Error() string {
	return "envdrift: baseline environment not found: " + e.Name
}
