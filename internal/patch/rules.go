package patch

// SupportedOps lists all valid patch operations.
var SupportedOps = []Op{OpSet, OpDelete, OpRename}

// IsValidOp returns true if the given op is supported.
func IsValidOp(op Op) bool {
	for _, s := range SupportedOps {
		if s == op {
			return true
		}
	}
	return false
}

// FilterApplied returns only results where Applied is true.
func FilterApplied(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Applied {
			out = append(out, r)
		}
	}
	return out
}

// FilterSkipped returns only results where Applied is false.
func FilterSkipped(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if !r.Applied {
			out = append(out, r)
		}
	}
	return out
}
