package envcount

import "strings"

// Status represents the count category of a key.
type Status string

const (
	StatusEmpty    Status = "empty"
	StatusShort    Status = "short"
	StatusMedium   Status = "medium"
	StatusLong     Status = "long"
	StatusSensitive Status = "sensitive"
)

// Result holds the analysis of a single env key.
type Result struct {
	Key       string
	Value     string
	Length    int
	Status    Status
	Sensitive bool
}

// Summary holds aggregate counts.
type Summary struct {
	Total     int
	Empty     int
	Short     int
	Medium    int
	Long      int
	Sensitive int
}

var sensitivePatterns = []string{
	"PASSWORD", "SECRET", "TOKEN", "KEY", "APIKEY", "API_KEY",
	"PRIVATE", "CREDENTIAL", "AUTH", "PASS",
}

func isSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

func classify(length int) Status {
	switch {
	case length == 0:
		return StatusEmpty
	case length <= 8:
		return StatusShort
	case length <= 64:
		return StatusMedium
	default:
		return StatusLong
	}
}

// Apply analyses each key-value pair and returns results with length metadata.
func Apply(env map[string]string) []Result {
	results := make([]Result, 0, len(env))
	for k, v := range env {
		sens := isSensitive(k)
		results = append(results, Result{
			Key:       k,
			Value:     v,
			Length:    len(v),
			Status:    classify(len(v)),
			Sensitive: sens,
		})
	}
	return results
}

// GetSummary aggregates result counts.
func GetSummary(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusEmpty:
			s.Empty++
		case StatusShort:
			s.Short++
		case StatusMedium:
			s.Medium++
		case StatusLong:
			s.Long++
		}
		if r.Sensitive {
			s.Sensitive++
		}
	}
	return s
}
